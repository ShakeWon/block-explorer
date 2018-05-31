package xormimpl

import (
    "github.com/go-xorm/xorm"
    "github.com/shakewon/block-explorer/model/po"
    "github.com/shakewon/block-explorer/repository"
    "strconv"
)

type XormBlockRepoImpl struct {
    *xorm.Engine
    repository.BlockRepo
}

func (x *XormBlockRepoImpl) Count(height, hash string) (int64, error) {
    block := po.Block{}
    if len(height) > 0 {
        if h, err := strconv.Atoi(height); err != nil {
            return 0, err
        } else {
            block.Height = int64(h)
        }
    }
    if len(hash) > 0 {
        block.Hash = hash
    }
    count, error := x.Engine.Count(block)
    return count, error
}

func (x *XormBlockRepoImpl) Page(index, pageSize int, height, hash string) ([]po.Block, error) {
    start := 0
    if index >= 1 {
        start = (index - 1) * pageSize
    }
    resp := make([]po.Block, 0)

    block := po.Block{}
    if len(height) > 0 {
        if h, err := strconv.Atoi(height); err != nil {
            return nil, err
        } else {
            block.Height = int64(h)
        }
    }
    if len(hash) > 0 {
        block.Hash = hash
    }

    error := x.Engine.Desc("Height").Limit(pageSize, start).Find(&resp,block)
    return resp, error
}

func (x *XormBlockRepoImpl) Save(blocks []po.Block) error {
    if len(blocks) > 0 {
        _, error := x.Engine.Insert(blocks)
        return error
    } else {
        return nil
    }
}



