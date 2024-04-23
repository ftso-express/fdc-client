package payload

import (
	"errors"
	"flare-common/database"
	"strconv"
)

type Message struct {
	From        string
	Signature   string
	Protocol    uint64
	VotingRound uint64
	Length      uint64
	Payload     string
}

func ExtractPayloads(tx database.Transaction, protocol string) (map[uint64]Message, error) {

	messages := make(map[uint64]Message)

	data := tx.Input[10:]

	for len(data) > 0 {
		protocol, err := strconv.ParseUint(data[:2], 16, 64)

		if err != nil {
			return nil, errors.New("protocol id error")
		}
		votingRound, err := strconv.ParseUint(data[2:10], 16, 64)
		if err != nil {
			return nil, errors.New("voting round error")
		}

		length, err := strconv.ParseUint(data[10:14], 16, 64)
		if err != nil {
			return nil, errors.New("length error")
		}

		end := 14 + 2*length
		payload := data[14:end]

		message := Message{tx.FromAddress, tx.FunctionSig, protocol, votingRound, length, payload}

		messages[protocol] = message

		data = data[end:]

	}

	return messages, nil
}
