package tiger

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

func TestMarketStatus_HappyPath(t *testing.T) {
	fake := &fakeQuoteAPI{
		marketStates: []model.MarketState{
			{Market: "US", MarketStatus: "TRADING", Status: "OPEN", OpenTime: "09:30"},
		},
	}
	c := &Client{quoteClient: fake}

	got, err := c.MarketStatus(context.Background(), "US")
	if err != nil {
		t.Fatalf("MarketStatus returned unexpected error: %v", err)
	}
	want := []MarketState{{Market: "US", MarketStatus: "TRADING", Status: "OPEN", OpenTime: "09:30"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MarketStatus() = %+v, want %+v", got, want)
	}
}

func TestMarketStatus_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{quoteClient: &fakeQuoteAPI{marketErr: wantErr}}

	_, err := c.MarketStatus(context.Background(), "US")
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}

func TestKline_HappyPath(t *testing.T) {
	fake := &fakeQuoteAPI{
		klines: []model.Kline{
			{
				Symbol: "AAPL",
				Period: "day",
				Items:  []model.KlineItem{{Time: 1, Open: 100, Close: 101, High: 102, Low: 99, Volume: 500}},
			},
		},
	}
	c := &Client{quoteClient: fake}

	got, err := c.Kline(context.Background(), "AAPL", "day")
	if err != nil {
		t.Fatalf("Kline returned unexpected error: %v", err)
	}
	want := &Kline{
		Symbol: "AAPL",
		Period: "day",
		Items:  []KlineItem{{Time: 1, Open: 100, Close: 101, High: 102, Low: 99, Volume: 500}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Kline() = %+v, want %+v", got, want)
	}
}

func TestKline_EmptyResult(t *testing.T) {
	c := &Client{quoteClient: &fakeQuoteAPI{klines: nil}}

	got, err := c.Kline(context.Background(), "AAPL", "day")
	if err != nil {
		t.Fatalf("Kline returned unexpected error: %v", err)
	}
	want := &Kline{Symbol: "AAPL", Period: "day"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Kline() = %+v, want %+v", got, want)
	}
}

func TestKline_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{quoteClient: &fakeQuoteAPI{klineErr: wantErr}}

	_, err := c.Kline(context.Background(), "AAPL", "day")
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}

func TestDepth_HappyPath(t *testing.T) {
	fake := &fakeQuoteAPI{
		depths: []model.Depth{
			{
				Symbol: "AAPL",
				Asks:   []model.DepthLevel{{Price: 201, Count: 1, Volume: 100}},
				Bids:   []model.DepthLevel{{Price: 200, Count: 2, Volume: 200}},
			},
		},
	}
	c := &Client{quoteClient: fake}

	got, err := c.Depth(context.Background(), []string{"AAPL"}, "US")
	if err != nil {
		t.Fatalf("Depth returned unexpected error: %v", err)
	}
	want := []Depth{
		{
			Symbol: "AAPL",
			Asks:   []DepthLevel{{Price: 201, Count: 1, Volume: 100}},
			Bids:   []DepthLevel{{Price: 200, Count: 2, Volume: 200}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Depth() = %+v, want %+v", got, want)
	}
}

func TestDepth_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{quoteClient: &fakeQuoteAPI{depthErr: wantErr}}

	_, err := c.Depth(context.Background(), []string{"AAPL"}, "US")
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
