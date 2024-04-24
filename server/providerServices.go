package server

import "fmt"

func submit1Handler(round string, address string) PDPResponse {
	fmt.Println("Submit1Handler")
	fmt.Printf("round: %s\n", round)
	fmt.Printf("address: %s\n", address)
	return PDPResponse{Status: OK, Data: round + address}

}
