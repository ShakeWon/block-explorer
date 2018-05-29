package xormimpl

import (
    "github.com/go-xorm/xorm"
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/model/po"
)

type XormBlockRepoImpl struct {
    xorm.Engine
    repository.BlockRepo
}

func (x *XormBlockRepoImpl) Count() (int64, error) {
    count, error := x.Engine.Count(po.Block{})
    return count, error
}

func (x *XormBlockRepoImpl) Page(index, pageSize int) ([]po.Block, error) {
    start := 0
    if index >= 1 {
        start = (index - 1) * pageSize
    }
    var resp []po.Block
    error := x.Engine.OrderBy("Height").Desc().Limit(start, pageSize).Find(resp)
    return resp, error
}

func (x *XormBlockRepoImpl) Query(height int) (po.Block, error) {
    block := po.Block{Height: int64(height)}
    _, error := x.Engine.Get(&block)
    return block, error
}

func (x *XormBlockRepoImpl) save(blocks []po.Block) error {
    if len(blocks) > 0 {
        _, error := x.Engine.Insert(blocks)
        return error
    } else {
        return nil
    }
}
