// Package service holds the implementation of the kolide service interface and the HTTP endpoints
// for the API
package service

import (
	"io"

	"github.com/WatchBeam/clock"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/kolide/server/config"
	"github.com/kolide/kolide/server/kolide"
	"golang.org/x/net/context"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// NewService creates a new service from the config struct
func NewService(ds kolide.Datastore, resultStore kolide.QueryResultStore, logger kitlog.Logger, kolideConfig config.KolideConfig, mailService kolide.MailService, c clock.Clock, checker LicenseChecker) (kolide.Service, error) {
	var svc kolide.Service

	logFile := func(path string) io.Writer {
		return &lumberjack.Logger{
			Filename:   path,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
	}

	svc = service{
		ds:             ds,
		resultStore:    resultStore,
		logger:         logger,
		config:         kolideConfig,
		clock:          c,
		licenseChecker: checker,

		osqueryStatusLogWriter: logFile(kolideConfig.Osquery.StatusLogFile),
		osqueryResultLogWriter: logFile(kolideConfig.Osquery.ResultLogFile),
		mailService:            mailService,
	}
	svc = validationMiddleware{svc, ds}
	return svc, nil
}

type service struct {
	ds             kolide.Datastore
	resultStore    kolide.QueryResultStore
	logger         kitlog.Logger
	config         config.KolideConfig
	clock          clock.Clock
	licenseChecker LicenseChecker

	osqueryStatusLogWriter io.Writer
	osqueryResultLogWriter io.Writer

	mailService kolide.MailService
}

func (s service) SendEmail(mail kolide.Email) error {
	return s.mailService.SendEmail(mail)
}

func (s service) Clock() clock.Clock {
	return s.clock
}

// LicenseChecker allows checking that a license is valid by calling in to
// a remote URL.
type LicenseChecker interface {
	CheckinLicense(ctx context.Context)
}
