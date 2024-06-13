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
	Protocol              uint8
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

type collectorDB interface {
	FetchState(context.Context) (database.State, error)

	FetchLatestLogsByAddressAndTopic0(
		context.Context, common.Address, common.Hash, int,
	) ([]database.Log, error)

	FetchLogsByAddressAndTopic0Timestamp(
		context.Context, common.Address, common.Hash, int64, int64,
	) ([]database.Log, error)

	FetchLogsByAddressAndTopic0BlockNumber(
		context.Context, common.Address, common.Hash, int64, int64,
	) ([]database.Log, error)

	FetchLogsByAddressAndTopic0TimestampToBlockNumber(
		context.Context,
		common.Address,
		common.Hash,
		int64,
		int64,
	) ([]database.Log, error)

	FetchTransactionsByAddressAndSelectorTimestamp(
		context.Context, common.Address, [4]byte, int64, int64,
	) ([]database.Transaction, error)
}

type collectorDBGorm struct {
	db *gorm.DB
}

func (c collectorDBGorm) FetchState(ctx context.Context) (database.State, error) {
	return database.FetchState(ctx, c.db)
}
func (c collectorDBGorm) FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, addr common.Address, topic0 common.Hash, num int,
) ([]database.Log, error) {
	return database.FetchLatestLogsByAddressAndTopic0(ctx, c.db, addr, topic0, num)
}

func (c collectorDBGorm) FetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context, addr common.Address, topic0 common.Hash, from, to int64,
) ([]database.Log, error) {
	return database.FetchLogsByAddressAndTopic0Timestamp(ctx, c.db, addr, topic0, from, to)
}

func (c collectorDBGorm) FetchLogsByAddressAndTopic0BlockNumber(
	ctx context.Context,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]database.Log, error) {
	return database.FetchLogsByAddressAndTopic0BlockNumber(ctx, c.db, address, topic0, from, to)
}

func (c collectorDBGorm) FetchLogsByAddressAndTopic0TimestampToBlockNumber(
	ctx context.Context,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]database.Log, error) {
	return database.FetchLogsByAddressAndTopic0TimestampToBlockNumber(
		ctx, c.db, address, topic0, from, to,
	)
}

func (c collectorDBGorm) FetchTransactionsByAddressAndSelectorTimestamp(
	ctx context.Context,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]database.Transaction, error) {
	return database.FetchTransactionsByAddressAndSelectorTimestamp(
		ctx, c.db, toAddress, functionSel, from, to,
	)
}

func (r *Collector) Run(ctx context.Context) {
	state, err := database.FetchState(ctx, r.DB)
	if err != nil {
		log.Panic("database error:", err)
	}

	if k := time.Now().Unix() - int64(state.BlockTimestamp); k > 60 { //get 60 from config
		log.Panic("database not up to date")
	}

	chooseTrigger := make(chan uint64)

	db := collectorDBGorm{db: r.DB}

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(ctx, db, r.RelayContractAddress, 3)

	r.RoundManager.BitVotes = BitVoteListener(ctx, db, r.FdcContractAddress, r.submit1Sel, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = AttestationRequestListener(ctx, db, r.FdcContractAddress, requestsBufferSize, requestListenerInterval)

	go r.RoundManager.Run(ctx)

	prepareChooseTriggers(ctx, chooseTrigger, r.DB)

}
