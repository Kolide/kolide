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
	Queries() ([]*Query, error)
}

type QueryService interface {
	GetAllQueries(ctx context.Context) ([]*Query, error)
	GetQuery(ctx context.Context, id uint) (*Query, error)
	NewQuery(ctx context.Context, p QueryPayload) (*Query, error)
	ModifyQuery(ctx context.Context, id uint, p QueryPayload) (*Query, error)
	DeleteQuery(ctx context.Context, id uint) error
}

type QueryPayload struct {
	Name         *string
	Query        *string
	Interval     *uint
	Snapshot     *bool
	Differential *bool
	Platform     *string
	Version      *string
}

type PackPayload struct {
	Name     *string
	Platform *string
}

type Query struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string `gorm:"not null;unique_index:idx_query_unique_name"`
	Query        string `gorm:"not null"`
	Interval     uint
	Snapshot     bool
	Differential bool
	Platform     string
	Version      string
}

type DistributedQueryStatus int

const (
	QueryRunning  DistributedQueryStatus = iota
	QueryComplete DistributedQueryStatus = iota
	QueryError    DistributedQueryStatus = iota
)

type DistributedQueryCampaign struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	QueryID     uint
	MaxDuration time.Duration
	Status      DistributedQueryStatus
	UserID      uint
}

type DistributedQueryCampaignTarget struct {
	ID                         uint `gorm:"primary_key"`
	Type                       TargetType
	DistributedQueryCampaignID uint
	TargetID                   uint
}

type DistributedQueryExecutionStatus int

const (
	ExecutionWaiting   DistributedQueryExecutionStatus = iota
	ExecutionRequested DistributedQueryExecutionStatus = iota
	ExecutionSucceeded DistributedQueryExecutionStatus = iota
	ExecutionFailed    DistributedQueryExecutionStatus = iota
)

type DistributedQueryExecution struct {
	HostID             uint
	DistributedQueryID uint
	Status             DistributedQueryExecutionStatus
	Error              string `gorm:"size:1024"`
	ExecutionDuration  time.Duration
}

type Option struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"not null;unique_index:idx_option_unique_key"`
	Value     string `gorm:"not null"`
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad     DecoratorType = iota
	DecoratorAlways   DecoratorType = iota
	DecoratorInterval DecoratorType = iota
)

type Decorator struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType `gorm:"not null"`
	Interval  int
	Query     string
}
