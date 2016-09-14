package server

import (
	"context"
	"testing"

	"github.com/drone/drone/store/datastore"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/stretchr/testify/assert"
)

// TestEndpointPermissions tests that
// the endpoint.Middleware correctly grants or denies
// permissions to access or modify resources
func TestEndpointPermissions(t *testing.T) {
	req := struct{}{}
	ds, _ := datastore.New("inmem", "")
	createTestUsers(t, ds)
	admin1, _ := ds.User("admin1")
	user1, _ := ds.User("user1")
	user2, _ := ds.User("user2")
	user2.Enabled = false

	e := endpoint.Nop // a test endpoint
	var endpointTests = []struct {
		endpoint endpoint.Endpoint
		// who is making the request
		vc *viewerContext
		// what resource are we editing
		requestID uint
		// what error to expect
		wantErr interface{}
		// custom request struct
		request interface{}
	}{
		{
			endpoint: mustBeAdmin(e),
			wantErr:  errNoContext,
		},
		{
			endpoint: canReadUser(e),
			wantErr:  errNoContext,
		},
		{
			endpoint: canModifyUser(e),
			wantErr:  errNoContext,
		},
		{
			endpoint: mustBeAdmin(e),
			vc:       &viewerContext{user: admin1},
		},
		{
			endpoint: mustBeAdmin(e),
			vc:       &viewerContext{user: user1},
			wantErr:  forbiddenError{"must be an admin"},
		},
		{
			endpoint: canModifyUser(e),
			vc:       &viewerContext{user: admin1},
		},
		{
			endpoint: canModifyUser(e),
			vc:       &viewerContext{user: user1},
			wantErr:  forbiddenError{message: "no write permissions on user"},
		},
		{
			endpoint:  canModifyUser(e),
			vc:        &viewerContext{user: user1},
			requestID: admin1.ID,
			wantErr:   forbiddenError{message: "no write permissions on user"},
		},
		{
			endpoint:  canReadUser(e),
			vc:        &viewerContext{user: user1},
			requestID: admin1.ID,
		},
		{
			endpoint:  canReadUser(e),
			vc:        &viewerContext{user: user2},
			requestID: admin1.ID,
			wantErr:   forbiddenError{message: "no read permissions on user"},
		},
		{
			endpoint: validateModifyUserRequest(e),
			request:  modifyUserRequest{},
			wantErr:  errNoContext,
		},
		{
			endpoint: validateModifyUserRequest(e),
			request:  modifyUserRequest{payload: kolide.UserPayload{Enabled: boolPtr(true)}},
			vc:       &viewerContext{user: user1},
			wantErr:  "must be an admin",
		},
		{
			endpoint: canResetPassword(e),
			vc:       emptyVC(),
			request:  changePasswordRequest{},
		},
		{
			endpoint:  canResetPassword(e),
			vc:        &viewerContext{user: user1},
			requestID: user1.ID,
			request:   changePasswordRequest{},
		},
		{
			endpoint: canResetPassword(e),
			vc:       &viewerContext{user: user2},
			request:  changePasswordRequest{},
			wantErr:  "must be logged in",
		},
	}

	for i, tt := range endpointTests {
		ctx := context.Background()
		if tt.vc != nil {
			ctx = context.WithValue(ctx, "viewerContext", tt.vc)
		}
		if tt.requestID != 0 {
			ctx = context.WithValue(ctx, "request-id", tt.requestID)
		}
		var request interface{}
		if tt.request != nil {
			request = tt.request
		} else {
			request = req
		}
		_, eerr := tt.endpoint(ctx, request)
		assert.Equal(t, tt.wantErr, eerr)
	}
}
