package po

import "time"

type Block struct {
    ID             int64
    Height         int64  `xormimpl:"index"`
    Hash           string `xormimpl:"index"`
    ChainId        string
    Time           time.Time
    NumTxs         int64
    LastCommitHash string
    DataHash       string
    ValidatorsHash string
    AppHash        string
    Reward         int64
    CoinBase       string
}
