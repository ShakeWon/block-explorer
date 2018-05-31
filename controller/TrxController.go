package controller

import (
    "github.com/kataras/iris"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/model/vo/response"
)

type TrxController struct {
    Ctx iris.Context
    service.TransactionService
}


func (t *TrxController) GetPage() {

    pageIndex, error := t.Ctx.URLParamInt("pageIndex")

    if error != nil {
        t.Ctx.Application().Logger().Error(error)
        t.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    pageSize, error := t.Ctx.URLParamInt("pageSize")
    if error != nil {
        t.Ctx.Application().Logger().Error(error)
        t.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    height:=t.Ctx.URLParam("height")

    hash:= t.Ctx.URLParam("hash")

    count, error := t.Count(height,hash)
    if error != nil {
        t.Ctx.Application().Logger().Error(error)
        t.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
        )
        return
    }

    var data []response.Transaction
    if count > 0 {
        data, error = t.TransactionService.Page(pageIndex, pageSize,height,hash)
        if error != nil {
            t.Ctx.Application().Logger().Error(error)
            t.Ctx.JSON(response.BaseResponse{Success: false, Error: error,},
            )
            return
        }

    }

    t.Ctx.JSON(
        response.BaseResponse{
            Success: true, Error: error,
            Data: response.PageTrxResponse{
                Total: count,
                Data:  data,
            },
        },
    )

}



