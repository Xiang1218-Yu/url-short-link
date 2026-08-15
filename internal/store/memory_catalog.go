package store

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"url-short-link/internal/domain"
)

// MemoryCatalog is the local runtime representation of the checked-in JSON catalog.
type MemoryCatalog struct {
	mu    sync.RWMutex
	links map[string]domain.Link
}

func NewMemoryCatalog(links []domain.Link) (*MemoryCatalog, error) {
	catalog := &MemoryCatalog{links: make(map[string]domain.Link, len(links))}
	for _, link := range links {
		if err := link.Validate(); err != nil {
			return nil, err
		}
		key := normalizeCode(link.Code)
		if _, exists := catalog.links[key]; exists {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateCode, key)
		}
		link.Code = key
		catalog.links[key] = link.Clone()
	}
	return catalog, nil
}

func (c *MemoryCatalog) Find(ctx context.Context, code string) (domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return domain.Link{}, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	link, ok := c.links[normalizeCode(code)]
	if !ok {
		return domain.Link{}, fmt.Errorf("%w: %s", ErrNotFound, code)
	}
	return link.Clone(), nil
}

func (c *MemoryCatalog) List(ctx context.Context) ([]domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	links := make([]domain.Link, 0, len(c.links))
	for _, link := range c.links {
		links = append(links, link.Clone())
	}
	sort.Slice(links, func(i, j int) bool { return links[i].Code < links[j].Code })
	return links, nil
}

func (c *MemoryCatalog) RecordVisit(ctx context.Context, code string) (int, error) {
	key := normalizeCode(code)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	link, ok := c.links[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, code)
	}
	link.Visits++
	c.links[key] = link
	return link.Visits, nil
}

func normalizeCode(code string) string { return strings.ToLower(strings.TrimSpace(code)) }
