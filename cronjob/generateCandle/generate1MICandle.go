package generateCandle

import (
	"context"
	"strconv"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/paper-trade-chatbot/be-candle/dao/candleDao"
	"github.com/paper-trade-chatbot/be-candle/database"
	"github.com/paper-trade-chatbot/be-candle/logging"
	"github.com/paper-trade-chatbot/be-candle/models/dbModels"
	"github.com/paper-trade-chatbot/be-candle/service"
	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/product"
	"github.com/paper-trade-chatbot/be-proto/quote"
	"github.com/shopspring/decimal"
)

func Generate1MICandle(ctx context.Context) error {

	now := time.Now().Truncate(time.Minute)

	db := database.GetDB()

	enabled := product.Status_Status_Enabled
	productRes, err := pagination.IteratePageGRPC[*product.GetProductsReq, *product.GetProductsRes](
		&product.GetProductsReq{
			Status:     &enabled,
			Pagination: pagination.NewPagination(3000),
		},
		func(req *product.GetProductsReq) (*product.GetProductsRes, error) {
			productData, err := service.Impl.ProductIntf.GetProducts(ctx, req)
			if err != nil {
				logging.Error(ctx, "[Generate1MICandle] GetProducts err: %v", err)
				return nil, err
			}
			return productData, nil
		},
	)
	if err != nil {
		logging.Error(ctx, "[Generate1MICandle] IteratePage err: %v", err)
		return err
	}

	productIDSet := mapset.NewSet[int64]()

	for _, r := range productRes {
		for _, p := range r.Product {
			productIDSet.Add(p.Id)
		}
	}

	from := now.Add(-time.Minute).Format("150405")
	to := now.Format("150405")
	quoteData, err := service.Impl.QuoteIntf.GetQuotes(ctx, &quote.GetQuotesReq{
		ProductIDs: productIDSet.ToSlice(),
		Flag:       quote.GetQuotesReq_GetFlag_Quote | quote.GetQuotesReq_GetFlag_Latest,
		GetFrom:    &from,
		GetTo:      &to,
	})
	if err != nil {
		logging.Error(ctx, "[Generate1MICandle] GetQuotes err: %v", err)
		return err
	}

	models := []*dbModels.CandleModel{}

	for _, q := range quoteData.Quotes {

		from := time.Date(0, 0, 0, now.Add(-time.Minute).Hour(), now.Add(-time.Minute).Minute(), now.Add(-time.Minute).Second(), 0, time.UTC)
		to := time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

		if _, ok := q.Quotes["latest"]; !ok {
			logging.Error(ctx, "[Generate1MICandle] quote [%s] not having latest quote.", q.ProductID)
			return err
		}
		latestPrice, _ := decimal.NewFromString(q.Quotes["latest"])
		delete(q.Quotes, "latest")

		closestToFrom := to // 用來找出最接近from的時間點
		tempOpenQuote := latestPrice

		closestToTo := from // 用來找出最接近to的時間點
		tempCloseQuote := latestPrice

		open := latestPrice
		close := latestPrice
		high := latestPrice
		low := latestPrice

		if isCrossingDate := (now.Hour() == 0 && now.Minute() == 0); isCrossingDate {
			to.Add(time.Hour * 24)
		}

		for k, v := range q.Quotes {
			quoteTime, err := time.Parse("150405", k)
			if err != nil {
				logging.Warn(ctx, "[Generate1MICandle] parse time error: %v", err)
				continue
			}

			price, _ := decimal.NewFromString(v)

			if k == "000000" {
				quoteTime = quoteTime.Add(time.Hour * 24)
			}

			if quoteTime.After(closestToTo) && !quoteTime.After(to) {
				closestToTo = quoteTime
				tempCloseQuote = price
			}

			if quoteTime.Before(closestToFrom) && quoteTime.After(from) {
				closestToFrom = quoteTime
				tempOpenQuote = price
			}

			if price.LessThan(low) {
				low = price
			}

			if high.LessThan(price) {
				high = price
			}
		}
		open = tempOpenQuote
		close = tempCloseQuote

		candleChart := &dbModels.CandleModel{
			ProductID:    uint64(q.ProductID),
			IntervalType: dbModels.IntervalType_1MI,
			Start:        now.Add(-time.Minute),
			Open:         open,
			Close:        close,
			High:         high,
			Low:          low,
		}

		models = append(models, candleChart)
	}

	if _, err := candleDao.News(db, models); err != nil {
		logging.Error(ctx, "[Generate1MICandle] news error: %v", err)
		return err
	}

	return nil
}

func Generate1MICandleKey() string {
	now := time.Now()
	key := "Generate1MICandle:" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute())
	return key
}
