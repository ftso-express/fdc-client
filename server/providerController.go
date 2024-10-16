package server

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/go-flare-common/pkg/restserver"
	"github.com/flare-foundation/go-flare-common/pkg/storage"

	"gitlab.com/flarenetwork/fdc/fdc-client/client/round"
	"gitlab.com/flarenetwork/fdc/fdc-client/client/timing"
)

type FDCProtocolProviderController struct {
	rounds     *storage.Cyclic[uint32, *round.Round]
	protocolID uint8
}

type submitXParams struct {
	votingRoundID uint32
	submitAddress string
}

func newFDCProtocolProviderController(rounds *storage.Cyclic[uint32, *round.Round], protocolID uint8) *FDCProtocolProviderController {
	return &FDCProtocolProviderController{rounds: rounds, protocolID: protocolID}
}

const hexPrefix = "0x"

func validateEVMAddressString(address string) bool {
	address = strings.TrimPrefix(address, hexPrefix)
	dec, err := hex.DecodeString(address)
	if err != nil {
		return false
	}
	if len(dec) != 20 {
		return false
	}
	return err == nil
}

func validateSubmitXParams(params map[string]string) (submitXParams, error) {
	if _, ok := params["votingRoundID"]; !ok {
		return submitXParams{}, fmt.Errorf("missing votingRound param")
	}
	votingRoundID, err := strconv.ParseUint(params["votingRoundID"], 10, 32)
	if err != nil {
		return submitXParams{}, fmt.Errorf("votingRound param is not a number")
	}
	if _, ok := params["submitAddress"]; !ok {
		return submitXParams{}, fmt.Errorf("missing submitAddress param")
	}
	submitAddress := params["submitAddress"]
	if !validateEVMAddressString(submitAddress) {
		return submitXParams{}, fmt.Errorf("submitAddress param is not a valid EVM address")
	}
	return submitXParams{votingRoundID: uint32(votingRoundID), submitAddress: submitAddress}, nil
}

func submitXController(
	params map[string]string,
	_ interface{},
	_ interface{},
	service func(uint32, string) (string, bool, error),
	timeLock func(uint32) uint64,
) (PDPResponse, *restserver.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)

	if err != nil {
		logger.Error(err)
		return PDPResponse{}, restserver.BadParamsErrorHandler(err)
	}

	atTheEarliest := timeLock(pathParams.votingRoundID)
	now := uint64(time.Now().Unix())
	if atTheEarliest > now {
		return PDPResponse{}, restserver.ToEarlyErrorHandler(fmt.Errorf("to early %v before %d", atTheEarliest-now, atTheEarliest))
	}

	rsp, exists, err := service(pathParams.votingRoundID, pathParams.submitAddress)
	if err != nil {
		logger.Error(err)
		return PDPResponse{}, restserver.InternalServerErrorHandler(err)
	}
	if !exists {
		return PDPResponse{}, restserver.NotAvailableErrorHandler(fmt.Errorf("commit data for round id %d not available", pathParams.votingRoundID))
	}
	response := PDPResponse{Data: rsp, Status: Ok}
	return response, nil
}

func (controller *FDCProtocolProviderController) submit1Controller(
	params map[string]string,
	queryParams interface{},
	body interface{},
) (PDPResponse, *restserver.ErrorHandler) {
	return submitXController(params, queryParams, body, controller.submit1Service, timing.RoundStartTime)
}

func (controller *FDCProtocolProviderController) submit2Controller(
	params map[string]string,
	queryParams interface{},
	body interface{},
) (PDPResponse, *restserver.ErrorHandler) {
	return submitXController(params, queryParams, body, controller.submit2Service, timing.ChooseStartTimestamp)
}

func (controller *FDCProtocolProviderController) submitSignaturesController(
	params map[string]string,
	queryParams interface{},
	body interface{},
) (PDPResponse, *restserver.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)
	if err != nil {
		logger.Error(err)
		return PDPResponse{}, restserver.BadParamsErrorHandler(err)
	}
	response, exists, err := controller.submitSignaturesService(pathParams.votingRoundID, pathParams.submitAddress)
	if err != nil || !exists {
		return PDPResponse{}, restserver.NotAvailableErrorHandler(fmt.Errorf("round id %d not available", pathParams.votingRoundID))

	}
	return response, nil
}
