package sys

import (
    "github.com/shakewon/block-explorer/service"
    "github.com/shakewon/block-explorer/third"
    "sync"
    "github.com/Sirupsen/logrus"
    "time"
)

type SysJob struct {
    service.BlockService
    service.TransactionService
    Convert third.DataConverter
}

var (
    gape int64 = 200
)

func (job *SysJob) Start() {

    for range time.NewTicker(1 * time.Second).C {
        blocks, error := job.Bs.Page(1,1,"","")
        if error != nil {
            panic(error)
        }
        var heightInDB int64 = 0
        if len(blocks)>0 {
            heightInDB = blocks[0].Height
        }
        heightOfNode := int64((job.Convert).Height())

        for i := heightInDB + 1; i+gape <= heightOfNode; i = i + gape {
            var group sync.WaitGroup
            for j := i; j < i+gape; j++ {
                group.Add(1)
                go job.work(&group, j)
            }
            group.Wait()
        }
    }
}

func (job *SysJob) work(group *sync.WaitGroup, h int64) {
    defer group.Done()
    if data := (job.Convert).Block(int(h)); data != nil {
        if error := job.Bs.Save(data.Blocks); error != nil {
            logrus.Error(error)
        }

        if error := job.Ts.Save(data.Txs); error != nil {
            logrus.Error(error)
        }
    }

}
