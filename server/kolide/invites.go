package kolide

import (
	"bytes"
	"html/template"

	"golang.org/x/net/context"
)

// InviteStore contains the methods for
// managing user invites in a datastore.
type InviteStore interface {
	// NewInvite creates and stores a new invitation in a DB.
	NewInvite(i *Invite) (*Invite, error)

	// Invites lists all invites in the datastore.
	ListInvites(opt ListOptions) ([]*Invite, error)

	// Invite retrieves an invite by it's ID.
	Invite(id uint) (*Invite, error)

	// InviteByEmail retrieves an invite for a specific email address.
	InviteByEmail(email string) (*Invite, error)

	// InviteByToken retrieves and invite using the token string.
	InviteByToken(token string) (*Invite, error)

	// SaveInvite saves an invitation in the datastore.
	SaveInvite(i *Invite) error
}

// InviteService contains methods for a service which deals with
// user invites.
type InviteService interface {
	// InviteNewUser creates an invite for a new user to join Kolide.
	InviteNewUser(ctx context.Context, payload InvitePayload) (invite *Invite, err error)

	// DeleteInvite removes an invite.
	DeleteInvite(ctx context.Context, id uint) (err error)

	// Invites returns a list of all invites.
	ListInvites(ctx context.Context, opt ListOptions) (invites []*Invite, err error)

	// VerifyInvite verifies that an invite exists and that it matches the
	// invite token.
	VerifyInvite(ctx context.Context, token string) (invite *Invite, err error)
}

// InvitePayload contains fields required to create a new user invite.
type InvitePayload struct {
	InvitedBy *uint `json:"invited_by"`
	Email     *string
	Admin     *bool
	Name      *string
	Position  *string
}

// Invite represents an invitation for a user to join Kolide.
type Invite struct {
	UpdateCreateTimestamps
	DeleteFields
	ID        uint   `json:"id"`
	InvitedBy uint   `json:"invited_by" db:"invited_by"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
	Name      string `json:"name"`
	Position  string `json:"position,omitempty"`
	Token     string `json:"-"`
}

func (i *Invite) EntityID() uint {
	return i.ID
}

func (i *Invite) EntityType() string {
	return "invites"
}

// InviteMailer is used to build an email template for the invite email.
type InviteMailer struct {
	*Invite
	KolideServerURL   template.URL
	InvitedByUsername string
}

func (i *InviteMailer) Message() ([]byte, error) {
	t, err := getTemplate("server/mail/templates/invite_token.html")
	if err != nil {
		return nil, err
	}

	var msg bytes.Buffer
	if err = t.Execute(&msg, i); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}
