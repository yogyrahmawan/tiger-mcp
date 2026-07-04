package tiger

import (
	"context"
	"fmt"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

// MarketState reports a market's current trading status.
type MarketState struct {
	Market       string `json:"market"`
	MarketStatus string `json:"marketStatus"`
	Status       string `json:"status"`
	OpenTime     string `json:"openTime"`
}

// MarketStatusFetcher fetches the current status of a market.
type MarketStatusFetcher interface {
	MarketStatus(ctx context.Context, market string) ([]MarketState, error)
}

// KlineItem is a single historical bar.
type KlineItem struct {
	Time   int64   `json:"time"`
	Volume int64   `json:"volume"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Amount float64 `json:"amount"`
}

// Kline is a series of historical bars for one symbol at one period.
type Kline struct {
	Symbol string      `json:"symbol"`
	Period string      `json:"period"`
	Items  []KlineItem `json:"items"`
}

// KlineFetcher fetches historical K-line bars for a single symbol.
type KlineFetcher interface {
	Kline(ctx context.Context, symbol, period string) (*Kline, error)
}

// DepthLevel is a single price level on one side of an order book.
type DepthLevel struct {
	Price  float64 `json:"price"`
	Count  int     `json:"count"`
	Volume int64   `json:"volume"`
}

// Depth is the order book snapshot for one symbol.
type Depth struct {
	Symbol string       `json:"symbol"`
	Asks   []DepthLevel `json:"asks"`
	Bids   []DepthLevel `json:"bids"`
}

// DepthFetcher fetches order book depth for a batch of symbols.
type DepthFetcher interface {
	Depth(ctx context.Context, symbols []string, market string) ([]Depth, error)
}

// MarketStatus returns the current status of market via Tiger's
// market_status endpoint.
func (c *Client) MarketStatus(_ context.Context, market string) ([]MarketState, error) {
	states, err := c.quoteClient.GetMarketState(market)
	if err != nil {
		return nil, fmt.Errorf("tiger: get market state: %w", err)
	}

	result := make([]MarketState, 0, len(states))
	for _, s := range states {
		result = append(result, MarketState{
			Market:       s.Market,
			MarketStatus: s.MarketStatus,
			Status:       s.Status,
			OpenTime:     s.OpenTime,
		})
	}
	return result, nil
}

// Kline returns historical bars for symbol at period via Tiger's kline
// endpoint.
func (c *Client) Kline(_ context.Context, symbol, period string) (*Kline, error) {
	klines, err := c.quoteClient.GetKline(symbol, period)
	if err != nil {
		return nil, fmt.Errorf("tiger: get kline: %w", err)
	}
	if len(klines) == 0 {
		return &Kline{Symbol: symbol, Period: period}, nil
	}

	k := klines[0]
	items := make([]KlineItem, 0, len(k.Items))
	for _, it := range k.Items {
		items = append(items, KlineItem{
			Time:   it.Time,
			Volume: it.Volume,
			Open:   it.Open,
			Close:  it.Close,
			High:   it.High,
			Low:    it.Low,
			Amount: it.Amount,
		})
	}
	return &Kline{Symbol: k.Symbol, Period: k.Period, Items: items}, nil
}

// Depth returns order book depth for symbols via Tiger's quote_depth
// endpoint. market may be empty to use Tiger's default.
func (c *Client) Depth(_ context.Context, symbols []string, market string) ([]Depth, error) {
	depths, err := c.quoteClient.GetQuoteDepth(model.DepthQuoteRequest{
		Symbols: symbols,
		Market:  market,
	})
	if err != nil {
		return nil, fmt.Errorf("tiger: get quote depth: %w", err)
	}

	result := make([]Depth, 0, len(depths))
	for _, d := range depths {
		result = append(result, Depth{
			Symbol: d.Symbol,
			Asks:   mapDepthLevels(d.Asks),
			Bids:   mapDepthLevels(d.Bids),
		})
	}
	return result, nil
}

func mapDepthLevels(levels []model.DepthLevel) []DepthLevel {
	result := make([]DepthLevel, 0, len(levels))
	for _, l := range levels {
		result = append(result, DepthLevel{Price: l.Price, Count: l.Count, Volume: l.Volume})
	}
	return result
}
