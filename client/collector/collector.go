package collector

import (
	"context"

	"github.com/flare-foundation/go-flare-common/pkg/contracts/registry"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/relay"
	"github.com/flare-foundation/go-flare-common/pkg/contracts/submission"

	"github.com/flare-foundation/go-flare-common/pkg/contracts/fdchub"
	"github.com/flare-foundation/go-flare-common/pkg/database"
	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/go-flare-common/pkg/payload"

	"gitlab.com/flarenetwork/fdc/fdc-client/client/config"
	"gitlab.com/flarenetwork/fdc/fdc-client/client/shared"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

const (
	bitVoteOffChainTriggerSeconds = 15
	outOfSyncTolerance            = 60
	requestListenerInterval       = 2 * time.Second
	databasePollTime              = 2 * time.Second
	bitVoteHeadStart              = 5 * time.Second
)

var signingPolicyInitializedEventSel common.Hash
var AttestationRequestEventSel common.Hash
var voterRegisteredEventSel common.Hash

var Submit2FuncSel [4]byte

func init() {
	relayABI, err := relay.RelayMetaData.GetAbi()
	if err != nil {
		logger.Panic("cannot get relayAby:", err)
	}

	signingPolicyEvent, ok := relayABI.Events["SigningPolicyInitialized"]
	if !ok {
		logger.Panic("cannot get SigningPolicyInitialized event:", err)
	}
	signingPolicyInitializedEventSel = signingPolicyEvent.ID

	fdcABI, err := fdchub.FdcHubMetaData.GetAbi()

	if err != nil {
		logger.Panic("cannot get fdcABI:", err)
	}

	requestEvent, ok := fdcABI.Events["AttestationRequest"]
	if !ok {
		logger.Panic("cannot get AttestationRequest event:", err)
	}

	AttestationRequestEventSel = requestEvent.ID

	registryABI, err := registry.RegistryMetaData.GetAbi()
	if err != nil {
		logger.Panic("cannot get registryABI:", err)
	}

	voterRegisteredEvent, ok := registryABI.Events["VoterRegistered"]

	if !ok {
		logger.Panic("cannot get VoterRegistered event:", err)
	}

	voterRegisteredEventSel = voterRegisteredEvent.ID

	submissionABI, err := submission.SubmissionMetaData.GetAbi()
	if err != nil {
		logger.Panic("cannot get submission ABI:", err)
	}
	copy(Submit2FuncSel[:], submissionABI.Methods["submit2"].ID[:4])
}

type Collector struct {
	ProtocolID                   uint8
	SubmitContractAddress        common.Address
	FdcContractAddress           common.Address
	RelayContractAddress         common.Address
	VoterRegistryContractAddress common.Address

	DB              *gorm.DB
	Requests        chan<- []database.Log
	BitVotes        chan<- payload.Round
	SigningPolicies chan<- []shared.VotersData
}

// New creates new Collector from user and system configs.
func New(user *config.UserRaw, system *config.System, sharedDataPipes *shared.DataPipes) *Collector {
	db, err := database.Connect(&user.DB)
	if err != nil {
		logger.Panic("Could not connect to database:", err)
	}

	runner := Collector{
		ProtocolID:                   user.ProtocolID,
		SubmitContractAddress:        system.Addresses.SubmitContract,
		FdcContractAddress:           system.Addresses.FdcContract,
		RelayContractAddress:         system.Addresses.RelayContract,
		VoterRegistryContractAddress: system.Addresses.VoterRegistryContract,

		DB:              db,
		SigningPolicies: sharedDataPipes.Voters,
		BitVotes:        sharedDataPipes.BitVotes,
		Requests:        sharedDataPipes.Requests,
	}

	return &runner
}

// Run starts SigningPolicyInitializedListener, BitVoteListener, and AttestationRequestListener,
// assigns their channels to the RoundManager, and starts the RoundManager.
func (c *Collector) Run(ctx context.Context) {
	state, err := database.FetchState(ctx, c.DB, nil)
	if err != nil {
		logger.Panic("database error:", err)
	}

	if k := time.Now().Unix() - int64(state.BlockTimestamp); k > outOfSyncTolerance {
		logger.Panicf("database not up to date. lags for %d minutes", k/60)
	}

	go SigningPolicyInitializedListener(ctx, c.DB, c.RelayContractAddress, c.VoterRegistryContractAddress, c.SigningPolicies)
	go AttestationRequestListener(ctx, c.DB, c.FdcContractAddress, requestListenerInterval, c.Requests)

	chooseTrigger := make(chan uint32)
	go BitVoteListener(ctx, c.DB, c.SubmitContractAddress, Submit2FuncSel, c.ProtocolID, chooseTrigger, c.BitVotes)
	PrepareChooseTrigger(ctx, chooseTrigger, c.DB)
}
