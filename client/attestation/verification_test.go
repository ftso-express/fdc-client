package attestation_test

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"local/fdc/client/attestation"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const response string = "42616c616e636544656372656173696e675472616e73616374696f6e000000004254430000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000664cbf4c2a3ce5fb95fa6b436fbed49cbccc6dcbb9ee166a3ef217d227cbe5add6898dd20000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000664cbf4c06fa5d68b3284548b849dca2ffd9a59350c7440c5be121fe4b4ae0941dcae638000000000000000000000000000000000000000000000000000000000131a3c0000000000000000000000000000000000000000000000000000000add6898dd2"

func TestIsStaticType(t *testing.T) {

	_, err := attestation.IsStaticType([]byte{1, 1, 1})

	if err == nil {
		t.Error("fail")

	}

	resp, _ := hex.DecodeString(response)

	ok, err := attestation.IsStaticType(resp)

	if err != nil || !ok {
		t.Error("fail")

	}

}

func TestAddRound(t *testing.T) {

	var resp attestation.Response

	resp, _ = hex.DecodeString(response)

	resp, _ = resp.AddRound(9)

	if resp[95] != byte(9) {
		t.Error("fail")
	}

	resp.AddRound(257)

	if resp[95] != byte(1) || resp[94] != byte(1) {
		t.Error("fail")
	}

}

func TestComputeMic(t *testing.T) {

	var resp attestation.Response

	resp, _ = hex.DecodeString(response)

	file, err := os.ReadFile("../../configs/abis/BalanceDecreasingTransaction.json")

	fmt.Println(err)

	var arg abi.Argument

	err = arg.UnmarshalJSON(file)

	fmt.Println(err)

	args := abi.Arguments{arg}

	mic, _ := resp.ComputeMic(args)

	if mic.String() != "0x2f51362aef7ff57fa4aa74ecca9a5fbaffc123416d7df97226e8635776f06d0b" {
		t.Error("wrong mic")
	}

}

func TestHash(t *testing.T) {

	var resp attestation.Response

	resp, _ = hex.DecodeString(response)

	hash, _ := resp.Hash(0)

	if hash.String() != "0xf014a6220c448b59f8067a99397b4d3506d17a182df101f23c52825e6c5e1f17" {
		t.Error("wrong hash")
	}
}
