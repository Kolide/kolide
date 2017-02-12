package kolide

import "golang.org/x/net/context"

type OsqueryService interface {
	EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (nodeKey string, err error)
	AuthenticateHost(ctx context.Context, nodeKey string) (host *Host, err error)
	GetClientConfig(ctx context.Context) (config *OsqueryConfig, err error)
	GetDistributedQueries(ctx context.Context) (queries map[string]string, err error)
	SubmitDistributedQueryResults(ctx context.Context, results OsqueryDistributedQueryResults, statuses map[string]string) (err error)
	SubmitStatusLogs(ctx context.Context, logs []OsqueryStatusLog) (err error)
	SubmitResultLogs(ctx context.Context, logs []OsqueryResultLog) (err error)
}

type OsqueryDistributedQueryResults map[string][]map[string]string

type QueryContent struct {
	Query       string `json:"query"`
	Description string `json:"description,omitempty"`
	Interval    uint   `json:"interval"`
	Platform    string `json:"platform,omitempty"`
	Version     string `json:"version,omitempty"`
	Snapshot    bool   `json:"snapshot,omitempty"`
	Removed     bool   `json:"removed"`
	Shard       uint   `json:"shard"`
}

type Queries map[string]QueryContent

type PackContent struct {
	Platform  string   `json:"platform,omitempty"`
	Version   string   `json:"version,omitempty"`
	Shard     uint     `json:"shard"`
	Discovery []string `json:"discovery,omitempty"`
	Queries   Queries  `json:"queries"`
}

type Packs map[string]PackContent

type Decorators struct {
	Load     []string            `json:"load,omitempty"`
	Always   []string            `json:"always,omitempty"`
	Interval map[string][]string `json:"interval,omitempty"`
}

type OsqueryConfig struct {
	Options    map[string]interface{} `json:"options"`
	Decorators Decorators             `json:"decorators,omitempty"`
	Packs      Packs                  `json:"packs,omitempty"`
}

type OsqueryResultLog struct {
	Name           string `json:"name"`
	HostIdentifier string `json:"hostIdentifier"`
	UnixTime       string `json:"unixTime"`
	CalendarTime   string `json:"calendarTime"`
	// Columns stores the columns of differential queries
	Columns map[string]string `json:"columns,omitempty"`
	// Snapshot stores the rows and columns of snapshot queries
	Snapshot    []map[string]string `json:"snapshot,omitempty"`
	Action      string              `json:"action"`
	Decorations map[string]string   `json:"decorations"`
}

type OsqueryStatusLog struct {
	Severity    string            `json:"severity"`
	Filename    string            `json:"filename"`
	Line        string            `json:"line"`
	Message     string            `json:"message"`
	Version     string            `json:"version"`
	Decorations map[string]string `json:"decorations"`
}
