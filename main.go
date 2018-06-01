package main

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/mvc"
    "github.com/shakewon/block-explorer/controller"
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/repository/xormimpl"
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/model"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "github.com/kataras/golog"
    "runtime"
    "fmt"
    "github.com/iris-contrib/middleware/cors"
    "github.com/shakewon/block-explorer/third/bubuji"
    "github.com/shakewon/block-explorer/sys"
)

func main() {
    app := iris.New()

    app.Logger().Handle(func(l *golog.Log) bool {
        prefix := golog.GetTextForLevel(l.Level, true)
        pc, fn, line, _ := runtime.Caller(6)
        message := fmt.Sprintf("%s line %d (%s) (%s) %s: %s",
            prefix, line, runtime.FuncForPC(pc).Name(), fn, l.FormatTime(), l.Message)

        if l.NewLine {
            message += "\n"
        }

        fmt.Print(message)
        return true
    })

    app.Logger().SetLevel("debug")

    mvc.Configure(app.Party("/api",cors.AllowAll()).AllowMethods(iris.MethodOptions), basicMVC)

    app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./configs/iris.yml")))
}

func basicMVC(app *mvc.Application) {
    config := ReadAppConfigFile("./configs/app.yml")
    engine, error := repository.InitDataSouce(config)
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

    searchController := new(controller.SearchController)
    searchController.BlockSerivce = blockService
    searchController.TransactionService = transactionService
    app.Party("/search").Handle(searchController)

    convert := &bubuji.BubujiChainConvert{
       URL:     config.Sys.Url,
       ChainId: config.Sys.ChainId,
    }
    convert.Init()
    job := sys.SysJob{
       BlockService:       blockService,
       TransactionService: transactionService,
       Convert:            convert,
    }
    go job.Start()
}

func ReadAppConfigFile(path string) model.AppConfig {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }
    var config model.AppConfig
    if error := yaml.Unmarshal(data, &config); error != nil {
        panic(error)
    }
    return config
}
