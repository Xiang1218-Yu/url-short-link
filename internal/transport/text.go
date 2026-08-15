package transport

import (
	"fmt"
	"io"
	"strings"

	"github.com/Xiang1218-Yu/url-short-link/internal/domain"
)

func WriteResolution(out io.Writer, resolution domain.Resolution) error {
	_, err := fmt.Fprintf(out, "code=%s\ntarget=%s\nowner=%s\ntags=%s\nvisits=%d\n", resolution.Code, resolution.Target, resolution.Owner, strings.Join(resolution.Tags, ","), resolution.VisitCount)
	return err
}

func WriteOwnerLinks(out io.Writer, links []domain.Link) error {
	for _, link := range links {
		if _, err := fmt.Fprintf(out, "%s\t%s\t%s\n", link.Code, link.Owner, link.Target); err != nil {
			return err
		}
	}
	return nil
}
