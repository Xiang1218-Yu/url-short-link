package service

import (
	"context"
	"fmt"
	"time"

	"url-short-link/internal/domain"
	"url-short-link/internal/store"
)

// Resolver applies expiration policy before a target is exposed and a visit is counted.
type Resolver struct{ catalog store.Catalog }

func NewResolver(catalog store.Catalog) *Resolver { return &Resolver{catalog: catalog} }

func (r *Resolver) Resolve(ctx context.Context, code string, at time.Time) (domain.Resolution, error) {
	if err := ctx.Err(); err != nil {
		return domain.Resolution{}, err
	}
	link, err := r.catalog.Find(ctx, code)
	if err != nil {
		return domain.Resolution{}, fmt.Errorf("find short link: %w", err)
	}
	if err := link.ActiveAt(at); err != nil {
		return domain.Resolution{}, err
	}
	visits, err := r.catalog.RecordVisit(ctx, link.Code)
	if err != nil {
		return domain.Resolution{}, fmt.Errorf("record short-link visit: %w", err)
	}
	return domain.NewResolution(link, visits, at), nil
}
