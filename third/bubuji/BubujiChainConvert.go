package bubuji

import (
	"github.com/astaxie/beego"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"math/big"
	"strconv"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
    "github.com/shakewon/block-explorer/third"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/Sirupsen/logrus"
)

type BubujiChainConvert struct {
	URL     string
	ChainId string
}

var EthSigner = HomesteadSigner{}

type Metas struct {
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockMeta struct {
	Hash   string  `json:"hash"`   // The block hash
	Header *Header `json:"header"` // The block's Header
}
type Header struct {
	ChainID        string   `json:"chain_id"`
	Height         int      `json:"height"`
	Time           JsonTime `json:"time"`
	NumTxs         int      `json:"num_txs"`          // XXX: Can we get rid of this?
	LastCommitHash string   `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       string   `json:"data_hash"`        // transactions
	ValidatorsHash string   `json:"validators_hash"`  // validators for the current block
	AppHash        string   `json:"app_hash"`         // state after txs from the previous block
	ReceiptsHash   string   `json:"recepits_hash"`    // recepits_hash from previous block
	LastBlockID    BlockID  `json:"last_block_id"`
}

type BlockID struct {
	Hash string `json:"hash"`
}

func (c *BubujiChainConvert) Init() {

}

func (c *BubujiChainConvert) Height()  int{
	return third.GetStatus(c.URL,c.ChainId).LatestBlockHeight
}

func (c *BubujiChainConvert) Block(h int) * third.BlockRepo  {
    logrus.Info("current block height : ", h)
	url := fmt.Sprintf("%s/blockchain?minHeight=%d&maxHeight=%d&chainid=\"%s\"", c.URL, h, h, c.ChainId)

	bytez, err := third.GetHTTPResp(url)
	if err != nil {
		beego.Info(err)
		return &third.BlockRepo{}
	}
	var metas Metas
	err = json.Unmarshal(bytez, &metas)
	if err != nil {
		beego.Error(err.Error())
		return &third.BlockRepo{}
	}
	//save block
	br := &third.BlockRepo{
		Blocks: []po.Block{},
		Txs:    []po.Transaction{},
	}
	for _, o := range metas.BlockMetas {
		rb := po.Block{
			Hash:           strings.ToLower(o.Hash),
			ChainId:        o.Header.ChainID,
			Height:         int64(o.Header.Height),
			Time:           time.Time(o.Header.Time),
			NumTxs:         int64(o.Header.NumTxs),
			LastCommitHash: strings.ToLower(o.Header.LastCommitHash),
			DataHash:       strings.ToLower(o.Header.DataHash),
			ValidatorsHash: strings.ToLower(o.Header.ValidatorsHash),
			AppHash:        strings.ToLower(o.Header.AppHash),
		}
		br.Blocks = append(br.Blocks, rb)
		if o.Header.NumTxs > 0 {
			//		blockHash := common.ToHex(o.Header.Hash())
			resultBlock, err := GetBlock(o.Header.Height, c.URL, c.ChainId)
			if err != nil {
                logrus.Error("GetBlock(height:%d) Error :%v\n", o.Header.Height, err)
				return &third.BlockRepo{}
			} else {
				for _, v := range resultBlock.Block.Data.Txs {
					tx := new(Transaction)
					err := rlp.DecodeBytes(common.FromHex(v), tx)
					if err != nil {
						logrus.Error("Decode Transaction tx bytes: [%v], error : %v\n", v, err)
						return &third.BlockRepo{}
					}

					from, err := Sender(EthSigner, tx)
					if err != nil {
                        logrus.Error("Error : Get Transaction From error: %v\n", err)
						return &third.BlockRepo{}
					}
					rtx := po.Transaction{
						Payload: string(tx.Data()),
						Hash:    tx.Hash().Hex(),
						FromAddr:    from.Hex(),
						Amount:  new(big.Int).Set(tx.Value()).String(),
						Nonce:   string(tx.Nonce()),
						Size:    int64(tx.Size()),
						Block:   rb.Hash,
						Time:    rb.Time,
						Height:  int64(o.Header.Height),
					}

					to := tx.To()
					if to == nil {
						// contract create
						contractAddr := crypto.CreateAddress(from, tx.Nonce())
						rtx.Contract = contractAddr.Hex()
					} else {
						rtx.ToAddr = to.Hex()
					}

					br.Txs = append(br.Txs, rtx)

				}
			}
		}
	}

	return br

}

type Block struct {
	*Header `json:"header"`
	*Data   `json:"data"`
}

type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

type Data struct {
	Txs []string `json:"txs"`

	// Volatile
	hash []byte
}

func GetBlock(height int, httpUrl, chanid string) (result ResultBlock, err error) {
	url := fmt.Sprintf("%s/block?height=%d&chainid=\"%s\"", httpUrl, height, chanid)
	bytez, errB := third.GetHTTPResp(url)
	if errB != nil {
		err = errB
		return
	}
	err = json.Unmarshal(bytez, &result)
	return
}

type JsonTime time.Time

func (j *JsonTime) UnmarshalJSON(data []byte) error {
	nano, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*j = JsonTime(time.Unix(0, nano))
	return nil
}
