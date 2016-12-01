package kolide

// NotFoundError is returned when the datastore resource cannot be found.
type NotFoundError interface {
	error
	IsNotFound() bool
}

// AlreadyExists is returned when creating a datastore resource that already
// exists.
type AlreadyExistsError interface {
	error
	IsExists() bool
}

// Datastore combines all the interfaces in the Kolide DAL
type Datastore interface {
	UserStore
	QueryStore
	CampaignStore
	PackStore
	LabelStore
	HostStore
	PasswordResetStore
	SessionStore
	AppConfigStore
	InviteStore
	Name() string
	Drop() error
	Migrate() error
}
