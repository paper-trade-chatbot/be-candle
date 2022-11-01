package candle

import (
	"context"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/paper-trade-chatbot/be-candle/dao/candleDao"
	"github.com/paper-trade-chatbot/be-candle/database"
	"github.com/paper-trade-chatbot/be-candle/logging"
	"github.com/paper-trade-chatbot/be-candle/models/dbModels"
	"github.com/paper-trade-chatbot/be-candle/service"
	common "github.com/paper-trade-chatbot/be-common"
	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/candle"
	"github.com/paper-trade-chatbot/be-proto/product"
	"github.com/shopspring/decimal"
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

	productIDSet := mapset.NewSet[int64]()

	for _, c := range in.GetCandleCharts() {
		productIDSet.Add(c.GetProductID())
	}

	res, err := pagination.IteratePageGRPC[*product.GetProductsReq, *product.GetProductsRes](
		&product.GetProductsReq{
			Id:         productIDSet.ToSlice(),
			Pagination: pagination.NewPagination(3000),
		},
		func(req *product.GetProductsReq) (*product.GetProductsRes, error) {
			productData, err := service.Impl.ProductIntf.GetProducts(ctx, req)
			if err != nil {
				logging.Error(ctx, "[CreateCandles] GetProducts err: %v", err)
				return nil, err
			}
			return productData, nil
		},
	)
	if err != nil {
		logging.Error(ctx, "[CreateCandles] IteratePage err: %v", err)
		return nil, err
	}

	for _, r := range res {
		for _, p := range r.Product {
			productIDSet.Remove(p.Id)
		}
	}
	if productIDSet.Cardinality() > 0 {
		logging.Error(ctx, "[CreateCandles] no such product: %#q", productIDSet.ToSlice())
		return nil, common.ErrNoSuchProduct
	}

	models := []*dbModels.CandleModel{}
	for _, cc := range in.GetCandleCharts() {
		for _, s := range cc.GetCandleSticks() {

			open, _ := decimal.NewFromString(s.GetOpen())
			close, _ := decimal.NewFromString(s.GetClose())
			high, _ := decimal.NewFromString(s.GetHigh())
			low, _ := decimal.NewFromString(s.GetLow())
			volume, _ := decimal.NewFromString(s.GetVolume())

			model := &dbModels.CandleModel{
				ProductID:    uint64(cc.GetProductID()),
				IntervalType: dbModels.IntervalType(cc.GetIntervalType()),
				Start:        time.Unix(s.GetStart(), 0),
				Open:         open,
				Close:        close,
				High:         high,
				Low:          low,
				Volume:       volume,
			}
			models = append(models, model)
		}

	}

	db := database.GetDB()
	count, err := candleDao.News(db, models)
	if err != nil {
		logging.Error(ctx, "[CreateCandles] candleDao.News error: %v", err)
		return nil, err
	}

	return &candle.CreateCandlesRes{
		TotalSuccess: int32(count),
	}, nil
}

func (impl *CandleImpl) GetCandles(ctx context.Context, in *candle.GetCandlesReq) (*candle.GetCandlesRes, error) {
	logging.Info(ctx, "[GetCandles] test")
	db := database.GetDB()

	startTime := time.Unix(in.StartTime, 0)
	endTime := time.Unix(in.EndTime, 0)
	productIDIn := []uint64{}
	for _, p := range in.ProductID {
		productIDIn = append(productIDIn, uint64(p))
	}
	orders := []*candleDao.Order{}
	for i := range in.OrderBy {
		order := &candleDao.Order{
			Column:    candleDao.OrderColumn(in.OrderBy[i]),
			Direction: candleDao.OrderDirection(in.OrderDirection[i]),
		}
		orders = append(orders, order)
	}

	queryModel := &candleDao.QueryModel{
		IntervalType: dbModels.IntervalType(in.IntervalType),
		StartFrom:    &startTime,
		StartTo:      &endTime,
		ProductIDIn:  productIDIn,
		OrderBy:      orders,
	}

	models, paginationInfo, err := candleDao.GetsWithPagination(db, queryModel, in.Pagination)
	if err != nil {
		return nil, err
	}

	candles := []*candle.GetCandlesResElement{}

	if len(models) == 0 {
		return &candle.GetCandlesRes{
			Candles:        candles,
			PaginationInfo: paginationInfo,
		}, nil
	}

	for _, m := range models {

		c := &candle.GetCandlesResElement{
			ProductID:    int64(m.ProductID),
			IntervalType: candle.IntervalType(m.IntervalType),
			CandleSticks: &candle.CandleStick{
				Start:  m.Start.Unix(),
				Open:   m.Open.String(),
				Close:  m.Close.String(),
				High:   m.High.String(),
				Low:    m.Low.String(),
				Volume: m.Volume.String(),
			},
		}
		candles = append(candles, c)
	}

	return &candle.GetCandlesRes{
		Candles:        candles,
		PaginationInfo: paginationInfo,
	}, nil
}
