package mocks

import (
	"encoding/json"
	"flare-common/database"
	"fmt"
	"io"
	"local/fdc/client/attestation"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func MockVerifier(t *testing.T, port int, response string, testLog database.Log) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		MockResponse(t, writer, request, response, testLog)
	})

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	fmt.Println("Mock verifier starting")
	err := server.ListenAndServe()
	require.NoError(t, err)
}

func MockResponse(t *testing.T, writer http.ResponseWriter, request *http.Request, response string, testLog database.Log) {
	body, err := io.ReadAll(request.Body)
	require.NoError(t, err)

	var requestStruct attestation.AbiEncodedRequestBody
	err = json.Unmarshal(body, &requestStruct)
	require.NoError(t, err)
	require.Equal(t, "0x"+testLog.Data[192:len(testLog.Data)-1], requestStruct.AbiEncodedRequest[:len(requestStruct.AbiEncodedRequest)-1]) // todo: is it expected to be trimmed?

	responseStruct := attestation.AbiEncodedResponseBody{Status: "VALID", AbiEncodedResponse: response}
	responseBytes, err := json.Marshal(responseStruct)
	require.NoError(t, err)

	_, err = writer.Write(responseBytes)
	require.NoError(t, err)
}
