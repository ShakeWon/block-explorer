package main

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/middleware/logger"
    "github.com/kataras/iris/mvc"
    "github.com/shakewon/block-explorer/controller"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/repository/xormimpl"
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/sys"
    "github.com/shakewon/block-explorer/third/bubuji"
)

func main() {
    app := iris.New()

    customLogger := logger.New(logger.Config{
        Status:             true,
        IP:                 true,
        Method:             true,
        Path:               true,
        MessageContextKeys: []string{"logger_message"},
        MessageHeaderKeys:  []string{"User-Agent"},
    })

    app.Use(customLogger)
    app.Logger().SetLevel("debug")

    mvc.Configure(app.Party("/api"), basicMVC)

    app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
}

func basicMVC(app *mvc.Application) {
    engine, error := repository.InitDataSouce()
    if error != nil {
        panic(error)
    }
    trxController := new(controller.TrxController)

    transactionService := service.TransactionService{Ts: &xormimpl.XormTransactionRepoImpl{Engine: engine}}
    trxController.TransactionService = transactionService
    app.Party("/trx").
        Handle(trxController)

    blockController := new(controller.BlockController)

    blockService := service.BlockService{Bs: &xormimpl.XormBlockRepoImpl{Engine: engine}}
    blockController.BlockService = blockService
    app.Party("/block").Handle(blockController)
    convert := &bubuji.BubujiChainConvert{
       URL:     "http://rpc-bbjchain.zhonganinfo.com",
       ChainId: "prover",
    }
    convert.Init()
    job := sys.SysJob{
       BlockService:       blockService,
       TransactionService: transactionService,
       Convert:            convert,
    }
    go job.Start()
}
