package tiger

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

// fakeQuoteAPI fakes quoteAPI for tests, avoiding any real Tiger network call.
type fakeQuoteAPI struct {
	briefs   []model.Brief
	briefErr error

	marketStates []model.MarketState
	marketErr    error

	klines   []model.Kline
	klineErr error

	depths   []model.Depth
	depthErr error
}

func (f *fakeQuoteAPI) GetBrief(_ model.BriefRequest) ([]model.Brief, error) {
	if f.briefErr != nil {
		return nil, f.briefErr
	}
	return f.briefs, nil
}

func (f *fakeQuoteAPI) GetMarketState(_ string) ([]model.MarketState, error) {
	if f.marketErr != nil {
		return nil, f.marketErr
	}
	return f.marketStates, nil
}

func (f *fakeQuoteAPI) GetKline(_, _ string) ([]model.Kline, error) {
	if f.klineErr != nil {
		return nil, f.klineErr
	}
	return f.klines, nil
}

func (f *fakeQuoteAPI) GetQuoteDepth(_ model.DepthQuoteRequest) ([]model.Depth, error) {
	if f.depthErr != nil {
		return nil, f.depthErr
	}
	return f.depths, nil
}

// fakeTradeAPI fakes tradeAPI for tests, avoiding any real Tiger network call.
type fakeTradeAPI struct {
	assets    []model.Asset
	assetsErr error

	positions    []model.Position
	positionsErr error

	orders    []model.Order
	ordersErr error
}

func (f *fakeTradeAPI) Assets(_ model.AssetsRequest) ([]model.Asset, error) {
	if f.assetsErr != nil {
		return nil, f.assetsErr
	}
	return f.assets, nil
}

func (f *fakeTradeAPI) Positions(_ model.PositionsRequest) ([]model.Position, error) {
	if f.positionsErr != nil {
		return nil, f.positionsErr
	}
	return f.positions, nil
}

func (f *fakeTradeAPI) Orders(_ model.OrdersRequest) ([]model.Order, error) {
	if f.ordersErr != nil {
		return nil, f.ordersErr
	}
	return f.orders, nil
}

func TestRealTimeQuotes_HappyPath(t *testing.T) {
	fake := &fakeQuoteAPI{
		briefs: []model.Brief{
			{Symbol: "AAPL", LatestPrice: 200.5, Open: 199, High: 201, Low: 198, Close: 200.5, Volume: 1000},
		},
	}
	c := &Client{quoteClient: fake}

	got, err := c.RealTimeQuotes(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatalf("RealTimeQuotes returned unexpected error: %v", err)
	}
	want := []Quote{{Symbol: "AAPL", LatestPrice: 200.5, Open: 199, High: 201, Low: 198, Close: 200.5, Volume: 1000}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RealTimeQuotes() = %+v, want %+v", got, want)
	}
}

func TestRealTimeQuotes_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{quoteClient: &fakeQuoteAPI{briefErr: wantErr}}

	_, err := c.RealTimeQuotes(context.Background(), []string{"AAPL"})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
