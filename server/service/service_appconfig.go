package service

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/contexts/viewer"
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, err
	}
	newConfig, err := svc.ds.NewAppConfig(fromPayload(p, *config))
	if err != nil {
		return nil, err
	}
	return newConfig, nil
}

func (svc service) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	return svc.ds.AppConfig()
}

func (svc service) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	oldConfig, err := svc.AppConfig(ctx)
	if err != nil {
		return nil, err
	}
	newConfig := fromPayload(p, *oldConfig)
	if p.SMTPSettings != nil && !p.SMTPSettings.SMTPDisabled {
		oldSettings := smtpSettingsFromAppConfig(oldConfig)
		// anything changed?
		if *oldSettings != *p.SMTPSettings {
			v, ok := viewer.FromContext(ctx)
			if !ok {
				return nil, errors.InternalServerError(fmt.Errorf("unable to retrieve viewer from context"))
			}
			testMail := kolide.Email{
				Subject: "Hello from Kolide",
				To:      []string{v.User.Email},
				Mailer: &kolide.SMTPTestMailer{
					KolideServerURL: newConfig.KolideServerURL,
				},
			}
			// test mail set SMTPConfigured so we know if we can send mail
			err = svc.mailService.SendEmail(testMail)
			if err != nil {
				// if the provided SMTP parameters don't work with the targeted SMTP server
				// capture the error and return it to the front end so that GUI can
				// display the problem to the end user to aid in diagnosis
				newConfig.SMTPLastError = err.Error()
			}
			newConfig.SMTPConfigured = (err == nil)

			// if testing is indicated we don't persist anything, otherwise
			// email is marked as unconfigured
			if p.SMTPTest != nil && *p.SMTPTest {
				return newConfig, nil
			}
		}
	}
	if err = svc.ds.SaveAppConfig(newConfig); err != nil {
		return nil, err
	}
	return newConfig, nil
}

func fromPayload(p kolide.AppConfigPayload, config kolide.AppConfig) *kolide.AppConfig {
	if p.OrgInfo != nil && p.OrgInfo.OrgLogoURL != nil {
		config.OrgLogoURL = *p.OrgInfo.OrgLogoURL
	}
	if p.OrgInfo != nil && p.OrgInfo.OrgName != nil {
		config.OrgName = *p.OrgInfo.OrgName
	}
	if p.ServerSettings != nil && p.ServerSettings.KolideServerURL != nil {
		config.KolideServerURL = *p.ServerSettings.KolideServerURL
	}
	if p.SMTPSettings != nil {
		config.SMTPAuthenticationMethod = p.SMTPSettings.SMTPAuthenticationMethod
		config.SMTPAuthenticationType = p.SMTPSettings.SMTPAuthenticationType
		config.SMTPConfigured = p.SMTPSettings.SMTPConfigured
		config.SMTPDisabled = p.SMTPSettings.SMTPDisabled
		config.SMTPDomain = p.SMTPSettings.SMTPDomain
		config.SMTPEnableStartTLS = p.SMTPSettings.SMTPEnableStartTLS
		config.SMTPEnableTLS = p.SMTPSettings.SMTPEnableTLS
		config.SMTPPassword = p.SMTPSettings.SMTPPassword
		config.SMTPPort = p.SMTPSettings.SMTPPort
		config.SMTPSenderAddress = p.SMTPSettings.SMTPSenderAddress
		config.SMTPServer = p.SMTPSettings.SMTPServer
		config.SMTPUserName = p.SMTPSettings.SMTPUserName
		config.SMTPVerifySSLCerts = p.SMTPSettings.SMTPVerifySSLCerts
	}
	return &config
}

func smtpSettingsFromAppConfig(config *kolide.AppConfig) *kolide.SMTPSettings {
	return &kolide.SMTPSettings{
		SMTPConfigured:           config.SMTPConfigured,
		SMTPSenderAddress:        config.SMTPSenderAddress,
		SMTPServer:               config.SMTPServer,
		SMTPPort:                 config.SMTPPort,
		SMTPAuthenticationType:   config.SMTPAuthenticationType,
		SMTPUserName:             config.SMTPUserName,
		SMTPPassword:             config.SMTPPassword,
		SMTPEnableTLS:            config.SMTPEnableTLS,
		SMTPAuthenticationMethod: config.SMTPAuthenticationMethod,
		SMTPDomain:               config.SMTPDomain,
		SMTPVerifySSLCerts:       config.SMTPVerifySSLCerts,
		SMTPEnableStartTLS:       config.SMTPEnableStartTLS,
		SMTPDisabled:             config.SMTPDisabled,
	}
}

func fromAppConfig(config *kolide.AppConfig) *kolide.AppConfigPayload {
	return &kolide.AppConfigPayload{
		OrgInfo: &kolide.OrgInfo{
			OrgLogoURL: &config.OrgLogoURL,
			OrgName:    &config.OrgName,
		},
		ServerSettings: &kolide.ServerSettings{
			KolideServerURL: &config.KolideServerURL,
		},
		SMTPSettings: smtpSettingsFromAppConfig(config),
	}
}
