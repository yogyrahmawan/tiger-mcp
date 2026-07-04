package tiger

import (
	"context"
	"fmt"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

// Asset is a curated view of one account's asset summary.
type Asset struct {
	Account        string  `json:"account"`
	Currency       string  `json:"currency"`
	BuyingPower    float64 `json:"buyingPower"`
	CashValue      float64 `json:"cashValue"`
	NetLiquidation float64 `json:"netLiquidation"`
	RealizedPnL    float64 `json:"realizedPnL"`
	UnrealizedPnL  float64 `json:"unrealizedPnL"`
}

// AssetsFetcher fetches account asset summaries.
type AssetsFetcher interface {
	Assets(ctx context.Context) ([]Asset, error)
}

// Position is a curated view of one held position.
type Position struct {
	Account       string  `json:"account"`
	Symbol        string  `json:"symbol"`
	SecType       string  `json:"secType"`
	Market        string  `json:"market"`
	Currency      string  `json:"currency"`
	Name          string  `json:"name"`
	PositionQty   float64 `json:"positionQty"`
	AverageCost   float64 `json:"averageCost"`
	LatestPrice   float64 `json:"latestPrice"`
	MarketValue   float64 `json:"marketValue"`
	RealizedPnl   float64 `json:"realizedPnl"`
	UnrealizedPnl float64 `json:"unrealizedPnl"`
}

// PositionsFetcher fetches currently held positions.
type PositionsFetcher interface {
	Positions(ctx context.Context) ([]Position, error)
}

// Order is a curated view of one order (historical or active).
type Order struct {
	ID             int64   `json:"id"`
	Symbol         string  `json:"symbol"`
	SecType        string  `json:"secType"`
	Market         string  `json:"market"`
	Currency       string  `json:"currency"`
	Action         string  `json:"action"`
	OrderType      string  `json:"orderType"`
	Status         string  `json:"status"`
	TotalQuantity  int64   `json:"totalQuantity"`
	FilledQuantity int64   `json:"filledQuantity"`
	LimitPrice     float64 `json:"limitPrice"`
	AvgFillPrice   float64 `json:"avgFillPrice"`
	Commission     float64 `json:"commission"`
	RealizedPnl    float64 `json:"realizedPnl"`
	OpenTime       int64   `json:"openTime"`
	UpdateTime     int64   `json:"updateTime"`
}

// OrdersFetcher fetches order history.
type OrdersFetcher interface {
	Orders(ctx context.Context) ([]Order, error)
}

// Assets returns account asset summaries via Tiger's assets endpoint.
func (c *Client) Assets(_ context.Context) ([]Asset, error) {
	assets, err := c.tradeClient.Assets(model.AssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("tiger: get assets: %w", err)
	}

	result := make([]Asset, 0, len(assets))
	for _, a := range assets {
		result = append(result, Asset{
			Account:        a.Account,
			Currency:       a.Currency,
			BuyingPower:    a.BuyingPower,
			CashValue:      a.CashValue,
			NetLiquidation: a.NetLiquidation,
			RealizedPnL:    a.RealizedPnL,
			UnrealizedPnL:  a.UnrealizedPnL,
		})
	}
	return result, nil
}

// Positions returns currently held positions via Tiger's positions endpoint.
func (c *Client) Positions(_ context.Context) ([]Position, error) {
	positions, err := c.tradeClient.Positions(model.PositionsRequest{})
	if err != nil {
		return nil, fmt.Errorf("tiger: get positions: %w", err)
	}

	result := make([]Position, 0, len(positions))
	for _, p := range positions {
		result = append(result, Position{
			Account:       p.Account,
			Symbol:        p.Symbol,
			SecType:       p.SecType,
			Market:        p.Market,
			Currency:      p.Currency,
			Name:          p.Name,
			PositionQty:   p.PositionQty,
			AverageCost:   p.AverageCost,
			LatestPrice:   p.LatestPrice,
			MarketValue:   p.MarketValue,
			RealizedPnl:   p.RealizedPnl,
			UnrealizedPnl: p.UnrealizedPnl,
		})
	}
	return result, nil
}

// Orders returns order history via Tiger's orders endpoint.
func (c *Client) Orders(_ context.Context) ([]Order, error) {
	orders, err := c.tradeClient.Orders(model.OrdersRequest{})
	if err != nil {
		return nil, fmt.Errorf("tiger: get orders: %w", err)
	}

	result := make([]Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, Order{
			ID:             o.ID,
			Symbol:         o.Symbol,
			SecType:        o.SecType,
			Market:         o.Market,
			Currency:       o.Currency,
			Action:         o.Action,
			OrderType:      o.OrderType,
			Status:         o.Status,
			TotalQuantity:  o.TotalQuantity,
			FilledQuantity: o.FilledQuantity,
			LimitPrice:     o.LimitPrice,
			AvgFillPrice:   o.AvgFillPrice,
			Commission:     o.Commission,
			RealizedPnl:    o.RealizedPnl,
			OpenTime:       o.OpenTime,
			UpdateTime:     o.UpdateTime,
		})
	}
	return result, nil
}
