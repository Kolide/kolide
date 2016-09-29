package datastore

import (
	"sync"

	"github.com/kolide/kolide-ose/server/kolide"
)

type inmem struct {
	Driver               string
	mtx                  sync.RWMutex
	users                map[uint]*kolide.User
	sessions             map[uint]*kolide.Session
	passwordResets       map[uint]*kolide.PasswordResetRequest
	invites              map[uint]*kolide.Invite
	labels               map[uint]*kolide.Label
	labelQueryExecutions map[uint]*kolide.LabelQueryExecution
	queries              map[uint]*kolide.Query
	orginfo              *kolide.OrgInfo
}

func (orm *inmem) Name() string {
	return "inmem"
}

func (orm *inmem) Migrate() error {
	return nil
}

func (orm *inmem) Drop() error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()
	orm.users = make(map[uint]*kolide.User)
	orm.sessions = make(map[uint]*kolide.Session)
	orm.passwordResets = make(map[uint]*kolide.PasswordResetRequest)
	return nil
}
