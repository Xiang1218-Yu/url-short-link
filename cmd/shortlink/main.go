package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Xiang1218-Yu/url-short-link/internal/service"
	"github.com/Xiang1218-Yu/url-short-link/internal/store"
	"github.com/Xiang1218-Yu/url-short-link/internal/transport"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	flags := flag.NewFlagSet("shortlink", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	input := flags.String("input", "examples/links.json", "local JSON catalog")
	code := flags.String("code", "", "short-link code to resolve")
	owner := flags.String("owner", "", "owner whose links should be listed")
	atText := flags.String("at", "", "RFC3339 resolution time (optional)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if (*code == "" && *owner == "") || (*code != "" && *owner != "") {
		return errors.New("provide exactly one of --code or --owner")
	}
	links, err := store.LoadJSON(ctx, *input)
	if err != nil {
		return err
	}
	catalog, err := store.NewMemoryCatalog(links)
	if err != nil {
		return fmt.Errorf("prepare catalog: %w", err)
	}
	if *owner != "" {
		links, err := service.NewReporter(catalog).LinksForOwner(ctx, *owner)
		if err != nil {
			return err
		}
		return transport.WriteOwnerLinks(out, links)
	}
	at := time.Now().UTC()
	if *atText != "" {
		at, err = time.Parse(time.RFC3339, *atText)
		if err != nil {
			return fmt.Errorf("parse --at: %w", err)
		}
	}
	resolution, err := service.NewResolver(catalog).Resolve(ctx, *code, at)
	if err != nil {
		return err
	}
	return transport.WriteResolution(out, resolution)
}
