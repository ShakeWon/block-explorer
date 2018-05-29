package service

import (
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/model/vo/response"
    "github.com/shakewon/block-explorer/model/po"
)

type BlockService struct {
    Bs repository.BlockRepo
}

func (bs *BlockService) Count() (int64, error) {
    return bs.Bs.Count()
}

func (bs *BlockService) Page(index, pageSize int) ([]response.Block, error) {
    data, error := bs.Bs.Page(index, pageSize)
    if error != nil {
        return nil, error
    }
    var resp []response.Block
    for _, p := range data {
        resp = append(resp, response.Block{
            Height:         p.Height,
            Hash:           p.Hash,
            ChainId:        p.ChainId,
            Time:           p.Time,
            NumTxs:         p.NumTxs,
            LastCommitHash: p.LastCommitHash,
            DataHash:       p.DataHash,
            ValidatorsHash: p.ValidatorsHash,
            AppHash:        p.AppHash,
            Reward:         p.Reward,
            CoinBase:       p.CoinBase,
        })
    }
    return resp, error
}

/**
 *  按照高度查询
 */
func (bs *BlockService) Query(height int) (*response.Block, error) {
    if p, error := bs.Bs.Query(height); error == nil && p!=nil {
        return &response.Block{
            Height:         p.Height,
            Hash:           p.Hash,
            ChainId:        p.ChainId,
            Time:           p.Time,
            NumTxs:         p.NumTxs,
            LastCommitHash: p.LastCommitHash,
            DataHash:       p.DataHash,
            ValidatorsHash: p.ValidatorsHash,
            AppHash:        p.AppHash,
            Reward:         p.Reward,
            CoinBase:       p.CoinBase,
        },error
    } else {
        return nil,error
    }
}

func (bs *BlockService) Save(blocks []po.Block) error {
    return bs.Bs.Save(blocks)
}
