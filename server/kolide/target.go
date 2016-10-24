package kolide

import (
	"golang.org/x/net/context"
)

type TargetSearchResults struct {
    Hosts []Host
    Labels []Label
}

type TargetService interface {
    Search(ctx context.Context, query string, omit []Target) (TargetSearchResults, uint, error)
    CountHostsInTargets(ctx context.Context, targets []Target) (uint, error)
}

type TargetType int

const (
	TargetLabel TargetType = iota
	TargetHost
)

type Target struct {
    Type TargetType
    TargetID uint
}