package kolide

import (
	"time"

	"golang.org/x/net/context"
)

type QueryStore interface {
	// Query methods
	NewQuery(query *Query) (*Query, error)
	SaveQuery(query *Query) error
	DeleteQuery(query *Query) error
	Query(id uint) (*Query, error)
	ListQueries(opt ListOptions) ([]*Query, error)
	// LoadPacksForQueries loads the packs associated with the provided
	// queries.
	LoadPacksForQueries(queries []*Query) error
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
	// Packs is loaded via LoadPacksForQueries (requires a join in the
	// MySQL backend)
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
