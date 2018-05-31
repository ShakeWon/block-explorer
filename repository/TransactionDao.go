package repository

import "github.com/shakewon/block-explorer/model/po"

type TransactionsRepo interface {
    Count(height,hash string) (int64, error)
    Page(index, pageSize int,height,hash string) ([]po.Transaction, error)
    Save(trxs []po.Transaction) error
    Search(hash string) ([]po.Transaction, error)
}
