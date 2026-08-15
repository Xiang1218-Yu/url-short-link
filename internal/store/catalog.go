package store

import (
	"context"
	"errors"

	"github.com/Xiang1218-Yu/url-short-link/internal/domain"
)

var (
	ErrNotFound      = errors.New("short link not found")
	ErrDuplicateCode = errors.New("duplicate short-link code")
)

// Catalog is the service boundary for reading links and recording successful resolutions.
type Catalog interface {
	Find(context.Context, string) (domain.Link, error)
	List(context.Context) ([]domain.Link, error)
	RecordVisit(context.Context, string) (int, error)
}
