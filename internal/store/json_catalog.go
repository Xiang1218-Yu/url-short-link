package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"url-short-link/internal/domain"
)

type catalogFile struct {
	Links []domain.Link `json:"links"`
}

// LoadJSON reads only local catalog data; it performs no network or database access.
func LoadJSON(ctx context.Context, path string) ([]domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open catalog: %w", err)
	}
	defer file.Close()

	var payload catalogFile
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode catalog: %w", err)
	}
	for i, link := range payload.Links {
		if strings.TrimSpace(link.Owner) == "" {
			return nil, fmt.Errorf("link %d (%q): %w", i+1, link.Code, domain.ErrMissingOwner)
		}
	}
	if len(payload.Links) == 0 {
		return nil, fmt.Errorf("catalog contains no links")
	}
	return payload.Links, nil
}
