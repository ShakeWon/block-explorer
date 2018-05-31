package bubuji

import (
    "fmt"
    "encoding/json"
    "strings"
    "time"
    "strconv"
    "github.com/ethereum/go-ethereum/common"
    "github.com/shakewon/block-explorer/third"
    "github.com/shakewon/block-explorer/model/po"
    "bytes"
    "encoding/gob"
    "encoding/hex"
    "gitlab.zhonganonline.com/ann/prover/src/chain/app"
    "github.com/annchain/angine/types"
    "github.com/kataras/golog"
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
    Hash   string  `json:"Hash"`   // The block hash
    Header *Header `json:"Header"` // The block's Header
}
type BlockID struct {
    Hash []byte `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
}
type Header struct {
    ChainID            string   `protobuf:"bytes,1,opt,name=ChainID,proto3" json:"ChainID,omitempty"`
    Height             int      `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
    Time               int64    `protobuf:"varint,3,opt,name=Time,proto3" json:"Time,omitempty"`
    NumTxs             int64    `protobuf:"varint,4,opt,name=NumTxs,proto3" json:"NumTxs,omitempty"`
    Maker              []byte   `protobuf:"bytes,5,opt,name=Maker,proto3" json:"Maker,omitempty"`
    LastBlockID        *BlockID `protobuf:"bytes,6,opt,name=LastBlockID" json:"LastBlockID,omitempty"`
    LastCommitHash     string   `protobuf:"bytes,7,opt,name=LastCommitHash,proto3" json:"LastCommitHash,omitempty"`
    DataHash           string   `protobuf:"bytes,8,opt,name=DataHash,proto3" json:"DataHash,omitempty"`
    ValidatorsHash     string   `protobuf:"bytes,9,opt,name=ValidatorsHash,proto3" json:"ValidatorsHash,omitempty"`
    AppHash            string   `protobuf:"bytes,10,opt,name=AppHash,proto3" json:"AppHash,omitempty"`
    ReceiptsHash       string   `protobuf:"bytes,11,opt,name=ReceiptsHash,proto3" json:"ReceiptsHash,omitempty"`
    LastNonEmptyHeight int64    `protobuf:"varint,12,opt,name=LastNonEmptyHeight,proto3" json:"LastNonEmptyHeight,omitempty"`
    CoinBase           string   `json:"CoinBase,omitempty"`
    BlockRewards       uint64   `protobuf:"varint,14,opt,name=BlockRewards,proto3" json:"BlockRewards,omitempty"`
}

type Block struct {
    *Header `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
    *Data   `protobuf:"bytes,2,opt,name=Data" json:"Data,omitempty"`
}

type ResultBlock struct {
    BlockMeta *BlockMeta `json:"block_meta"`
    Block     *Block     `json:"block"`
}
type Data struct {
    Txs   []string `json:"Txs,omitempty"`
    ExTxs []string `json:"ExTxs,omitempty"`
}

func (c *BubujiChainConvert) Init() {

}

func (c *BubujiChainConvert) Height() int {
    return third.GetStatus(c.URL, c.ChainId).LatestBlockHeight
}

func (c *BubujiChainConvert) Block(h int) *third.BlockRepo {
    golog.Info("current block height : ", h)

    block, errG := GetBlock(h, c.URL, c.ChainId)
    if errG != nil {
        golog.Error(errG)
        return nil
    }
    //save block
    br := &third.BlockRepo{
        Blocks: []po.Block{},
        Txs:    []po.Transaction{},
    }

    rb := po.Block{
        Hash:           strings.ToLower(block.BlockMeta.Hash),
        ChainId:        block.Block.ChainID,
        Height:         int64(block.Block.Height),
        Time:           time.Unix(0, block.BlockMeta.Header.Time),
        LastCommitHash: strings.ToLower(block.Block.LastCommitHash),
        DataHash:       strings.ToLower(block.Block.DataHash),
        ValidatorsHash: strings.ToLower(block.Block.ValidatorsHash),
        AppHash:        strings.ToLower(block.Block.AppHash),
        Reward:         int64(block.Block.BlockRewards),
        CoinBase:       block.Block.CoinBase,
    }
    if len(block.Block.Data.Txs) == 0 && len(block.Block.Data.ExTxs) == 0 {
        br.Blocks = append(br.Blocks, rb)
    } else {
        TxsData := make(map[string]int)
        //unique
        for _, v := range block.Block.Data.Txs {
            TxsData[v] = 0
        }
        rb.NumTxs = int64(len(TxsData))
        br.Blocks = append(br.Blocks, rb)
        for k, _ := range TxsData {

            rtx, errP := parseTransaction(&rb, k)
            if errP != nil {
                golog.Error(errP)
                return nil
            }
            br.Txs = append(br.Txs, rtx)
        }
    }

    return br

}

func parseTransaction(rb *po.Block, v string) (rtx po.Transaction, err error) {
    rtx = po.Transaction{
        Block:  rb.Hash,
        Time:   rb.Time,
        Height: rb.Height,
    }

    tag := string(common.FromHex(v)[:4])
    bytez := types.UnwrapTx(common.FromHex(v))
    golog.Info(tag)
    switch tag {

    case string(app.ProverConfirmTag):
        tx := &app.ConfirmTx{}
        decoder := gob.NewDecoder(bytes.NewReader(bytez))
        err = decoder.Decode(tx)
        if err != nil {
            return
        }

        buf := bytes.Buffer{}
        encoder := gob.NewEncoder(&buf)
        if err := encoder.Encode(tx); err != nil {
            return rtx, err
        }

        txBytes := types.WrapTx(app.ProverConfirmTag, buf.Bytes())

        hash := types.Tx(txBytes).Hash()
        rtx = po.Transaction{
            Hash:       hex.EncodeToString(hash),
            Amount:     strconv.FormatInt(tx.Value, 10),
            PayloadHex: tx.Payload,
            Block:      rb.Hash,
            Time:       rb.Time,
            Height:     rb.Height,
        }

    default:
        golog.Warnf("unknown tag: %s\n", tag)
        return
    }

    return
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
