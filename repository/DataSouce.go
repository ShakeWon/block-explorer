package repository

import (
    "github.com/go-xorm/xorm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/shakewon/block-explorer/model"
    "strings"
    "fmt"
)

var (
    sqlite = "sqlite"
    mysql = "mysql"
    mongodb = "mongodb"
)

func InitDataSouce(config model.AppConfig) (*xorm.Engine, error) {

    if strings.EqualFold(sqlite,config.Store.DriverName){
        orm, err := xorm.NewEngine("sqlite3", "./block.db?cache=shared&_busy_timeout=30000")
        if err != nil {
            return orm, err
        }
        orm.SetMaxOpenConns(1)
        err = orm.Sync2(new(po.Transaction), new(po.Block))
        return orm, err
    }else {
        panic(fmt.Sprintf("%v driver is not supported ",config.Store.DriverName))
    }
}
