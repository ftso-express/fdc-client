package payload

import (
	"errors"
	"flare-common/database"
	"strconv"
	"strings"
)

type Round struct {
	Messages []Message
	ID       uint64
}

type Message struct {
	From             string
	Signature        string
	Protocol         uint64
	VotingRound      uint64
	Timestamp        uint64
	BlockNumber      uint64
	TransactionIndex uint64
	Length           uint64
	Payload          string
}

func ExtractPayloads(tx database.Transaction) (map[uint64]Message, error) {

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

		message := Message{tx.FromAddress, tx.FunctionSig, protocol, votingRound, tx.Timestamp, tx.BlockNumber, tx.TransactionIndex, length, payload}

		messages[protocol] = message

		data = data[end:]

	}

	return messages, nil
}

func prependZerosToLength(hexString string, finalLength int) (string, error) {
	p := len(hexString) - finalLength

	if p < 0 {
		return hexString, errors.New("string too long")
	}

	prefix := strings.Repeat("0", p)

	return prefix + hexString, nil

}

func BuildMessage(protocol, votingRound uint64, payload string) (string, error) {

	if len(payload)%2 != 0 {
		return "", errors.New("uneven payload")
	}

	protocolStr, err := prependZerosToLength(strconv.FormatUint(protocol, 16), 2)

	if err != nil {
		return "", errors.New("invalid protocol")
	}

	votingRoundStr, err := prependZerosToLength(strconv.FormatUint(votingRound, 16), 8)

	if err != nil {
		return "", errors.New("invalid voting round")
	}

	lenStr, err := prependZerosToLength(strconv.FormatInt(int64(len(payload)), 16), 4)

	if err != nil {
		return "", errors.New("invalid payload length")
	}

	return protocolStr + votingRoundStr + lenStr + payload, nil
}
