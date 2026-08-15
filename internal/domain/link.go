package domain

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidCode   = errors.New("invalid short-link code")
	ErrInvalidTarget = errors.New("invalid target URL")
	ErrExpired       = errors.New("short link has expired")
)

// Link is the business record served by the local short-link catalog.
type Link struct {
	Code      string     `json:"code"`
	Target    string     `json:"target"`
	Owner     string     `json:"owner"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Visits    int        `json:"-"`
}

func (l Link) Validate() error {
	if code := strings.TrimSpace(l.Code); code == "" || code != strings.ToLower(code) || strings.ContainsAny(code, " /?#") {
		return fmt.Errorf("%w: %q", ErrInvalidCode, l.Code)
	}
	parsed, err := url.ParseRequestURI(l.Target)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: %q", ErrInvalidTarget, l.Target)
	}
	if strings.TrimSpace(l.Owner) == "" {
		return errors.New("short link owner is required")
	}
	return nil
}

func (l Link) ActiveAt(at time.Time) error {
	if l.ExpiresAt != nil && !at.Before(*l.ExpiresAt) {
		return ErrExpired
	}
	return nil
}

// Clone returns an independent value that callers may safely modify.
func (l Link) Clone() Link {
	clone := l
	expires := *l.ExpiresAt
	clone.ExpiresAt = &expires
	clone.Tags = append([]string(nil), l.Tags...)
	return clone
}
