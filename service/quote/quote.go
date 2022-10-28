package quote

import (
	"context"

	"github.com/paper-trade-chatbot/be-candle/config"
	"github.com/paper-trade-chatbot/be-proto/quote"
)

type QuoteIntf interface {
	AddProductQuoteSources(ctx context.Context, in *quote.AddProductQuoteSourcesReq) (*quote.AddProductQuoteSourcesRes, error)
	ModifyProductQuoteSource(ctx context.Context, in *quote.ModifyProductQuoteSourceReq) (*quote.ModifyProductQuoteSourceRes, error)
	GetQuotes(ctx context.Context, in *quote.GetQuotesReq) (*quote.GetQuotesRes, error)
	DeleteQuotes(ctx context.Context, in *quote.DeleteQuotesReq) (*quote.DeleteQuotesRes, error)
}

type QuoteImpl struct {
	QuoteClient quote.QuoteServiceClient
}

var (
	QuoteServiceHost    = config.GetString("QUOTE_GRPC_HOST")
	QuoteServerGRpcPort = config.GetString("QUOTE_GRPC_PORT")
)

func New(quoteClient quote.QuoteServiceClient) QuoteIntf {
	return &QuoteImpl{
		QuoteClient: quoteClient,
	}
}

func (impl *QuoteImpl) AddProductQuoteSources(ctx context.Context, in *quote.AddProductQuoteSourcesReq) (*quote.AddProductQuoteSourcesRes, error) {
	return impl.AddProductQuoteSources(ctx, in)
}

func (impl *QuoteImpl) ModifyProductQuoteSource(ctx context.Context, in *quote.ModifyProductQuoteSourceReq) (*quote.ModifyProductQuoteSourceRes, error) {
	return impl.ModifyProductQuoteSource(ctx, in)
}

func (impl *QuoteImpl) GetQuotes(ctx context.Context, in *quote.GetQuotesReq) (*quote.GetQuotesRes, error) {
	return impl.GetQuotes(ctx, in)
}

func (impl *QuoteImpl) DeleteQuotes(ctx context.Context, in *quote.DeleteQuotesReq) (*quote.DeleteQuotesRes, error) {
	return impl.DeleteQuotes(ctx, in)
}
