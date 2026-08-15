package domain

import "time"

// Resolution is the public result of a successful redirect lookup.
type Resolution struct {
	Code       string
	Target     string
	Owner      string
	Tags       []string
	VisitCount int
	ResolvedAt time.Time
}

func NewResolution(link Link, visitCount int, at time.Time) Resolution {
	return Resolution{
		Code: link.Code, Target: link.Target, Owner: link.Owner,
		Tags: append([]string(nil), link.Tags...), VisitCount: visitCount, ResolvedAt: at.UTC(),
	}
}
