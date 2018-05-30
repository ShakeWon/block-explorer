package xormimpl

import (
    "github.com/go-xorm/xorm"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/shakewon/block-explorer/repository"
)

type XormTransactionRepoImpl struct {
    *xorm.Engine
    repository.TransactionsRepo
}

func (x *XormTransactionRepoImpl) Count() (int64, error) {
    total, err := x.Engine.Count(&po.Transaction{})
    return total, err
}

func (x *XormTransactionRepoImpl) Page(index, pageSize int) ([]po.Transaction, error) {
    start := 0
    if index >= 1 {
        start = (index - 1) * pageSize
    }
    var resp  = make([]po.Transaction,0)
    error := x.Engine.Desc("Height").Limit(pageSize, start).Find(&resp)
    return resp, error
}

func (x *XormTransactionRepoImpl) Query(trxHash string) (*po.Transaction, error) {
    trx := &po.Transaction{Hash: trxHash}
    exists, err := x.Engine.Get(trx)
    if exists {
        return trx,err
    } else {
        return nil,err
    }
}

func (x *XormTransactionRepoImpl) Save(trxs []po.Transaction) error {
    if len(trxs) > 0 {
        _, error := x.Engine.Insert(trxs)
        return error
    } else {
        return nil
    }
}
