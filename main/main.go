package main

import (
	"encoding/hex"
	"fmt"
	"local/fdc/client/attestation"
	"local/fdc/client/verification"
	verificationServer "local/fdc/client/verifierServer"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
)

// func main() {

// 	// Prepare context
// 	// Empty context
// 	context := context.Background()

// 	cancelChan := make(chan os.Signal, 1)
// 	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)

// 	// run server
// 	go server.RunProviderServer(context)

// 	<-cancelChan
// 	fmt.Printf("Shutting down server")
// }

// Temp main testing
// TODO: Luka Move this to test (assuming local evm verifier)
func main() {

	validRequestData := "45564d5472616e73616374696f6e00000000000000000000000000000000000045544800000000000000000000000000000000000000000000000000000000002f9e9bc1059c5c49403b5cdbe2b314787f626797d98c9ef101ecbe0786106c9f0000000000000000000000000000000000000000000000000000000000000020ef4514befee7f686f494273a7df083f180e76459b300bcfaf8fb8d3ae1b55f3800000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000"

	req, _ := hex.DecodeString(validRequestData)

	fee := big.NewInt(0)
	hash := crypto.Keccak256Hash([]byte("Hello, World!"))

	request := verification.Request(req)

	attestation := attestation.Attestation{
		Index:     attestation.IndexLog{BlockNumber: 1, LogIndex: 1},
		RoundID:   1,
		Request:   request,
		Response:  nil,
		Fee:       fee,
		Status:    attestation.Waiting,
		Consensus: false,
		Hash:      hash,
	}

	url := "http://localhost:4500/eth/EVMTransaction/verifyFDC"
	err := verificationServer.ResolveAttestationRequest(&attestation, url, "12345")

	if err != nil {
		fmt.Println("Error resolving attestation request")
	} else {
		fmt.Println("Done")
		printStructFields(attestation)
	}

}

func printStructFields(s interface{}) {
	// Get the type of the struct
	t := reflect.TypeOf(s)

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("%s: %v\n", field.Name, reflect.ValueOf(s).Field(i))
	}
}
