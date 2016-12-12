package service

import (
	"testing"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestCreateAppConfig(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)
	svc, err := newTestService(ds, nil)
	require.Nil(t, err)
	var appConfigTests = []struct {
		configPayload kolide.AppConfigPayload
	}{
		{
			configPayload: kolide.AppConfigPayload{
				OrgInfo: &kolide.OrgInfo{
					OrgLogoURL: stringPtr("acme.co/images/logo.png"),
					OrgName:    stringPtr("Acme"),
				},
				ServerSettings: &kolide.ServerSettings{
					KolideServerURL: stringPtr("https://acme.co:8080/"),
				},
			},
		},
	}

	for _, tt := range appConfigTests {
		result, err := svc.NewAppConfig(context.Background(), tt.configPayload)
		require.Nil(t, err)

		payload := tt.configPayload
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, *payload.OrgInfo.OrgLogoURL, result.OrgLogoURL)
		assert.Equal(t, *payload.OrgInfo.OrgName, result.OrgName)
		assert.Equal(t, *payload.ServerSettings.KolideServerURL, result.KolideServerURL)
	}
}

func TestVerfiyNoAuthSMTPConnection(t *testing.T) {
	config := &kolide.AppConfig{
		SMTPConfig: kolide.SMTPConfig{
			SenderAddress:      "foo@bar.com",
			Server:             "localhost",
			Port:               1025,
			AuthenticationType: kolide.AuthTypeNone,
		},
	}

	response := testSMTPConfiguration("bobdobbs@subgenuis.com", config)
	require.NotNil(t, response)
	assert.True(t, response.Success)

}
