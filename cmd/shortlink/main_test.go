package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunResolvesCatalogLink(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"--input", "../../examples/links.json", "--code", "docs", "--at", "2026-08-15T12:00:00Z"}, &output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "target=https://docs.example.com/getting-started") {
		t.Fatalf("output=%q", output.String())
	}
}

func TestRunListsOwnerLinks(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"--input", "../../examples/links.json", "--owner", "platform"}, &output)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(strings.TrimSpace(output.String()), "\n") + 1; got != 2 {
		t.Fatalf("lines=%d output=%q", got, output.String())
	}
}

func TestRunReturnsCancellationBeforeReadingInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	err := run(ctx, []string{"--input", "missing-catalog.json", "--code", "docs"}, &output)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want context cancellation", err)
	}
}
