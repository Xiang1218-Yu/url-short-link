package main

import (
	"bytes"
	"context"
	"os"
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

func TestRunRejectsCatalogLinkWithoutOwner(t *testing.T) {
	file := t.TempDir() + "/links.json"
	if err := os.WriteFile(file, []byte(`{"links":[{"code":"docs","target":"https://docs.example.com","owner":"  "}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := run(context.Background(), []string{"--input", file, "--code", "docs", "--at", "2026-08-15T12:00:00Z"}, &output); err == nil {
		t.Fatal("catalog link without an owner was accepted")
	}
}
