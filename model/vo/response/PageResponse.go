package response

import "time"

type BaseResponse struct {
    Success bool
    Error   error
    Data    interface{}
}

type PageTrxResponse struct {
    BaseResponse
    Total int64
    Data  []Transaction
}

type PageBlockResponse struct {
    Total int64
    Data  []Block
}

type Block struct {
    Height         int64
    Hash           string
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

type Transaction struct {
    Hash       string
    Payload    string
    PayloadHex string
    FromAddr   string
    ToAddr     string
    Receipt    string
    Amount     string
    Nonce      string
    Gas        string
    Size       int64
    Block      string
    Contract   string
    Time       time.Time
    Height     int64
    TxType     string
    Fee        int64
}
