package service

import (
	"context"
	"sort"
	"strings"

	"github.com/Xiang1218-Yu/url-short-link/internal/domain"
	"github.com/Xiang1218-Yu/url-short-link/internal/store"
)

// Reporter produces owner-scoped management views without exposing catalog internals.
type Reporter struct{ catalog store.Catalog }

func NewReporter(catalog store.Catalog) *Reporter { return &Reporter{catalog: catalog} }

func (r *Reporter) LinksForOwner(ctx context.Context, owner string) ([]domain.Link, error) {
	links, err := r.catalog.List(ctx)
	if err != nil {
		return nil, err
	}
	selected := make([]domain.Link, 0, len(links))
	for _, link := range links {
		if strings.EqualFold(link.Owner, strings.TrimSpace(owner)) {
			selected = append(selected, link.Clone())
		}
	}
	sort.Slice(selected, func(i, j int) bool { return selected[i].Code < selected[j].Code })
	return selected, nil
}
