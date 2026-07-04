package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakePositionsFetcher struct {
	positions []tiger.Position
	err       error
}

func (f *fakePositionsFetcher) Positions(_ context.Context) ([]tiger.Position, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.positions, nil
}

func TestGetPositionsHandler_HappyPath(t *testing.T) {
	fake := &fakePositionsFetcher{
		positions: []tiger.Position{{Symbol: "AAPL", PositionQty: 10, MarketValue: 2000}},
	}
	handler := getPositionsHandler(fake)

	_, output, err := handler(context.Background(), nil, GetPositionsInput{})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(output.Positions, fake.positions) {
		t.Errorf("output.Positions = %+v, want %+v", output.Positions, fake.positions)
	}
}

func TestGetPositionsHandler_EmptyPositionsIsNotAnError(t *testing.T) {
	fake := &fakePositionsFetcher{positions: []tiger.Position{}}
	handler := getPositionsHandler(fake)

	_, output, err := handler(context.Background(), nil, GetPositionsInput{})
	if err != nil {
		t.Fatalf("empty positions should not be an error, got: %v", err)
	}
	if len(output.Positions) != 0 {
		t.Errorf("output.Positions = %+v, want empty", output.Positions)
	}
}

func TestGetPositionsHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakePositionsFetcher{err: wantErr}
	handler := getPositionsHandler(fake)

	_, _, err := handler(context.Background(), nil, GetPositionsInput{})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
