package repository

import "github.com/shakewon/block-explorer/model/po"

type TransactionsRepo interface {
    Count() (int64, error)
    Page(index, pageSize int) ([]po.Transaction, error)
    Query(trxHash string) (*po.Transaction, error)
    Save(trxs []po.Transaction) error
}
