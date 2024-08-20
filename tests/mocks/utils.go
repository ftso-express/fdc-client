package mocks

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetFailingMessage(client ethclient.Client, hash common.Hash) (string, error) {
	tx, _, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return "", err
	}

	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func PrivKeyToAddress(privateKey string) (common.Address, *ecdsa.PrivateKey) {
	if privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return fromAddress, privateKeyECDSA
}
