package repository

import "github.com/shakewon/block-explorer/model/po"

type BlockRepo interface {
    /**
     * total number of blocks saved in db
     */
    Count() (int64,error)

    /**
     *  分页查询
     *  @index 页码
     *  @pageSize 页大小
     */
    Page(index, pageSize int) ([]po.Block,error)

    /**
     *  按照高度查询
     */
    Query(height int) (*po.Block,error)

    Save(blocks []po.Block)error
<<<<<<< HEAD
=======

    Height() (int64,error)
>>>>>>> c6fb292d2362ce2a42572dccaa159c1dd6fd551f
}
