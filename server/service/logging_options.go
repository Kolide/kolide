package service

import (
	"context"
	"time"

	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
)

func (mw loggingMiddleware) GetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.loggerForError(err).Log(
			"method", "GetOptions",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.GetOptions(ctx)
	return options, err
}

func (mw loggingMiddleware) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	var (
		options      []kolide.Option
		err          error
		loggedInUser = "unauthenticated"
	)

	if vc, ok := viewer.FromContext(ctx); ok {

		loggedInUser = vc.Username()
	}

	defer func(begin time.Time) {
		_ = mw.loggerForError(err).Log(
			"method", "ModifyOptions",
			"err", err,
			"user", loggedInUser,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.ModifyOptions(ctx, req)
	return options, err
}

func (mw loggingMiddleware) ResetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)
	defer func(begin time.Time) {
		_ = mw.loggerForError(err).Log(
			"method", "ResetOptions",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.ResetOptions(ctx)
	return options, err
}
