package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/sso"
	"github.com/pkg/errors"
)

func (mw validationMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (string, error) {
	invalid := &invalidArgumentError{}
	session, err := mw.ssoSessionStore.Get(auth.RequestID())
	if err != nil {
		invalid.Append("session", "missing for request")
		return "", invalid
	}
	validator, err := sso.NewValidator(session.Metadata)
	if err != nil {
		return "", errors.Wrap(err, "creating validator from metadata")
	}
	// make sure the response hasn't been tampered with
	auth, err = validator.ValidateSignature(auth)
	if err != nil {
		invalid.Appendf("sso response", "signature validation failed %s", err.Error())
		return "", invalid
	}
	// make sure the response isn't stale
	err = validator.ValidateResponse(auth)
	if err != nil {
		invalid.Appendf("sso response", "response validation failed %s", err.Error())
	}

	return mw.Service.CallbackSSO(ctx, auth)
}
