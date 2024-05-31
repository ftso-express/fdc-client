package collector

import (
	"flare-common/database"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// AttestationRequestListener returns a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	db *gorm.DB, fdcContractAddress common.Address, bufferSize int, ListenerInterval time.Duration,
) <-chan []database.Log {

	out := make(chan []database.Log, bufferSize)

	go func() {

		trigger := time.NewTicker(ListenerInterval)

		_, startTimestamp := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))

		state, err := database.FetchState(db)
		if err != nil {
			log.Panic("fetch initial state error:", err)
		}

		lastQueriedBlock := state.Index

		logs, err := database.FetchLogsByAddressAndTopic0TimestampToBlockNumber(
			db, fdcContractAddress, attestationRequestEventSel, int64(startTimestamp), int64(state.Index),
		)
		if err != nil {
			log.Panic("fetch initial logs error")
		}

		if len(logs) > 0 {
			out <- logs
		}

		for {
			<-trigger.C

			state, err = database.FetchState(db)
			if err != nil {
				log.Error("fetch state error:", err)
				continue
			}

			logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
				db, fdcContractAddress, attestationRequestEventSel, int64(lastQueriedBlock), int64(state.Index),
			)
			if err != nil {
				log.Error("fetch logs error:", err)
				continue
			}

			lastQueriedBlock = state.Index

			if len(logs) > 0 {
				log.Debugf("Adding %d request logs to channel", len(logs))
				out <- logs
			}

		}

	}()

	return out
}
