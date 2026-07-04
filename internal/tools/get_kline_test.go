package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeKlineFetcher struct {
	kline            *tiger.Kline
	err              error
	calledWithSymbol string
	calledWithPeriod string
}

func (f *fakeKlineFetcher) Kline(_ context.Context, symbol, period string) (*tiger.Kline, error) {
	f.calledWithSymbol = symbol
	f.calledWithPeriod = period
	if f.err != nil {
		return nil, f.err
	}
	return f.kline, nil
}

func TestGetKlineHandler_HappyPath(t *testing.T) {
	fake := &fakeKlineFetcher{
		kline: &tiger.Kline{
			Symbol: "AAPL",
			Period: "day",
			Items:  []tiger.KlineItem{{Time: 1, Open: 100, Close: 101}},
		},
	}
	handler := getKlineHandler(fake)

	_, output, err := handler(context.Background(), nil, GetKlineInput{Symbol: "AAPL", Period: "day"})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if fake.calledWithSymbol != "AAPL" || fake.calledWithPeriod != "day" {
		t.Errorf("fetcher called with symbol=%q period=%q, want AAPL/day", fake.calledWithSymbol, fake.calledWithPeriod)
	}
	if !reflect.DeepEqual(output.Kline, *fake.kline) {
		t.Errorf("output.Kline = %+v, want %+v", output.Kline, *fake.kline)
	}
}

func TestGetKlineHandler_EmptySymbol(t *testing.T) {
	fake := &fakeKlineFetcher{kline: &tiger.Kline{Symbol: "AAPL"}}
	handler := getKlineHandler(fake)

	_, _, err := handler(context.Background(), nil, GetKlineInput{Symbol: "", Period: "day"})
	if err == nil {
		t.Fatal("expected validation error for empty symbol, got nil")
	}
	if fake.calledWithSymbol != "" {
		t.Errorf("fetcher should not be called on validation error, got calledWithSymbol=%q", fake.calledWithSymbol)
	}
}

func TestGetKlineHandler_InvalidPeriod(t *testing.T) {
	fake := &fakeKlineFetcher{kline: &tiger.Kline{Symbol: "AAPL"}}
	handler := getKlineHandler(fake)

	_, _, err := handler(context.Background(), nil, GetKlineInput{Symbol: "AAPL", Period: "fortnight"})
	if err == nil {
		t.Fatal("expected validation error for invalid period, got nil")
	}
	if fake.calledWithSymbol != "" {
		t.Errorf("fetcher should not be called on validation error, got calledWithSymbol=%q", fake.calledWithSymbol)
	}
}

func TestGetKlineHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeKlineFetcher{err: wantErr}
	handler := getKlineHandler(fake)

	_, _, err := handler(context.Background(), nil, GetKlineInput{Symbol: "AAPL", Period: "day"})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
