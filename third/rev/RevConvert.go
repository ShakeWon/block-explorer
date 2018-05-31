package rev

import (
	"math/big"
	"github.com/astaxie/beego"
	"fmt"
	"encoding/json"
	"time"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
    "github.com/shakewon/block-explorer/third"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/kataras/golog"
)

type ResultBlockchainInfo struct {
	LastHeight int64        `json:"last_height"`
	BlockMetas []*BlockMeta `json:"block_metas"`
}

type Vote struct {
	ValidatorAddress []byte           `json:"validator_address"`
	ValidatorIndex   int              `json:"validator_index"`
	Height           int64            `json:"height"`
	Round            int64            `json:"round"`
	Type             byte             `json:"type"`
	BlockID          BlockID          `json:"block_id"` // zero if vote is nil.
}

type Data struct {

	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs   Txs `json:"txs"`
	ExTxs Txs `json:"extxs"` // this is for all other txs which won't be delivered to app

	// Volatile
	hash []byte
}

type Txs []Tx

type Tx []byte

type BlockMeta struct {
	Hash        string        `json:"hash"`         // The block hash
	Header      *Header       `json:"header"`       // The block's Header
	PartsHeader PartSetHeader `json:"parts_header"` // The PartSetHeader, for transfer
}

type Header struct {
	ChainID        string    `json:"chain_id"`
	Height         int64     `json:"height"`
	Time           time.Time `json:"time"`
	NumTxs         int64     `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID    BlockID   `json:"last_block_id"`
	LastCommitHash string    `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       string    `json:"data_hash"`        // transactions
	ValidatorsHash string    `json:"validators_hash"`  // validators for the current block
	AppHash        string    `json:"app_hash"`         // state after txs from the previous block
	ReceiptsHash   string    `json:"recepits_hash"`    // recepits_hash from previous block
}

type BlockID struct {
	Hash        []byte        `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}

type ResultBlock struct {
	BlockMeta *BlockMeta `json:"block_meta"`
	Block     *Block     `json:"block"`
}

type Block struct {
	*Header    `json:"header"`
	*Data      `json:"data"`
}

type PartSetHeader struct {
	Total int    `json:"total"`
	Hash  []byte `json:"hash"`
}

var EthSigner = HomesteadSigner{}

type RevChainConvert struct {
	URL     string
	ChainId string
}

func (c *RevChainConvert) Init() {
	c.URL = beego.AppConfig.String("api_addr")
	c.ChainId = beego.AppConfig.String("chain_id")
	if c.ChainId == "" {
		c.ChainId = third.GetStatus(c.URL, "").NodeInfo.NetWork
	}
}

func (c *RevChainConvert) Height()  int{
	return third.GetStatus(c.URL,c.ChainId).LatestBlockHeight
}

func (c *RevChainConvert) Block(h int) *third.BlockRepo {
    golog.Info("current block height : ", h)
	url := fmt.Sprintf("%s/blockchain?minHeight=%d&maxHeight=%d&chainid=\"%s\"", c.URL, h+1, h+5, c.ChainId)

	bytez, err := third.GetHTTPResp(url)
	if err != nil {
		beego.Info(err)
		return &third.BlockRepo{}
	}
	var metas ResultBlockchainInfo
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
			Hash:           o.Hash,
			ChainId:        o.Header.ChainID,
			Height:         int64(o.Header.Height),
			Time:           time.Time(o.Header.Time),
			NumTxs:         int64(o.Header.NumTxs),
			LastCommitHash: o.Header.LastCommitHash,
			DataHash:       o.Header.DataHash,
			ValidatorsHash: o.Header.ValidatorsHash,
			AppHash:        o.Header.AppHash,
		}
		br.Blocks = append(br.Blocks, rb)
		if o.Header.NumTxs > 0 {
			//		blockHash := common.ToHex(o.Header.Hash())
			resultBlock, err := GetBlock(int(o.Header.Height), c.URL, c.ChainId)
			if err != nil {
				golog.Error("GetBlock(height:%d) Error :%v\n", o.Header.Height, err)
				return &third.BlockRepo{}
			} else {
				for _, v := range resultBlock.Block.Data.Txs {
					tx := new(Transaction)
					err := rlp.DecodeBytes(common.FromHex(string(v)), tx)
					if err != nil {
						beego.Error("Decode Transaction tx bytes: [%v], error : %v\n", v, err)
						return &third.BlockRepo{}
					}

					from, err := Sender(EthSigner, tx)
					if err != nil {
						beego.Error("Error : Get Transaction From error: %v\n", err)
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
