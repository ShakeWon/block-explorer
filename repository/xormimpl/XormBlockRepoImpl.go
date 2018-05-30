package xormimpl

import (
    "github.com/go-xorm/xorm"
    "github.com/shakewon/block-explorer/model/po"
)

type XormBlockRepoImpl struct {
    *xorm.Engine
<<<<<<< HEAD
    repository.BlockRepo
=======
>>>>>>> c6fb292d2362ce2a42572dccaa159c1dd6fd551f
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
<<<<<<< HEAD
    var resp []po.Block
    error := x.Engine.OrderBy("Height").Desc("Height").Limit(start, pageSize).Find(resp)
=======
    resp:= make([]po.Block,0)
    error := x.Engine.Desc("Height").Limit(pageSize, start).Find(&resp)
>>>>>>> c6fb292d2362ce2a42572dccaa159c1dd6fd551f
    return resp, error
}

func (x *XormBlockRepoImpl) Query(height int) (*po.Block, error) {
    block := &po.Block{Height: int64(height)}
    exist, error := x.Engine.Get(block)
    if exist {
        return block, error
    } else {
        return nil, error
    }

}

func (x *XormBlockRepoImpl) Save(blocks []po.Block) error {
    if len(blocks) > 0 {
        _, error := x.Engine.Insert(blocks)
        return error
    } else {
        return nil
    }
}

func (x *XormBlockRepoImpl) Height() (int64,error) {
    var resp = make([]po.Block,0)
    if error := x.Engine.Desc("Height").Limit(1, 0).Find(&resp); error == nil && len(resp) > 0 {
        return resp[0].Height, error
    } else {
        return int64(0), error
    }

}
