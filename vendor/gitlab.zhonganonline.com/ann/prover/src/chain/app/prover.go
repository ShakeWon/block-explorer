package app

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	pbtypes "gitlab.zhonganonline.com/ann/angine/protos/types"
	agtypes "gitlab.zhonganonline.com/ann/angine/types"
	cmn "gitlab.zhonganonline.com/ann/ann-module/lib/go-common"
	"gitlab.zhonganonline.com/ann/ann-module/lib/go-crypto"
	"gitlab.zhonganonline.com/ann/ann-module/lib/go-merkle"
	"gitlab.zhonganonline.com/ann/prover/src/tools"
	ptypes "gitlab.zhonganonline.com/ann/prover/src/types"
)

const (
	APP_NAME = "proverstate"
)

var (
	ProverTag        = []byte{'p', 'r', 'o'}
	ProverConfirmTag = append(ProverTag, 0x01)
)

type (
	BlockState struct {
		Height   agtypes.INT
		Treehash []byte
	}

	ProverApp struct {
		agtypes.BaseApplication

		core ptypes.Core

		logger           *zap.Logger
		privkey          crypto.PrivKeyEd25519
		blockState       *BlockState
		PendingProverTxs map[string]*Tx

		attributes map[string]string

		config      *viper.Viper
		AngineHooks agtypes.Hooks
	}

	Tx struct {
		Value   int64
		Payload string
		Time    time.Time
	}

	ConfirmTx struct {
		Value   int64
		Payload string
		TxHash  []byte
		Time    time.Time

		// ECDSA.Pubkey
		X *big.Int
		Y *big.Int

		// Signature
		R *big.Int
		S *big.Int
	}

	LastBlockInfo struct {
		Height  agtypes.INT
		AppHash []byte
	}
)

func init() {}

// NewLastBlockInfo just a convience to generate an empty LastBlockInfo
func NewLastBlockInfo() *LastBlockInfo {
	return &LastBlockInfo{
		Height:  0,
		AppHash: make([]byte, 0),
	}
}

func NewProverApp(logger *zap.Logger, config *viper.Viper, privkey crypto.PrivKey) (*ProverApp, error) {
	datadir := config.GetString("db_dir")

	app := ProverApp{
		logger:  logger,
		config:  config,
		privkey: *privkey.(*crypto.PrivKeyEd25519),
		blockState: &BlockState{
			Treehash: make([]byte, 0),
			Height:   0,
		},

		PendingProverTxs: make(map[string]*Tx),
	}

	var err error
	if err = app.BaseApplication.InitBaseApplication(APP_NAME, datadir); err != nil {
		app.logger.Error("InitBaseApplication error", zap.Error(err))
		cmn.PanicCrisis(err)
	}

	app.AngineHooks = agtypes.Hooks{
		OnExecute: agtypes.NewHook(app.OnExecute),
		OnCommit:  agtypes.NewHook(app.OnCommit),
	}

	return &app, nil
}

func (app *ProverApp) GetAttributes() map[string]string {
	return app.attributes
}

func (app *ProverApp) GetAngineHooks() agtypes.Hooks {
	return app.AngineHooks
}

func (app *ProverApp) CompatibleWithAngine() {}

func (app *ProverApp) CheckTx(bs []byte) error {
	txBytes := agtypes.UnwrapTx(bs)

	if IsConfirmTx(bs) {
		tx := &ConfirmTx{}
		decoder := gob.NewDecoder(bytes.NewReader(txBytes))
		if err := decoder.Decode(tx); err != nil {
			return err
		}
		if err := app.checkConfirm(tx); err != nil {
			app.logger.Error("wrong decode checkTx")
			return errors.Wrap(err, "fail to checkTx")
		}
	}

	return nil
}

func (app *ProverApp) checkConfirm(tx *ConfirmTx) error {
	hash := tx.TxHash
	pubkey := ecdsa.PublicKey{Curve: secp256k1.S256(), X: tx.X, Y: tx.Y}
	e := tx.Verify(hash, &pubkey)
	if !e {
		return fmt.Errorf("signature failed")
	}
	return nil
}

func (ct *ConfirmTx) Verify(hash []byte, pubKey *ecdsa.PublicKey) bool {
	return ecdsa.Verify(pubKey, hash, ct.R, ct.S)
}

func (app *ProverApp) ExecuteTx(bs []byte) (hash []byte, err error) {
	txBytes := agtypes.UnwrapTx(bs)

	if IsConfirmTx(bs) {
		tx := &ConfirmTx{}
		decoder := gob.NewDecoder(bytes.NewReader(txBytes))
		if err := decoder.Decode(tx); err != nil {
			return nil, err
		}
		return bs, nil
	}
	return nil, nil
}

