package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/sirupsen/logrus"
)

func main() {
	InitLogger()

	config, err := GetConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Minute)

	blockHeigth := new(int32)

	// main loop
	for {
		CollectData(config, blockHeigth)
		<-ticker.C
	}
}

func CollectData(config *GlobalConfig, blockHeigth *int32) {
	logrus.Info("start collecting data")
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         config.Node.Host,
		User:         config.Node.User,
		Pass:         config.Node.Password,
	}, nil)
	if err != nil {
		logrus.Errorf("cannot create btc client: %v", err)
		return
	}

	blockChainInfo, err := client.GetBlockChainInfo()
	if err != nil {
		logrus.Errorf("cannot get blockchain info: %v", err)
		return
	}

	logrus.Infof("head block: %d", blockChainInfo.Blocks)

	if blockChainInfo.Blocks == *blockHeigth {
		logrus.Warn("no new blocks. skipped")
		return
	}

	*blockHeigth = blockChainInfo.Blocks

	blockStatsReq := &[]string{"height", "time", "maxfeerate", "avgfeerate", "minfeerate", "maxfee", "avgfee", "minfee"}
	blockStatsRes, err := client.GetBlockStats(blockChainInfo.Blocks, blockStatsReq)
	if err != nil {
		logrus.Errorf("cannot get block stats: %v", err)
		return
	}

	nodeFee, err := CollectNodeFee(&config.Node, client)
	if err != nil {
		logrus.Errorf("cannot collect node fee: %v", err)
	}

	apiFee, err := CollectApiFee(&config.API)
	if err != nil {
		logrus.Errorf("cannot collect node fee: %v", err)
	}

	currentBlockReport := &FeeReport{
		BlockNumber: blockStatsRes.Height,
		BlockTime:   time.Unix(blockStatsRes.Time, 0),
		MaxFeeRate:  blockStatsRes.MaxFeeRate,
		AvgFeeRate:  blockStatsRes.AverageFeeRate,
		MinFeeRate:  blockStatsRes.MinFeeRate,
		MaxFee:      blockStatsRes.MaxFee,
		AvgFee:      blockStatsRes.AverageFee,
		MinFee:      blockStatsRes.MinFee,
		Suggestions: SuggestionsInfo{
			SuggestedBlock: blockStatsRes.Height + 1,
			RequestTime:    time.Now().UTC(),
			API:            apiFee,
			Node:           nodeFee,
		},
	}

	err = WriteReport(currentBlockReport)
	if err != nil {
		logrus.Errorf("cannot write report to file: %v", err)
		return
	}
}

func CollectNodeFee(config *NodeFeeConfig, client *rpcclient.Client) (int64, error) {
	fee, err := client.EstimateSmartFee(6, &btcjson.EstimateModeConservative)
	if err != nil {
		logrus.Errorf("cannot estimate fee: %v", err)
	}

	// s/kb -> s/b
	return int64(*fee.FeeRate * SatoshiInBTC / ByteInKilobyte), nil
}

func CollectApiFee(config *APIFeeConfig) (*APIFee, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, config.Url, nil)
	if err != nil {
		err = fmt.Errorf("got connection error %s", err.Error())
		logrus.Error(err)
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("cannot make request %s", err.Error())
		logrus.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("cannot read response body: %s", err.Error())
		logrus.Error(err)
		return nil, err
	}

	data := &APIFee{}
	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		if err != nil {
			err = fmt.Errorf("cannot unmarshal response body: %s", err.Error())
			logrus.Error(err)
			return nil, err
		}
	}

	return data, nil
}
