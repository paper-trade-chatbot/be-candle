package candle

import (
	"context"

	"github.com/paper-trade-chatbot/be-proto/candle"
)

type CandleIntf interface {
	CreateCandles(ctx context.Context, in *candle.CreateCandlesReq) (*candle.CreateCandlesRes, error)
	GetCandles(ctx context.Context, in *candle.GetCandlesReq) (*candle.GetCandlesRes, error)
}

type CandleImpl struct {
	CandleClient candle.CandleServiceClient
}

func New() CandleIntf {
	return &CandleImpl{}
}

func (impl *CandleImpl) CreateCandles(ctx context.Context, in *candle.CreateCandlesReq) (*candle.CreateCandlesRes, error) {
	return nil, nil
}

func (impl *CandleImpl) GetCandles(ctx context.Context, in *candle.GetCandlesReq) (*candle.GetCandlesRes, error) {
	return nil, nil
}
