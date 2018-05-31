package service

import (
    "github.com/shakewon/block-explorer/repository"
    "github.com/shakewon/block-explorer/model/vo/response"
    "github.com/shakewon/block-explorer/model/po"
)

type BlockService struct {
    Bs repository.BlockRepo
}

func (bs *BlockService) Count(height,hash string) (int64, error) {
    return bs.Bs.Count(height,hash)
}

func (bs *BlockService) Page(index, pageSize int,height,hash string) ([]response.Block, error) {
    data, error := bs.Bs.Page(index, pageSize,height,hash)
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

func (bs *BlockService) Save(blocks []po.Block) error {
    return bs.Bs.Save(blocks)
}


