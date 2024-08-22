package collector

import (
	"context"
	"encoding/hex"
	"errors"
	"flare-common/contracts/registry"
	"flare-common/contracts/relay"

	"flare-common/database"
	"flare-common/logger"
	"flare-common/payload"
	"local/fdc/client/config"
	"local/fdc/client/shared"
	"local/fdc/contracts/fdc"

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
	Submit1FuncSelHex             = "6c532fae"
)

var signingPolicyInitializedEventSel common.Hash
var AttestationRequestEventSel common.Hash
var voterRegisteredEventSel common.Hash

var log = logger.GetLogger()

func init() {
	relayABI, err := relay.RelayMetaData.GetAbi()
	if err != nil {
		log.Panic("cannot get relayAby:", err)
	}

	signingPolicyEvent, ok := relayABI.Events["SigningPolicyInitialized"]
	if !ok {
		log.Panic("cannot get SigningPolicyInitialized event:", err)
	}
	signingPolicyInitializedEventSel = signingPolicyEvent.ID

	fdcABI, err := fdc.FdcMetaData.GetAbi()

	if err != nil {
		log.Panic("cannot get fdcABI:", err)
	}

	requestEvent, ok := fdcABI.Events["AttestationRequest"]
	if !ok {
		log.Panic("cannot get AttestationRequest event:", err)
	}

	AttestationRequestEventSel = requestEvent.ID

	registryABI, err := registry.RegistryMetaData.GetAbi()
	if err != nil {
		log.Panic("cannot get registryABI:", err)
	}

	voterRegisteredEvent, ok := registryABI.Events["VoterRegistered"]

	if !ok {
		log.Panic("cannot get VoterRegistered event:", err)
	}

	voterRegisteredEventSel = voterRegisteredEvent.ID

}

type Collector struct {
	ProtocolID                   uint8
	SubmitContractAddress        common.Address
	FdcContractAddress           common.Address
	RelayContractAddress         common.Address
	VoterRegistryContractAddress common.Address

	DB              *gorm.DB
	submit1Sel      [4]byte
	Requests        chan []database.Log
	BitVotes        chan payload.Round
	SigningPolicies chan []shared.VotersData
}

// New creates new Collector from user and system configs.
func New(user *config.UserRaw, system *config.System, sharedDataPipes *shared.DataPipes) *Collector {
	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Panic("Could not connect to database:", err)
	}

	submit1FuncSel, err := ParseFuncSel(Submit1FuncSelHex)
	if err != nil {
		log.Panic("Could not parse submit1FuncSel:", err)
	}

	runner := Collector{
		ProtocolID:                   user.ProtocolID,
		SubmitContractAddress:        system.Addresses.SubmitContract,
		FdcContractAddress:           system.Addresses.FdcContract,
		RelayContractAddress:         system.Addresses.RelayContract,
		VoterRegistryContractAddress: system.Addresses.VoterRegistryContract,

		DB:              db,
		submit1Sel:      submit1FuncSel,
		SigningPolicies: sharedDataPipes.Voters,
		BitVotes:        sharedDataPipes.BitVotes,
		Requests:        sharedDataPipes.Requests,
	}

	return &runner
}

func ParseFuncSel(sigInput string) ([4]byte, error) {
	var ret [4]byte
	inputBytes := []byte(sigInput)

	if hex.DecodedLen(len(inputBytes)) != 4 {
		return ret, errors.New("invalid length for function selector")
	}

	_, err := hex.Decode(ret[:], inputBytes)
	return ret, err
}

// Run starts SigningPolicyInitializedListener, BitVoteListener, and AttestationRequestListener,
// assigns their channels to the RoundManager, and starts the RoundManager.
func (c *Collector) Run(ctx context.Context) {
	state, err := database.FetchState(ctx, c.DB)
	if err != nil {
		log.Panic("database error:", err)
	}

	if k := time.Now().Unix() - int64(state.BlockTimestamp); k > outOfSyncTolerance {
		log.Panicf("database not up to date. lags for %d minutes", k/60)
	}

	chooseTrigger := make(chan uint64)

	go SigningPolicyInitializedListener(ctx, c.DB, c.RelayContractAddress, c.VoterRegistryContractAddress, c.SigningPolicies)
	go BitVoteListener(ctx, c.DB, c.SubmitContractAddress, c.submit1Sel, c.ProtocolID, chooseTrigger, c.BitVotes)
	go AttestationRequestListener(ctx, c.DB, c.FdcContractAddress, requestListenerInterval, c.Requests)

	PrepareChooseTriggers(ctx, chooseTrigger, c.DB)
}
