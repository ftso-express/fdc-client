package collector

import (
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/attestation"
	"local/fdc/client/timing"

	"gorm.io/gorm"
)

type Runner struct {
	Protocol              uint64
	SubmitContractAddress string
	RequestEventSig       string
	FdcContractAddress    string
	RelayContractAddress  string
	SigningPolicyEventSig string
	DB                    *gorm.DB
	submit1Sig            string
	roundManager          *attestation.Manager
}

func (r *Runner) Run() {

	var initialStart = timing.RoundLatest().Start.Unix()

	state, _ := database.FetchState(r.DB)

	latestTimestamp := int64(state.BlockTimestamp)
	latestBlock := state.Index

	r.roundManager.Timestamps <- state.BlockTimestamp

	requestEvents, _ := database.FetchLogsByAddressAndTopic0Timestamp(r.DB, r.FdcContractAddress, r.RequestEventSig, initialStart-1, latestTimestamp)

	processRequests(requestEvents)

	for {
		state, _ := database.FetchState(r.DB)

		r.roundManager.Timestamps <- state.BlockTimestamp

		// requestEvents, _ := database.FetchLogsByAddressAndTopic0BlockNumber(r.DB, r.FdcContractAddress, r.RequestEventSig, int64(latestBlock), int64(state.Index))

		// processRequests(requestEvents)

		submit1Txs, _ := database.FetchTransactionsByAddressAndSelectorBlockNumber(r.DB, r.SubmitContractAddress, r.submit1Sig, int64(latestBlock), int64(state.Index))

		for _, tx := range submit1Txs {
			payloads, _ := payload.ExtractPayloads(tx)
			bitvote := payloads[r.Protocol]

			r.roundManager.BitVotes <- bitvote
		}

		//latestTimestamp = int64(state.BlockTimestamp)
		latestBlock = state.Index
	}

}

func processRequests(logs []database.Log) {
}
