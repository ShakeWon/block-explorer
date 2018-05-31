package controller

import (
    "github.com/kataras/iris"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/model/vo/response"
    "github.com/kataras/iris/core/errors"
)

var (
    BLOCK       = "BLOCK"
    TRANSACTION = "TRANSACTION"
)

type SearchController struct {
    Ctx iris.Context
    BlockSerivce service.BlockService
    TransactionService service.TransactionService
}

func (s *SearchController) GetSearch() {
    keywords := s.Ctx.URLParam("keywords")
    if len(keywords) == 0 {
        s.Ctx.JSON(response.BaseResponse{Success: false, Error: errors.New("keywords can not be null")})
        return
    }
}
