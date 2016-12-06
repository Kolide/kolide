package kolide

import (
	"time"

	"golang.org/x/net/context"
)

type QueryStore interface {
	// NewQuery creates a new query object in thie datastore. The returned
	// query should have the ID updated.
	NewQuery(query *Query) (*Query, error)
	// SaveQuery saves changes to an existing query object.
	SaveQuery(query *Query) error
	// DeleteQuery (soft) deletes an existing query object.
	DeleteQuery(query *Query) error
	// Query returns the query associated with the provided ID. Associated
	// packs should also be loaded.
	Query(id uint) (*Query, error)
	// ListQueries returns a list of queries with the provided sorting and
	// paging options. Associated packs should also be loaded.
	ListQueries(opt ListOptions) ([]*Query, error)
}

type QueryService interface {
	ListQueries(ctx context.Context, opt ListOptions) ([]*Query, error)
	GetQuery(ctx context.Context, id uint) (*Query, error)
	NewQuery(ctx context.Context, p QueryPayload) (*Query, error)
	ModifyQuery(ctx context.Context, id uint, p QueryPayload) (*Query, error)
	DeleteQuery(ctx context.Context, id uint) error
}

type QueryPayload struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Query        *string `json:"query"`
	Interval     *uint   `json:"interval"`
	Snapshot     *bool   `json:"snapshot"`
	Differential *bool   `json:"differential"`
	Platform     *string `json:"platform"`
	Version      *string `json:"version"`
}

type Query struct {
	UpdateCreateTimestamps
	DeleteFields
	ID           uint   `json:"id"`
	Saved        bool   `json:"saved"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Query        string `json:"query"`
	Interval     uint   `json:"interval"`
	Snapshot     bool   `json:"snapshot"`
	Differential bool   `json:"differential"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
	// Packs is loaded when retrieving queries, but is stored in a join
	// table in the MySQL backend.
	Packs []Pack `json:"packs" db:"-"`
}

type Option struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad DecoratorType = iota
	DecoratorAlways
	DecoratorInterval
)

type Decorator struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType
	Interval  int
	Query     string
}
