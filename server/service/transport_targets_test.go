package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestDecodeSearchTargetsRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/targets", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeSearchTargetsRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(searchTargetsRequest)
		assert.Equal(t, "bar", params.Query)
	}).Methods("POST")
	var body bytes.Buffer

	body.Write([]byte(`{
        "query": "bar"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/targets", &body),
	)
}
