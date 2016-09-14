package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDecodeCreateUserRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/users", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeCreateUserRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(createUserRequest)
		assert.Equal(t, "foo", *params.payload.Name)
		assert.Equal(t, "foo@kolide.co", *params.payload.Email)
	}).Methods("POST")

	var body bytes.Buffer
	body.Write([]byte(`{
        "name": "foo",
        "email": "foo@kolide.co"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/users", &body),
	)
}

func TestDecodeGetUserRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/users/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetUserRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getUserRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/users/1", nil),
	)
}

func TestDecodeChangePasswordRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/users/{id}/password", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeChangePasswordRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(changePasswordRequest)
		assert.Equal(t, "foo", params.CurrentPassword)
		assert.Equal(t, "bar", params.NewPassword)
		assert.Equal(t, "baz", params.PasswordResetToken)
		assert.Equal(t, uint(1), params.UserID)
	}).Methods("POST")

	var body bytes.Buffer
	body.Write([]byte(`{
        "current_password": "foo",
        "new_password": "bar",
        "password_reset_token": "baz"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/users/1/password", &body),
	)
}

func TestDecodeModifyUserRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/users/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeModifyUserRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(modifyUserRequest)
		assert.Equal(t, "foo", *params.payload.Name)
		assert.Equal(t, "foo@kolide.co", *params.payload.Email)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("PATCH")

	var body bytes.Buffer
	body.Write([]byte(`{
        "name": "foo",
        "email": "foo@kolide.co"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("PATCH", "/api/v1/kolide/users/1", &body),
	)
}
