package policy

import (
	"flare-common/contacts/relay"
	"flare-common/database"
	"flare-common/events"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

const (
	listenerBufferSize               = 10
	ListenerInterval   time.Duration = 2 * time.Second
)

type signingPolicyListenerResponse struct {
	policyData *relay.RelaySigningPolicyInitialized
	timestamp  int64
}

type relayContractClient struct {
	address common.Address

	relay *relay.Relay

	topic0SPI string // for SigningPolicyInitialized event
	topic0PMR string // for ProtocolMessageRelayed event
}

func (r *relayContractClient) FetchSigningPolicies(db *gorm.DB, from, to int64) ([]signingPolicyListenerResponse, error) {
	logs, err := database.FetchLogsByAddressAndTopic0Timestamp(db, r.address.String(), r.topic0SPI, from, to)
	if err != nil {
		return nil, err
	}

	result := make([]signingPolicyListenerResponse, 0, len(logs))
	for _, log := range logs {
		policyData, err := events.ParseSigningPolicyInitializedEvent(r.relay, log)
		if err != nil {
			return nil, err
		}
		result = append(result, signingPolicyListenerResponse{policyData, int64(log.Timestamp)})
	}
	return result, nil
}

func (r *relayContractClient) SigningPolicyInitializedListener(db *gorm.DB, startTime time.Time) <-chan signingPolicyListenerResponse {
	out := make(chan signingPolicyListenerResponse, listenerBufferSize)
	go func() {
		ticker := time.NewTicker(ListenerInterval)
		eventRangeStart := startTime.Unix()
		for {
			<-ticker.C
			now := time.Now().Unix()
			logs, err := database.FetchLogsByAddressAndTopic0Timestamp(db, r.address.String(), r.topic0SPI, eventRangeStart, now)
			if err != nil {
				continue
			}
			for _, log := range logs {
				policyData, err := events.ParseSigningPolicyInitializedEvent(r.relay, log)
				if err != nil {
					break
				}
				out <- signingPolicyListenerResponse{policyData, int64(log.Timestamp)}
				// continue with timestamps > log.Timestamp,
				// there should be only one such log per timestamp
				eventRangeStart = int64(log.Timestamp)
			}
		}
	}()
	return out
}
