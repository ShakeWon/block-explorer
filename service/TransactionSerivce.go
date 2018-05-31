package service

import (
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/shakewon/block-explorer/model/vo/response"
)

type TransactionService struct {
    Ts repository.TransactionsRepo
}

func (t *TransactionService) Count(height,hash string) (int64, error) {
    return t.Ts.Count(height,hash)
}
func (t *TransactionService) Page(index, pageSize int,height,hash string) ([]response.Transaction, error) {
    data, error := t.Ts.Page(index, pageSize,height,hash)
    var resp []response.Transaction
    if len(data) > 0 {
        for _, e := range data {
            resp = append(resp, response.Transaction{
                Hash:       e.Hash,
                Payload:    e.Payload,
                PayloadHex: e.PayloadHex,
                FromAddr:   e.FromAddr,
                ToAddr:     e.ToAddr,
                Receipt:    e.Receipt,
                Amount:     e.Amount,
                Nonce:      e.Nonce,
                Gas:        e.Gas,
                Size:       e.Size,
                Block:      e.Block,
                Contract:   e.Block,
                Time:       e.Time,
                Height:     e.Height,
                TxType:     e.TxType,
                Fee:        e.Fee,
            })
        }
    }
    return resp, error
}

func (t *TransactionService) Save(trxs []po.Transaction) error {
    return t.Ts.Save(trxs)
}

func (t *TransactionService) Search(hash string)([]response.Transaction, error){
    data, error := t.Ts.Search(hash)
    var resp []response.Transaction
    if len(data) > 0 {
        for _, e := range data {
            resp = append(resp, response.Transaction{
                Hash:       e.Hash,
                Payload:    e.Payload,
                PayloadHex: e.PayloadHex,
                FromAddr:   e.FromAddr,
                ToAddr:     e.ToAddr,
                Receipt:    e.Receipt,
                Amount:     e.Amount,
                Nonce:      e.Nonce,
                Gas:        e.Gas,
                Size:       e.Size,
                Block:      e.Block,
                Contract:   e.Block,
                Time:       e.Time,
                Height:     e.Height,
                TxType:     e.TxType,
                Fee:        e.Fee,
            })
        }
    }
    return resp, error
}
