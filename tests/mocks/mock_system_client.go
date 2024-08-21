package mocks

import (
	"encoding/json"
	"io"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	"local/fdc/server"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func MockSystemClient(systemConfig *config.System, userConfig *config.UserRaw, client *ethclient.Client, submitPrivateKey, submitSignaturePrivateKey string) {
	submitAddress, _ := PrivKeyToAddress(submitPrivateKey)
	submitSigAddress, _ := PrivKeyToAddress(submitSignaturePrivateKey)

	time.Sleep(time.Duration(timing.Chain.CollectDurationSec) * time.Second)
	for {
		now := time.Now()

		round, err := timing.RoundIdForTimestamp(uint64(now.Unix()))
		if err != nil {
			log.Fatal("Error: %s", err)
		}
		// two processes since the intervals are overlapping
		go SystemClientIteration(userConfig, submitAddress, submitSigAddress, round)
		go SystemClientIteration(userConfig, submitAddress, submitSigAddress, round+1)

		submit1Time := timing.ChooseStartTimestamp(round + 2)
		timer := time.NewTimer(time.Until(time.Unix(int64(submit1Time-5), 0)))
		<-timer.C
	}
}

func SystemClientIteration(userConfig *config.UserRaw, submitAddress, submitSigAddress common.Address, round uint64) {
	submit1Time := timing.ChooseStartTimestamp(round)
	submit2Time := submit1Time + timing.Chain.ChooseDurationSec
	submitSignature := submit2Time + timing.Chain.CommitDurationSec

	timer := time.NewTimer(time.Until(time.Unix(int64(submit1Time+2), 0)))
	<-timer.C

	rspData, err := MakeGetRequest("submit1", &userConfig.RestServer, round, submitAddress.Hex())
	if err != nil || rspData.Status != server.Ok {
		log.Error("error submit1 response ", rspData, err)
		return
	}
	log.Info("response submit1 ", rspData.Status)

	timer = time.NewTimer(time.Until(time.Unix(int64(submit2Time+6), 0)))
	<-timer.C
	rspData, err = MakeGetRequest("submit2", &userConfig.RestServer, round, submitAddress.Hex())
	if err != nil || rspData.Status != server.Ok {
		log.Error("error submit2 response ", rspData, err)
		return
	}
	log.Info("response submit2 ", rspData.Status)

	timer = time.NewTimer(time.Until(time.Unix(int64(submitSignature+5), 0)))
	<-timer.C
	rspData, err = MakeGetRequest("submitSignatures", &userConfig.RestServer, round, submitSigAddress.Hex())
	if err != nil || rspData.Status != server.Ok {
		log.Error("error submit2 response ", rspData, err)
		return
	}
	log.Info("response submitSignatures ", rspData.Status)
}

func MakeGetRequest(
	apiName string, cfg *config.RestServer, votingRoundID uint64, submitAddress string,
) (*server.PDPResponse, error) {
	p, err := url.JoinPath(
		cfg.FSPSubpath,
		apiName,
		strconv.FormatUint(votingRoundID, 10),
		submitAddress,
	)
	if err != nil {
		return nil, err
	}
	log.Info("making request to ", p)

	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   p,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-KEY", cfg.ApiKeys[0])

	var client http.Client
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %s", rsp.Status)
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	rspData := new(server.PDPResponse)
	if err = json.Unmarshal(body, rspData); err != nil {
		return nil, err
	}

	return rspData, nil
}
