package service

import (
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (mw loggingMiddleware) SaveLicense(ctx context.Context, jwtToken string) (*kolide.License, error) {
	var (
		lic *kolide.License
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SaveLicense",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	lic, err = mw.Service.SaveLicense(ctx, jwtToken)
	return lic, err
}

func (mw loggingMiddleware) License(ctx context.Context) (*kolide.License, error) {
	var (
		lic *kolide.License
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "License",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	lic, err = mw.Service.License(ctx)
	return lic, err

}
