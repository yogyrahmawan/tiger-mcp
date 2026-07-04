package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeQuoteFetcher struct {
	quotes     []tiger.Quote
	err        error
	calledWith []string
}

func (f *fakeQuoteFetcher) RealTimeQuotes(_ context.Context, symbols []string) ([]tiger.Quote, error) {
	f.calledWith = symbols
	if f.err != nil {
		return nil, f.err
	}
	return f.quotes, nil
}

func TestGetQuoteHandler_HappyPath(t *testing.T) {
	fake := &fakeQuoteFetcher{
		quotes: []tiger.Quote{
			{Symbol: "AAPL", LatestPrice: 200.5},
			{Symbol: "TSLA", LatestPrice: 300.1},
		},
	}
	handler := getQuoteHandler(fake)

	result, output, err := handler(context.Background(), nil, GetQuoteInput{Symbols: []string{"AAPL", "TSLA"}})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil CallToolResult on success, got %+v", result)
	}
	if !reflect.DeepEqual(fake.calledWith, []string{"AAPL", "TSLA"}) {
		t.Errorf("fetcher called with %v, want [AAPL TSLA]", fake.calledWith)
	}
	if !reflect.DeepEqual(output.Quotes, fake.quotes) {
		t.Errorf("output.Quotes = %+v, want %+v", output.Quotes, fake.quotes)
	}
}

func TestGetQuoteHandler_EmptySymbols(t *testing.T) {
	fake := &fakeQuoteFetcher{quotes: []tiger.Quote{{Symbol: "AAPL"}}}
	handler := getQuoteHandler(fake)

	_, _, err := handler(context.Background(), nil, GetQuoteInput{Symbols: nil})
	if err == nil {
		t.Fatal("expected validation error for empty symbols, got nil")
	}
	if fake.calledWith != nil {
		t.Errorf("fetcher should not be called on validation error, got calledWith=%v", fake.calledWith)
	}
}

func TestGetQuoteHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeQuoteFetcher{err: wantErr}
	handler := getQuoteHandler(fake)

	_, _, err := handler(context.Background(), nil, GetQuoteInput{Symbols: []string{"AAPL"}})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
