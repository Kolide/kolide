package kolide

import (
	"time"

	"golang.org/x/net/context"
)

type LabelStore interface {
	// Label methods
	NewLabel(Label *Label) (*Label, error)
	DeleteLabel(lid uint) error
	Label(lid uint) (*Label, error)
	ListLabels(opt ListOptions) ([]*Label, error)

	// LabelQueriesForHost returns the label queries that should be executed
	// for the given host. The cutoff is the minimum timestamp a query
	// execution should have to be considered "fresh". Executions that are
	// not fresh will be repeated. Results are returned in a map of label
	// id -> query
	LabelQueriesForHost(host *Host, cutoff time.Time) (map[string]string, error)

	// RecordLabelQueryExecutions saves the results of label queries. The
	// results map is a map of label id -> whether or not the label
	// matches. The time parameter is the timestamp to save with the query
	// execution.
	RecordLabelQueryExecutions(host *Host, results map[string]bool, t time.Time) error

	// LabelsForHost returns the labels that the given host is in.
	ListLabelsForHost(hid uint) ([]Label, error)

	ListHostsInLabel(lid uint) ([]Host, error)
	ListUniqueHostsInLabels(labels []uint) ([]Host, error)

	SearchLabels(query string, omit ...uint) ([]Label, error)
}

type LabelService interface {
	ListLabels(ctx context.Context, opt ListOptions) (labels []*Label, err error)
	GetLabel(ctx context.Context, id uint) (label *Label, err error)
	NewLabel(ctx context.Context, p LabelPayload) (label *Label, err error)
	DeleteLabel(ctx context.Context, id uint) (err error)
}

type LabelPayload struct {
	Name        *string `json:"name"`
	Query       *string `json:"query"`
	Platform    *string `json:"platform"`
	Description *string `json:"description"`
}

// LabelType is used to catagorize the kind of label
type LabelType uint

const (
	// LabelTypeRegular is for user created labels that can be modified.
	LabelTypeRegular LabelType = iota
	// LabelTypeBuiltIn is for labels built into Kolide that cannot be
	// modified by users.
	LabelTypeBuiltIn
)

type Label struct {
	UpdateCreateTimestamps
	DeleteFields
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Query       string    `json:"query"`
	Platform    string    `json:"platform"`
	LabelType   LabelType `json:"label_type" db:"label_type"`
	HostIDs     []uint    `json:"host_ids" db:"-"`
}

type LabelQueryExecution struct {
	ID        uint
	UpdatedAt time.Time
	Matches   bool
	LabelID   uint
	HostID    uint
}
