package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeOrdersFetcher struct {
	orders []tiger.Order
	err    error
}

func (f *fakeOrdersFetcher) Orders(_ context.Context) ([]tiger.Order, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.orders, nil
}

func TestGetOrdersHandler_HappyPath(t *testing.T) {
	fake := &fakeOrdersFetcher{
		orders: []tiger.Order{{ID: 1, Symbol: "AAPL", Status: "FILLED"}},
	}
	handler := getOrdersHandler(fake)

	_, output, err := handler(context.Background(), nil, GetOrdersInput{})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(output.Orders, fake.orders) {
		t.Errorf("output.Orders = %+v, want %+v", output.Orders, fake.orders)
	}
}

func TestGetOrdersHandler_EmptyOrdersIsNotAnError(t *testing.T) {
	fake := &fakeOrdersFetcher{orders: []tiger.Order{}}
	handler := getOrdersHandler(fake)

	_, output, err := handler(context.Background(), nil, GetOrdersInput{})
	if err != nil {
		t.Fatalf("empty orders should not be an error, got: %v", err)
	}
	if len(output.Orders) != 0 {
		t.Errorf("output.Orders = %+v, want empty", output.Orders)
	}
}

func TestGetOrdersHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeOrdersFetcher{err: wantErr}
	handler := getOrdersHandler(fake)

	_, _, err := handler(context.Background(), nil, GetOrdersInput{})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
