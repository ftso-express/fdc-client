package collector

import (
	"flare-common/database"
	"local/fdc/client/timing"

	"gorm.io/gorm"
)

type Runner struct {
	Protocol              string
	SubmitContractAddress string
	FdcContractAddress    string
	RequestEventSig       string
	DB                    *gorm.DB
	submit1Sig            string
}

func (r *Runner) Run() {

	var initialStart = timing.RoundLatest().Start.Unix()

	state, _ := database.FetchState(r.DB)

	latestTimestamp := int64(state.BlockTimestamp)
	latestBlock := state.Index

	requestEvents, _ := database.FetchLogsByAddressAndTopic0Timestamp(r.DB, r.FdcContractAddress, r.RequestEventSig, initialStart-1, latestTimestamp)

	parseRequests(requestEvents)

	for {
		state, _ := database.FetchState(r.DB)

		requestEvents, _ := database.FetchLogsByAddressAndTopic0BlockNumber(r.DB, r.FdcContractAddress, r.RequestEventSig, int64(latestBlock), int64(state.Index))

		parseRequests(requestEvents)

		bitvotes, _ := database.FetchTransactionsByAddressAndSelectorBlockNumber(r.DB, r.SubmitContractAddress, r.submit1Sig, int64(latestBlock), int64(state.Index))

		parseBitVotes(bitvotes)

		latestTimestamp = int64(state.BlockTimestamp)
		latestBlock = state.Index
	}
}
