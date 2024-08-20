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
	relayAbi, err := relay.RelayMetaData.GetAbi()
	if err != nil {
		log.Panic("cannot get relayAby:", err)
	}

	signingPolicyEvent, ok := relayAbi.Events["SigningPolicyInitialized"]
	if !ok {
		log.Panic("cannot get SigningPolicyInitialized event:", err)
	}
	signingPolicyInitializedEventSel = signingPolicyEvent.ID

	fdcAbi, err := fdc.FdcMetaData.GetAbi()

	if err != nil {
		log.Panic("cannot get fdcAbi:", err)
	}

	requestEvent, ok := fdcAbi.Events["AttestationRequest"]
	if !ok {
		log.Panic("cannot get AttestationRequest event:", err)
	}

	AttestationRequestEventSel = requestEvent.ID

	registryAbi, err := registry.RegistryMetaData.GetAbi()
	if err != nil {
		log.Panic("cannot get registryAbi:", err)
	}

	voterRegisteredEvent, ok := registryAbi.Events["VoterRegistered"]

	if !ok {
		log.Panic("cannot get VoterRegistered event:", err)
	}

	voterRegisteredEventSel = voterRegisteredEvent.ID

}

type Collector struct {
	ProtocolId                   uint8
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

// NewCollector creates new Collector from user and system configs.
func NewCollector(user *config.UserRaw, system *config.System, sharedDataPipes *shared.SharedDataPipes) *Collector {
	// CONSTANTS
	// requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	// signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"

	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Panic("Could not connect to database:", err)
	}

	submit1FuncSel, err := ParseFuncSel(Submit1FuncSelHex)
	if err != nil {
		log.Panic("Could not parse submit1FuncSel:", err)
	}

	runner := Collector{
		ProtocolId:                   user.ProtocolId,
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
func (r *Collector) Run(ctx context.Context) {
	state, err := database.FetchState(ctx, r.DB)
	if err != nil {
		log.Panic("database error:", err)
	}

	if k := time.Now().Unix() - int64(state.BlockTimestamp); k > outOfSyncTolerance {
		log.Panicf("database not up to date. lags for %d minutes", k/60)
	}

	chooseTrigger := make(chan uint64)

	go SigningPolicyInitializedListener(ctx, r.DB, r.RelayContractAddress, r.VoterRegistryContractAddress, r.SigningPolicies)
	go BitVoteListener(ctx, r.DB, r.SubmitContractAddress, r.submit1Sel, r.ProtocolId, chooseTrigger, r.BitVotes)
	go AttestationRequestListener(ctx, r.DB, r.FdcContractAddress, requestListenerInterval, r.Requests)

	PrepareChooseTriggers(ctx, chooseTrigger, r.DB)
}
