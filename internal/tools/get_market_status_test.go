package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeMarketStatusFetcher struct {
	states     []tiger.MarketState
	err        error
	calledWith string
}

func (f *fakeMarketStatusFetcher) MarketStatus(_ context.Context, market string) ([]tiger.MarketState, error) {
	f.calledWith = market
	if f.err != nil {
		return nil, f.err
	}
	return f.states, nil
}

func TestGetMarketStatusHandler_HappyPath(t *testing.T) {
	fake := &fakeMarketStatusFetcher{
		states: []tiger.MarketState{{Market: "US", MarketStatus: "TRADING", Status: "OPEN"}},
	}
	handler := getMarketStatusHandler(fake)

	_, output, err := handler(context.Background(), nil, GetMarketStatusInput{Market: "US"})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if fake.calledWith != "US" {
		t.Errorf("fetcher called with %q, want %q", fake.calledWith, "US")
	}
	if !reflect.DeepEqual(output.MarketStates, fake.states) {
		t.Errorf("output.MarketStates = %+v, want %+v", output.MarketStates, fake.states)
	}
}

func TestGetMarketStatusHandler_InvalidMarket(t *testing.T) {
	fake := &fakeMarketStatusFetcher{states: []tiger.MarketState{{Market: "US"}}}
	handler := getMarketStatusHandler(fake)

	_, _, err := handler(context.Background(), nil, GetMarketStatusInput{Market: "MARS"})
	if err == nil {
		t.Fatal("expected validation error for invalid market, got nil")
	}
	if fake.calledWith != "" {
		t.Errorf("fetcher should not be called on validation error, got calledWith=%q", fake.calledWith)
	}
}

func TestGetMarketStatusHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeMarketStatusFetcher{err: wantErr}
	handler := getMarketStatusHandler(fake)

	_, _, err := handler(context.Background(), nil, GetMarketStatusInput{Market: "US"})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
