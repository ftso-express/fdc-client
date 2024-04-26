package server

import "fmt"

func submit1Handler(round uint64, address string) PDPResponse {
	fmt.Println("Submit1Handler")
	fmt.Printf("round: %s\n", fmt.Sprint(round))
	fmt.Printf("address: %s\n", address)
	return PDPResponse{Status: OK, Data: fmt.Sprint(round) + address}
}
