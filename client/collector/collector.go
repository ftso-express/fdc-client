package collector

import (
	"context"
	"encoding/hex"
	"errors"
	"flare-common/contracts/relay"
	"flare-common/database"
	"flare-common/logger"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	hub "local/fdc/contracts/FDC"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

const (
	bitVoteBufferSize             = 10
	bitVoteOffChainTriggerSeconds = 15
	requestsBufferSize            = 10
	requestListenerInterval       = 2 * time.Second
	signingPolicyBufferSize       = 3
)

const (
	roundLength      = 90 * time.Second
	databasePollTime = 2 * time.Second
	bitVoteHeadStart = 5 * time.Second
)

var signingPolicyInitializedEventSel common.Hash
var attestationRequestEventSel common.Hash
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

	fdcAbi, err := hub.HubMetaData.GetAbi()

	if err != nil {
		log.Panic("cannot get fdcAbi:", err)
	}

	requestEvent, ok := fdcAbi.Events["AttestationRequest"]

	if !ok {
		log.Panic("cannot get AttestationRequest event:", err)
	}

	attestationRequestEventSel = requestEvent.ID

}

type Collector struct {
	Protocol              uint64
	SubmitContractAddress common.Address
	FdcContractAddress    common.Address
	RelayContractAddress  common.Address
	DB                    *gorm.DB
	submit1Sel            [4]byte
	RoundManager          *attestation.Manager
}

const submit1FuncSelHex = "6c532fae"

func New(user config.UserConfigRaw, system config.SystemConfig) *Collector {
	// CONSTANTS
	// requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	// signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"

	roundManager := attestation.NewManager(user)

	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Panic("Could not connect to database:", err)
	}

	submit1FuncSel, err := parseFuncSel(submit1FuncSelHex)
	if err != nil {
		log.Panic("Could not parse submit1FuncSel:", err)
	}

	runner := Collector{
		Protocol:              system.Listener.Protocol,
		SubmitContractAddress: system.Listener.SubmitContractAddress,
		FdcContractAddress:    system.Listener.FdcContractAddress,
		RelayContractAddress:  system.Listener.RelayContractAddress,
		DB:                    db,
		submit1Sel:            submit1FuncSel,
		RoundManager:          roundManager,
	}

	return &runner
}

func parseFuncSel(sigInput string) ([4]byte, error) {
	var ret [4]byte
	inputBytes := []byte(sigInput)

	if hex.DecodedLen(len(inputBytes)) != 4 {
		return ret, errors.New("invalid length for function selector")
	}

	_, err := hex.Decode(ret[:], inputBytes)
	return ret, err
}

func (r *Collector) Run(ctx context.Context) {

	chooseTrigger := make(chan uint64)

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(ctx, r.DB, r.RelayContractAddress, 3)

	r.RoundManager.BitVotes = BitVoteListener(ctx, r.DB, r.FdcContractAddress, r.submit1Sel, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = AttestationRequestListener(ctx, r.DB, r.FdcContractAddress, requestsBufferSize, requestListenerInterval)

	go r.RoundManager.Run(ctx)

	prepareChooseTriggers(ctx, chooseTrigger, r.DB)

}
