// Package tiger wraps the Tiger Open API Go SDK's quote client for use by
// this server's MCP tools.
package tiger

import (
	"context"
	"fmt"

	tigerclient "github.com/tigerfintech/openapi-go-sdk/client"
	tigerconfig "github.com/tigerfintech/openapi-go-sdk/config"
	"github.com/tigerfintech/openapi-go-sdk/model"
	"github.com/tigerfintech/openapi-go-sdk/quote"
	"github.com/tigerfintech/openapi-go-sdk/trade"

	"github.com/yogyrahmawan/tiger-mcp/internal/config"
)

// Quote is a real-time price snapshot for a single symbol.
type Quote struct {
	Symbol      string  `json:"symbol"`
	Open        float64 `json:"open"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Close       float64 `json:"close"`
	PreClose    float64 `json:"preClose"`
	LatestPrice float64 `json:"latestPrice"`
	LatestTime  int64   `json:"latestTime"`
	Volume      int64   `json:"volume"`
	Change      float64 `json:"change"`
	ChangeRate  float64 `json:"changeRate"`
	Status      string  `json:"status"`
}

// QuoteFetcher fetches real-time quotes for a batch of symbols. It is
// implemented by Client (backed by the real Tiger API) and can be faked in
// tests.
type QuoteFetcher interface {
	RealTimeQuotes(ctx context.Context, symbols []string) ([]Quote, error)
}

// Client is backed by the real Tiger Open API and implements this package's
// fetcher interfaces (quote and trade).
type Client struct {
	quoteClient *quote.QuoteClient
	tradeClient *trade.TradeClient
}

// NewClient builds Tiger quote and trade clients from the given credentials.
func NewClient(cfg *config.Config) (*Client, error) {
	tigerCfg, err := tigerconfig.NewClientConfig(
		tigerconfig.WithTigerID(cfg.TigerID),
		tigerconfig.WithPrivateKey(cfg.PrivateKey),
		tigerconfig.WithAccount(cfg.Account),
	)
	if err != nil {
		return nil, fmt.Errorf("tiger: build client config: %w", err)
	}

	quoteHTTPClient := tigerclient.NewQuoteHttpClient(tigerCfg)
	tradeHTTPClient := tigerclient.NewHttpClient(tigerCfg)

	return &Client{
		quoteClient: quote.NewQuoteClient(quoteHTTPClient),
		tradeClient: trade.NewTradeClientFromConfig(tradeHTTPClient, tigerCfg),
	}, nil
}

// RealTimeQuotes returns real-time quotes for the given symbols via Tiger's
// quote_real_time endpoint.
func (c *Client) RealTimeQuotes(_ context.Context, symbols []string) ([]Quote, error) {
	briefs, err := c.quoteClient.GetBrief(model.BriefRequest{Symbols: symbols})
	if err != nil {
		return nil, fmt.Errorf("tiger: get brief: %w", err)
	}

	quotes := make([]Quote, 0, len(briefs))
	for _, b := range briefs {
		quotes = append(quotes, Quote{
			Symbol:      b.Symbol,
			Open:        b.Open,
			High:        b.High,
			Low:         b.Low,
			Close:       b.Close,
			PreClose:    b.PreClose,
			LatestPrice: b.LatestPrice,
			LatestTime:  b.LatestTime,
			Volume:      b.Volume,
			Change:      b.Change,
			ChangeRate:  b.ChangeRate,
			Status:      b.Status,
		})
	}
	return quotes, nil
}
