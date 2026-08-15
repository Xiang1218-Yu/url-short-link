package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"url-short-link/internal/domain"
	"url-short-link/internal/store"
)

func newCatalog(t *testing.T, links []domain.Link) *store.MemoryCatalog {
	t.Helper()
	catalog, err := store.NewMemoryCatalog(links)
	if err != nil {
		t.Fatal(err)
	}
	return catalog
}

func TestResolveCountsOnlyActiveLinkVisits(t *testing.T) {
	expires := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	catalog := newCatalog(t, []domain.Link{{Code: "docs", Target: "https://docs.example.com", Owner: "platform"}, {Code: "old", Target: "https://old.example.com", Owner: "platform", ExpiresAt: &expires}})
	resolver := NewResolver(catalog)
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	resolution, err := resolver.Resolve(context.Background(), " docs ", at)
	if err != nil {
		t.Fatal(err)
	}
	if resolution.VisitCount != 1 || resolution.Target != "https://docs.example.com" {
		t.Fatalf("unexpected resolution: %#v", resolution)
	}
	_, err = resolver.Resolve(context.Background(), "old", expires)
	if !errors.Is(err, domain.ErrExpired) {
		t.Fatalf("error=%v, want expired", err)
	}
	link, err := catalog.Find(context.Background(), "old")
	if err != nil {
		t.Fatal(err)
	}
	if link.Visits != 0 {
		t.Fatalf("expired link visits=%d, want 0", link.Visits)
	}
}

func TestResolveHonorsCancelledContext(t *testing.T) {
	catalog := newCatalog(t, []domain.Link{{Code: "docs", Target: "https://docs.example.com", Owner: "platform"}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewResolver(catalog).Resolve(ctx, "docs", time.Now())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want canceled", err)
	}
}
