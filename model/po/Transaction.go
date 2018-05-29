package po

import "time"

type Transaction struct {
    ID         int64
    Hash       string `xormimpl:"unique"`
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
    Height     int64 `xormimpl:"index"`
    TxType     string
    Fee        int64
}