func (app *ProverApp) OnExecute(height, round agtypes.INT, block *agtypes.BlockCache) (interface{}, error) {
	var (
		res agtypes.ExecuteResult
		err error
	)

	for _, tx := range block.Data.Txs {
		if !bytes.Equal(tx[:4], ProverConfirmTag) {
			continue
		}

		if txHash, err := app.ExecuteTx(tx); err != nil {
			res.InvalidTxs = append(res.InvalidTxs, agtypes.ExecuteInvalidTx{Bytes: txHash, Error: err})
		} else {
			res.ValidTxs = append(res.ValidTxs, txHash)
		}
	}

	return res, err
}

func (app *ProverApp) OnCommit(height, round agtypes.INT, block *agtypes.BlockCache) (interface{}, error) {
	app.blockState.Height = height

	lastblock := LastBlockInfo{Height: app.blockState.Height, AppHash: app.AppHash()}
	app.Database.SetSync(lastblock.AppHash, app.blockState.ToBytes())
	app.SaveLastBlock(lastblock)

	return agtypes.CommitResult{AppHash: lastblock.AppHash}, nil
}

func (app *ProverApp) Query(query []byte) agtypes.Result {
	if len(query) == 0 {
		return agtypes.NewResultOK([]byte{}, "Empty query")
	}

	var res agtypes.Result

	qryRes, err := app.core.GetEngine().Query(query[0], query[1:])
	if err != nil {
		return agtypes.NewError(pbtypes.CodeType_InternalError, err.Error())
	}

	info, ok := qryRes.(*agtypes.TxExecutionResult)
	if !ok {
		return agtypes.NewError(pbtypes.CodeType_InternalError, err.Error())
	}
	res.Code = pbtypes.CodeType_OK
	res.Data, _ = info.ToBytes()
	return res
}

func (app *ProverApp) Info() (resInfo agtypes.ResultInfo) {
	lb := NewLastBlockInfo()
	if res, err := app.LoadLastBlock(lb); err == nil {
		lb = res.(*LastBlockInfo)
	}

	resInfo.LastBlockAppHash = lb.AppHash
	resInfo.LastBlockHeight = lb.Height
	return
}

func (app *ProverApp) Start() error {
	lastBlock := NewLastBlockInfo()
	if res, err := app.LoadLastBlock(lastBlock); err == nil {
		lastBlock = res.(*LastBlockInfo)
	}
	if lastBlock.AppHash == nil || len(lastBlock.AppHash) == 0 {
		return nil
	}
	if err := app.Load(lastBlock); err != nil {
		app.logger.Error("fail to load prover state", zap.Error(err))
		return err
	}

	return nil
}

func (app *ProverApp) Stop() {
	app.BaseApplication.Stop()
}

func (app *ProverApp) Load(lb *LastBlockInfo) error {
	state := &BlockState{}
	if lb == nil {
		return nil
	}

	bs := app.Database.Get(lb.AppHash)
	if err := state.FromBytes(bs); err != nil {
		app.logger.Error("fail to restore prover state", zap.Error(err))
		return err
	}
	app.blockState = state

	return nil
}

func (app *ProverApp) SetCore(core ptypes.Core) {
	app.core = core
}

func (app *ProverApp) BroadcastTx(tx []byte) error {
	return app.core.GetEngine().BroadcastTx(tx)
}

// GetNodePubKey gets our universal public key
func (app *ProverApp) GetNodePubKey() crypto.PubKey {
	return app.core.GetEngine().PrivValidator().GetPubKey()
}

func (app *ProverApp) AppHash() []byte {
	return merkle.SimpleHashFromBinary(app.blockState)
}

func (app *Tx) ToBytes() ([]byte, error) {
	return json.Marshal(app)
}

func (app *Tx) Hash() ([]byte, error) {
	txBytes, err := app.ToBytes()
	if err != nil {
		return nil, err
	}
	return tools.HashKeccak(txBytes)
}

func IsConfirmTx(tx []byte) bool {
	return bytes.Equal(tx[:4], ProverConfirmTag)
}

func (ct *ConfirmTx) FromBytes(bytes []byte) error {
	return json.Unmarshal(bytes, ct)
}

func (ct *ConfirmTx) ToBytes() ([]byte, error) {
	return json.Marshal(ct)
}

func (ct *ConfirmTx) Hash() []byte {
	return merkle.SimpleHashFromBinary(ct)
}

func (be *BlockState) FromBytes(bs []byte) error {
	bsReader := bytes.NewReader(bs)
	gdc := gob.NewDecoder(bsReader)
	return gdc.Decode(be)
}

func (be *BlockState) ToBytes() []byte {
	var bs []byte
	buf := bytes.NewBuffer(bs)
	gec := gob.NewEncoder(buf)
	if err := gec.Encode(be); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
