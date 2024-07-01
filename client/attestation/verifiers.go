package attestation

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Function used to resolve attestation requests

type AbiEncodedRequestBody struct {
	AbiEncodedRequest string `json:"abiEncodedRequest"`
}

type AbiEncodedResponseBody struct {
	Status             string `json:"status"`
	AbiEncodedResponse string `json:"abiEncodedResponse"`
}

type VerifierCredentials struct {
	Url    string
	apiKey string
}

// ResolveAttestationRequest sends the attestation request to the verifier server with verifierCred and stores the response.
// Returns true if the response is "VALID" and false otherwise.
func ResolveAttestationRequest(att *Attestation) (bool, error) {
	client := &http.Client{}
	requestBytes := att.Request

	encoded := hex.EncodeToString(requestBytes)
	payload := AbiEncodedRequestBody{AbiEncodedRequest: "0x" + encoded}

	encodedBody, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	request, err := http.NewRequest("POST", att.Credentials.Url, bytes.NewBuffer(encodedBody))
	if err != nil {
		return false, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", att.Credentials.apiKey)
	resp, err := client.Do(request)

	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("request responded with code %d", resp.StatusCode)
	}

	// close response body after function ends
	defer resp.Body.Close()

	responseBody := AbiEncodedResponseBody{}

	decoder := json.NewDecoder(resp.Body)

	decoder.DisallowUnknownFields()

	err = decoder.Decode(&responseBody)

	if err != nil {
		return false, err
	}

	if responseBody.Status != "VALID" {
		return false, nil
	}

	responseBytes, err := hex.DecodeString(strings.TrimPrefix(responseBody.AbiEncodedResponse, "0x"))
	if err != nil {
		return false, err
	}

	att.Response = responseBytes

	return true, nil
}
