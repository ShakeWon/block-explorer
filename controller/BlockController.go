package controller

import (
    "github.com/kataras/iris"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/model/vo/response"
)

type BlockController struct {
    Ctx iris.Context
    service.BlockService
}

func (b *BlockController) GetPage() {

    pageIndex, error := b.Ctx.URLParamInt("pageIndex")

    if error != nil {
        b.Ctx.Application().Logger().Error(error)
        b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    pageSize, error := b.Ctx.URLParamInt("pageSize")
    if error != nil {
        b.Ctx.Application().Logger().Error(error)
        b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    count, error := b.Count()
    if error != nil {
        b.Ctx.Application().Logger().Error(error)
        b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    var data []response.Block
    if count > 0 {
        data, error = b.BlockService.Page(pageIndex, pageSize)
        if error != nil {
            b.Ctx.Application().Logger().Error(error)
            b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
            )
            return
        }

    }

    b.Ctx.JSON(
        response.BaseResponse{
            Success: true, Error: error,
            Data: response.PageBlockResponse{
                Total: count,
                Data:  data,
            },
        },
    )

}

func (b *BlockController) GetQuery() {
    height, error := b.Ctx.URLParamInt("height")
    if error != nil {
        b.Ctx.Application().Logger().Error(error)
        b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    data, error := b.BlockService.Query(height)
    if error != nil {
        b.Ctx.Application().Logger().Error(error)
        b.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }
    b.Ctx.JSON(
        response.BaseResponse{
            Success: true, Error: error,
            Data:    data,
        },
    )
}
