package store

import (
	"context"
	"errors"
	"testing"

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
