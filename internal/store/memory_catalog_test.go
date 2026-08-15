package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"url-short-link/internal/domain"
)

func TestMemoryCatalogReturnsIndependentSnapshots(t *testing.T) {
	catalog, err := NewMemoryCatalog([]domain.Link{{Code: "docs", Target: "https://docs.example.com", Owner: "platform", Tags: []string{"docs"}}})
	if err != nil {
		t.Fatal(err)
	}
	first, err := catalog.Find(context.Background(), "docs")
	if err != nil {
		t.Fatal(err)
	}
	first.Tags[0] = "changed"
	second, err := catalog.Find(context.Background(), "docs")
	if err != nil {
		t.Fatal(err)
	}
	if second.Tags[0] != "docs" {
		t.Fatalf("catalog leaked caller mutation: %#v", second.Tags)
	}
}

func TestMemoryCatalogRejectsDuplicateCodes(t *testing.T) {
	_, err := NewMemoryCatalog([]domain.Link{
		{Code: "docs", Target: "https://docs.example.com", Owner: "platform"},
		{Code: "docs", Target: "https://other.example.com", Owner: "platform"},
	})
	if !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("error=%v, want duplicate code", err)
	}
}

func TestCancelledRecordDoesNotBlockLaterVisit(t *testing.T) {
	catalog, err := NewMemoryCatalog([]domain.Link{{Code: "docs", Target: "https://docs.example.com", Owner: "platform"}})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := catalog.RecordVisit(ctx, "docs"); !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
	done := make(chan error, 1)
	go func() { _, err := catalog.RecordVisit(context.Background(), "docs"); done <- err }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(150 * time.Millisecond):
		t.Fatal("a canceled record blocked the next visit")
	}
}
