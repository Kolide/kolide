package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kolide/kolide/server/sso"
	"github.com/pkg/errors"
	"github.com/y0ssar1an/q"
)

func decodeGetInfoAboutSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionRequest{ID: id}, nil
}

func decodeGetInfoAboutSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getInfoAboutSessionsForUserRequest{ID: id}, nil
}

func decodeDeleteSessionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionRequest{ID: id}, nil
}

func decodeDeleteSessionsForUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteSessionsForUserRequest{ID: id}, nil
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Username = strings.ToLower(req.Username)
	return req, nil
}

func decodeInitiateSSORequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req initiateSSORequest
	err := json.NewDecoder(r.Body).Decode(&req)
	q.Q("decode", err, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginSSORequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCallbackSSORequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authResponse, err := sso.DecodeAuthResponse(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decoding sso callback")
	}
	return authResponse, nil
}
