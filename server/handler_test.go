package server

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kolide/kolide-ose/datastore"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestAPIRoutes(t *testing.T) {
	ds, err := datastore.New("gorm-sqlite3", ":memory:")
	assert.Nil(t, err)

	svc, err := NewService(testConfig(ds))
	assert.Nil(t, err)

	ctx := context.Background()

	r := mux.NewRouter()
	attachKolideRoutes(r, ctx, svc, nil)
	handler := mux.NewRouter()
	handler.PathPrefix("/api/v1/kolide").Handler(r)

	var routes = []struct {
		verb string
		uri  string
	}{
		{
			verb: "POST",
			uri:  "/api/v1/kolide/users",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/users/1",
		},
		{
			verb: "POST",
			uri:  "/api/v1/kolide/users/1/password",
		},
		{
			verb: "POST",
			uri:  "/api/v1/kolide/users/1/role",
		},
		{
			verb: "POST",
			uri:  "/api/v1/kolide/users/1/status",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/queries/1",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/queries",
		},
		{
			verb: "POST",
			uri:  "/api/v1/kolide/queries",
		},
		{
			verb: "PATCH",
			uri:  "/api/v1/kolide/queries/1",
		},
		{
			verb: "DELETE",
			uri:  "/api/v1/kolide/queries/1",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/packs/1",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/packs",
		},
		{
			verb: "POST",
			uri:  "/api/v1/kolide/packs",
		},
		{
			verb: "PATCH",
			uri:  "/api/v1/kolide/packs/1",
		},
		{
			verb: "DELETE",
			uri:  "/api/v1/kolide/packs/1",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/packs/1/queries/2",
		},
		{
			verb: "GET",
			uri:  "/api/v1/kolide/packs/1/queries",
		},
		{
			verb: "DELETE",
			uri:  "/api/v1/kolide/packs/1/queries/2",
		},
	}

	for _, route := range routes {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(
			recorder,
			httptest.NewRequest(route.verb, route.uri, nil),
		)
		assert.NotEqual(t, 404, recorder.Code)
	}

}
