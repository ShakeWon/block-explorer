package repository

import (
    "github.com/go-xorm/xorm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/shakewon/block-explorer/model/po"
)

func InitDataSouce() (*xorm.Engine, error) {
    orm, err := xorm.NewEngine("sqlite3", "./block.db?cache=shared&_busy_timeout=30000")
    if err != nil {
        return orm, err
    }
    orm.SetMaxOpenConns(1)
    err = orm.Sync2(new(po.Transaction), new(po.Block))
    return orm, err
}
