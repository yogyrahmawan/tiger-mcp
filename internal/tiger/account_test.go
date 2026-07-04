package tiger

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/tigerfintech/openapi-go-sdk/model"
)

func TestAssets_HappyPath(t *testing.T) {
	fake := &fakeTradeAPI{
		assets: []model.Asset{
			{Account: "U123", Currency: "USD", BuyingPower: 5000, NetLiquidation: 10000},
		},
	}
	c := &Client{tradeClient: fake}

	got, err := c.Assets(context.Background())
	if err != nil {
		t.Fatalf("Assets returned unexpected error: %v", err)
	}
	want := []Asset{{Account: "U123", Currency: "USD", BuyingPower: 5000, NetLiquidation: 10000}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assets() = %+v, want %+v", got, want)
	}
}

func TestAssets_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{tradeClient: &fakeTradeAPI{assetsErr: wantErr}}

	_, err := c.Assets(context.Background())
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}

func TestPositions_HappyPath(t *testing.T) {
	fake := &fakeTradeAPI{
		positions: []model.Position{
			{Symbol: "AAPL", PositionQty: 10, AverageCost: 150, MarketValue: 2000},
		},
	}
	c := &Client{tradeClient: fake}

	got, err := c.Positions(context.Background())
	if err != nil {
		t.Fatalf("Positions returned unexpected error: %v", err)
	}
	want := []Position{{Symbol: "AAPL", PositionQty: 10, AverageCost: 150, MarketValue: 2000}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Positions() = %+v, want %+v", got, want)
	}
}

func TestPositions_EmptyResultIsNotAnError(t *testing.T) {
	c := &Client{tradeClient: &fakeTradeAPI{positions: []model.Position{}}}

	got, err := c.Positions(context.Background())
	if err != nil {
		t.Fatalf("empty positions should not be an error, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Positions() = %+v, want empty", got)
	}
}

func TestPositions_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{tradeClient: &fakeTradeAPI{positionsErr: wantErr}}

	_, err := c.Positions(context.Background())
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}

func TestOrders_HappyPath(t *testing.T) {
	fake := &fakeTradeAPI{
		orders: []model.Order{
			{ID: 1, Symbol: "AAPL", Status: "FILLED", TotalQuantity: 10, FilledQuantity: 10},
		},
	}
	c := &Client{tradeClient: fake}

	got, err := c.Orders(context.Background())
	if err != nil {
		t.Fatalf("Orders returned unexpected error: %v", err)
	}
	want := []Order{{ID: 1, Symbol: "AAPL", Status: "FILLED", TotalQuantity: 10, FilledQuantity: 10}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Orders() = %+v, want %+v", got, want)
	}
}

func TestOrders_EmptyResultIsNotAnError(t *testing.T) {
	c := &Client{tradeClient: &fakeTradeAPI{orders: []model.Order{}}}

	got, err := c.Orders(context.Background())
	if err != nil {
		t.Fatalf("empty orders should not be an error, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Orders() = %+v, want empty", got)
	}
}

func TestOrders_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	c := &Client{tradeClient: &fakeTradeAPI{ordersErr: wantErr}}

	_, err := c.Orders(context.Background())
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
