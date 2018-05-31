package controller

import (
    "github.com/kataras/iris"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/model/vo/response"
    "github.com/kataras/iris/core/errors"
    "strconv"
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

    if height, error := strconv.Atoi(keywords); error == nil {
        // search BLOCK height
        if block, error := s.BlockSerivce.Query(height); error != nil {
            s.Ctx.JSON(response.BaseResponse{Success: false, Error: error})
        } else {
            s.Ctx.JSON(response.BaseResponse{Success: true, Data: response.SearchResponse{Type: BLOCK, Data: block}})
        }
    } else {
        // search by hash
        if block, error := s.BlockSerivce.Search(keywords); error != nil {
            s.Ctx.JSON(response.BaseResponse{Success: false, Error: error})
        } else if block != nil {
            s.Ctx.JSON(response.BaseResponse{Success: true, Data: response.SearchResponse{Type: BLOCK, Data: block}})
        } else {
            //search transaction
            if trxs, error := s.TransactionService.Search(keywords); error != nil {
                s.Ctx.JSON(response.BaseResponse{Success: false, Error: error})
            }else {
                s.Ctx.JSON(response.BaseResponse{Success: true, Data: response.SearchResponse{Type: TRANSACTION, Data: trxs}})
            }
        }
    }
}
