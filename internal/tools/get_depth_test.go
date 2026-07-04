package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeDepthFetcher struct {
	depths           []tiger.Depth
	err              error
	calledWithSymbol []string
	calledWithMarket string
}

func (f *fakeDepthFetcher) Depth(_ context.Context, symbols []string, market string) ([]tiger.Depth, error) {
	f.calledWithSymbol = symbols
	f.calledWithMarket = market
	if f.err != nil {
		return nil, f.err
	}
	return f.depths, nil
}

func TestGetDepthHandler_HappyPath(t *testing.T) {
	fake := &fakeDepthFetcher{
		depths: []tiger.Depth{
			{Symbol: "AAPL", Asks: []tiger.DepthLevel{{Price: 201.0, Volume: 100}}},
		},
	}
	handler := getDepthHandler(fake)

	_, output, err := handler(context.Background(), nil, GetDepthInput{Symbols: []string{"AAPL", "TSLA"}, Market: "US"})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(fake.calledWithSymbol, []string{"AAPL", "TSLA"}) {
		t.Errorf("fetcher called with symbols=%v, want [AAPL TSLA]", fake.calledWithSymbol)
	}
	if fake.calledWithMarket != "US" {
		t.Errorf("fetcher called with market=%q, want US", fake.calledWithMarket)
	}
	if !reflect.DeepEqual(output.Depths, fake.depths) {
		t.Errorf("output.Depths = %+v, want %+v", output.Depths, fake.depths)
	}
}

func TestGetDepthHandler_EmptySymbols(t *testing.T) {
	fake := &fakeDepthFetcher{depths: []tiger.Depth{{Symbol: "AAPL"}}}
	handler := getDepthHandler(fake)

	_, _, err := handler(context.Background(), nil, GetDepthInput{Symbols: nil, Market: "US"})
	if err == nil {
		t.Fatal("expected validation error for empty symbols, got nil")
	}
	if fake.calledWithSymbol != nil {
		t.Errorf("fetcher should not be called on validation error, got calledWithSymbol=%v", fake.calledWithSymbol)
	}
}

func TestGetDepthHandler_InvalidMarket(t *testing.T) {
	fake := &fakeDepthFetcher{depths: []tiger.Depth{{Symbol: "AAPL"}}}
	handler := getDepthHandler(fake)

	_, _, err := handler(context.Background(), nil, GetDepthInput{Symbols: []string{"AAPL"}, Market: "MARS"})
	if err == nil {
		t.Fatal("expected validation error for invalid market, got nil")
	}
	if fake.calledWithSymbol != nil {
		t.Errorf("fetcher should not be called on validation error, got calledWithSymbol=%v", fake.calledWithSymbol)
	}
}

func TestGetDepthHandler_EmptyMarketAllowed(t *testing.T) {
	fake := &fakeDepthFetcher{depths: []tiger.Depth{{Symbol: "AAPL"}}}
	handler := getDepthHandler(fake)

	_, _, err := handler(context.Background(), nil, GetDepthInput{Symbols: []string{"AAPL"}, Market: ""})
	if err != nil {
		t.Fatalf("empty market should be allowed (Tiger default), got error: %v", err)
	}
}

func TestGetDepthHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeDepthFetcher{err: wantErr}
	handler := getDepthHandler(fake)

	_, _, err := handler(context.Background(), nil, GetDepthInput{Symbols: []string{"AAPL"}})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
