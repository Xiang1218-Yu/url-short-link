package domain

import (
	"errors"
	"testing"
)

func TestValidateRejectsBlankOwner(t *testing.T) {
	link := Link{Code: "docs", Target: "https://docs.example.com"}
	if err := link.Validate(); !errors.Is(err, ErrMissingOwner) {
		t.Fatalf("Validate blank owner: error=%v, want ErrMissingOwner", err)
	}
}

func TestValidateRejectsWhitespaceOwner(t *testing.T) {
	link := Link{Code: "docs", Target: "https://docs.example.com", Owner: "  \t "}
	if err := link.Validate(); !errors.Is(err, ErrMissingOwner) {
		t.Fatalf("Validate whitespace owner: error=%v, want ErrMissingOwner", err)
	}
}

func TestValidateAcceptsCompleteLink(t *testing.T) {
	link := Link{Code: "docs", Target: "https://docs.example.com", Owner: "platform"}
	if err := link.Validate(); err != nil {
		t.Fatalf("Validate complete link: %v", err)
	}
}
