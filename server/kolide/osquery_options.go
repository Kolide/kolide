package kolide

import "encoding/json"

type OsqueryOptionsStore interface {
	ApplyOptions(spec *OptionsSpec) error
	GetOptions() (*OptionsSpec, error)
	OptionsForHost(host *Host) (json.RawMessage, error)
}

type OsqueryOptionsService interface{}

type OptionsYaml struct {
	ApiVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Spec       OptionsSpec `json:"spec"`
}

type OptionsSpec struct {
	Config    json.RawMessage  `json:"config"`
	Overrides OptionsOverrides `json:"overrides"`
}

type OptionsOverrides struct {
	Platforms map[string]json.RawMessage `json:"platforms"`
}

type OptionsRow struct {
	ID                 int                `db:"id"`
	OverrideType       OptionOverrideType `db:"override_type"`
	OverrideIdentifier string             `db:"override_identifier"`
	Options            string             `db:"options"`
}

// OptionOverrideType is used to designate which override type a given set of
// options is used for. Currently the only supported override type is by
// platform.
type OptionOverrideType int

const (
	// OptionOverrideTypeDefault indicates that this is the default config
	// (provided to hosts when there is no override set for them).
	OptionOverrideTypeDefault OptionOverrideType = iota
	// OptionOverrideTypePlatform indicates that this is a
	// platform-specific config override (with precedence over the default
	// config).
	OptionOverrideTypePlatform
)
