package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yogyrahmawan/tiger-mcp/internal/tiger"
)

type fakeAssetsFetcher struct {
	assets []tiger.Asset
	err    error
}

func (f *fakeAssetsFetcher) Assets(_ context.Context) ([]tiger.Asset, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.assets, nil
}

func TestGetAccountAssetsHandler_HappyPath(t *testing.T) {
	fake := &fakeAssetsFetcher{
		assets: []tiger.Asset{{Account: "U123", Currency: "USD", NetLiquidation: 10000}},
	}
	handler := getAccountAssetsHandler(fake)

	_, output, err := handler(context.Background(), nil, GetAccountAssetsInput{})
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(output.Assets, fake.assets) {
		t.Errorf("output.Assets = %+v, want %+v", output.Assets, fake.assets)
	}
}

func TestGetAccountAssetsHandler_UpstreamError(t *testing.T) {
	wantErr := errors.New("tiger api unavailable")
	fake := &fakeAssetsFetcher{err: wantErr}
	handler := getAccountAssetsHandler(fake)

	_, _, err := handler(context.Background(), nil, GetAccountAssetsInput{})
	if err == nil {
		t.Fatal("expected error to be surfaced, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want it to wrap %v", err, wantErr)
	}
}
