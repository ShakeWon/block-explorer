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
    Ctx                iris.Context
    BlockSerivce       service.BlockService
    TransactionService service.TransactionService
}

func (s *SearchController) GetSearch() {

    keywords := s.Ctx.URLParam("keywords")
    if len(keywords) == 0 {
        s.Ctx.JSON(response.BaseResponse{Success: false, Error: errors.New("keywords can not be null")})
        return
    }

    if _, error := strconv.Atoi(keywords); error == nil {
        //block search
        if data, error := s.BlockSerivce.Page(1, 1, keywords, ""); error != nil {
            s.Ctx.JSON(response.BaseResponse{Success: false, Error: error})
            return
        } else {
            s.Ctx.JSON(response.BaseResponse{Success: true, Data: response.SearchResponse{Type: BLOCK, Data: data}})
        }
    } else {

        // trx search

        if data, error := s.TransactionService.Search(keywords);error!=nil{
            s.Ctx.JSON(response.BaseResponse{Success: false, Error: error})
            return
        }else {
            s.Ctx.JSON(response.BaseResponse{Success: true, Data: response.SearchResponse{Type: TRANSACTION, Data: data}})
        }

    }
}
