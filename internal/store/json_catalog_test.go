package store

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"url-short-link/internal/domain"
)

func writeCatalog(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadJSONAcceptsCompleteCatalog(t *testing.T) {
	path := writeCatalog(t, "links.json", `{"links":[
		{"code":"docs","target":"https://docs.example.com","owner":"platform"},
		{"code":"status","target":"https://status.example.com","owner":"operations"}
	]}`)
	links, err := LoadJSON(context.Background(), path)
	if err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("links=%d, want 2", len(links))
	}
}

func TestLoadJSONRejectsBlankOwner(t *testing.T) {
	cases := map[string]string{
		"empty":  `{"links":[{"code":"docs","target":"https://docs.example.com","owner":""}]}`,
		"spaces": `{"links":[{"code":"docs","target":"https://docs.example.com","owner":"   "}]}`,
		"tab":    `{"links":[{"code":"docs","target":"https://docs.example.com","owner":"\t"}]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			path := writeCatalog(t, "links.json", body)
			if _, err := LoadJSON(context.Background(), path); err == nil {
				t.Fatal("LoadJSON accepted a link without a valid owner")
			} else if !errors.Is(err, domain.ErrMissingOwner) {
				t.Fatalf("error=%v, want ErrMissingOwner", err)
			}
		})
	}
}

func TestLoadJSONRejectsBlankOwnerAmongCompleteLinks(t *testing.T) {
	path := writeCatalog(t, "links.json", `{"links":[
		{"code":"docs","target":"https://docs.example.com","owner":"platform"},
		{"code":"launch","target":"https://example.com/launch","owner":"  "}
	]}`)
	if _, err := LoadJSON(context.Background(), path); err == nil {
		t.Fatal("LoadJSON accepted a catalog containing an ownerless link")
	} else if !errors.Is(err, domain.ErrMissingOwner) {
		t.Fatalf("error=%v, want ErrMissingOwner", err)
	}
}
