package xormimpl

import (
    "github.com/go-xorm/xorm"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/shakewon/block-explorer/repository"
    "strconv"
)

type XormTransactionRepoImpl struct {
    *xorm.Engine
    repository.TransactionsRepo
}

func (x *XormTransactionRepoImpl) Count(height, hash string) (int64, error) {
    transaction := &po.Transaction{}

    if len(height) > 0 {
        if h, err := strconv.Atoi(height); err != nil {
            return 0, err
        } else {
            transaction.Height = int64(h)
        }
    }
    if len(hash) > 0 {
        transaction.Hash = hash
    }
    total, err := x.Engine.Count(transaction)
    return total, err
}

func (x *XormTransactionRepoImpl) Page(index, pageSize int,height, hash string) ([]po.Transaction, error) {
    start := 0
    if index >= 1 {
        start = (index - 1) * pageSize
    }

    transaction := &po.Transaction{}

    if len(height) > 0 {
        if h, err := strconv.Atoi(height); err != nil {
            return nil, err
        } else {
            transaction.Height = int64(h)
        }
    }
    if len(hash) > 0 {
        transaction.Hash = hash
    }

    var resp = make([]po.Transaction, 0)
    error := x.Engine.Desc("Height").Limit(pageSize, start).Find(&resp,transaction)
    return resp, error
}

func (x *XormTransactionRepoImpl) Save(trxs []po.Transaction) error {
    if len(trxs) > 0 {
        _, error := x.Engine.Insert(trxs)
        return error
    } else {
        return nil
    }
}

func (x *XormTransactionRepoImpl) Search(hash string) ([]po.Transaction, error) {
    var resp = make([]po.Transaction, 0)
    error := x.Engine.Where("Hash=?", hash).Or("From_addr=?", hash).Or("To_Addr=?", hash).Desc("Height").Limit(25,0).Find(&resp)
    return resp, error
}
