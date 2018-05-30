package third

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strings"
	"net/http"
	"errors"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/Sirupsen/logrus"
)

type HTTPResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
}

type BlockRepo struct {
    Blocks []po.Block
    Txs    []po.Transaction
}

type DataConverter interface {
	Block(height int) *BlockRepo
	Init()
	Height() int
}

type Status struct {
	NodeInfo          *NodeInfo `json:"node_info"`
	LatestBlockHeight int       `json:"latest_block_height"`
}

type NodeInfo struct {
	NetWork string `json:"network"`
}

func GetStatus(httpAddr, chainID string) (status Status) {

	url := fmt.Sprintf("%s/status?chainid=\"%s\"", httpAddr, chainID)
	bytez, err := GetHTTPResp(url)
	if err != nil {
		logrus.Fatalf("GetHTTPResp failed: %s", err.Error())
	}

	err = json.Unmarshal(bytez, &status)
	if err != nil {
        logrus.Fatalf("json.Unmarshal(Status) failed: %s", err.Error())
	}
	return
}

func GetHTTPResp(url string) (bytez []byte, err error) {

	resp, errR := http.Get(url)
	if errR != nil {
		err = errR
		return
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var hr HTTPResponse
	err = json.Unmarshal(bytez, &hr)
	if err != nil {
		return
	}
	if hr.Result == nil {
		err = errors.New(fmt.Sprintf("json.Unmarshal (%s)HTTPResponse wrong ,maybe you need config 'chain_id'", url))
		return
	}
	bytez, err = hr.Result.MarshalJSON()
	if err != nil {
		return
	}
	str := string(bytez)
	i := strings.Index(str, ",")
	bytez = bytez[i+1 : len(bytez)-1]
	return
}
