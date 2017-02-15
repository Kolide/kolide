package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/WatchBeam/clock"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kolide/kolide/server/config"
	"github.com/kolide/kolide/server/datastore/mysql"
	"github.com/kolide/kolide/server/kolide"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type testResource struct {
	server     *httptest.Server
	adminToken string
	userToken  string
	ds         kolide.Datastore
}

type endpointService struct {
	kolide.Service
}

func (svc endpointService) SendTestEmail(ctx context.Context, config *kolide.AppConfig) error {
	return nil
}
func setupEndpointTest(t *testing.T) *testResource {
	test := &testResource{}

	var err error
	test.ds, err = mysql.New(config.TestConfig(), clock.NewMockClock())
	require.Nil(t, err)
	require.Nil(t, test.ds.MigrateData())

	devOrgInfo := &kolide.AppConfig{
		OrgName:                "Kolide",
		OrgLogoURL:             "http://foo.bar/image.png",
		SMTPPort:               465,
		SMTPAuthenticationType: kolide.AuthTypeUserNamePassword,
		SMTPEnableTLS:          true,
		SMTPVerifySSLCerts:     true,
		SMTPEnableStartTLS:     true,
	}
	test.ds.NewAppConfig(devOrgInfo)
	svc, _ := newTestService(test.ds, nil)
	svc = endpointService{svc}
	createTestUsers(t, test.ds)
	logger := kitlog.NewLogfmtLogger(os.Stdout)

	jwtKey := "CHANGEME"
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(setRequestsContexts(svc, jwtKey)),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerAfter(kithttp.SetContentType("application/json; charset=utf-8")),
	}

	router := mux.NewRouter()
	ke := MakeKolideServerEndpoints(svc, jwtKey)
	ctxt := context.Background()
	kh := makeKolideKitHandlers(ctxt, ke, opts)
	attachKolideAPIRoutes(router, kh)

	test.server = httptest.NewServer(router)

	userParam := loginRequest{
		Username: "admin1",
		Password: testUsers["admin1"].PlaintextPassword,
	}

	marshalledUser, _ := json.Marshal(&userParam)

	requestBody := &nopCloser{bytes.NewBuffer(marshalledUser)}
	resp, _ := http.Post(test.server.URL+"/api/v1/kolide/login", "application/json", requestBody)

	var jsn = struct {
		User  *kolide.User `json:"user"`
		Token string       `json:"token"`
		Err   string       `json:"error,omitempty"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsn)
	test.adminToken = jsn.Token

	// log in non admin user
	userParam.Username = "user1"
	userParam.Password = testUsers["user1"].PlaintextPassword
	marshalledUser, _ = json.Marshal(userParam)
	requestBody = &nopCloser{bytes.NewBuffer(marshalledUser)}
	resp, err = http.Post(test.server.URL+"/api/v1/kolide/login", "application/json", requestBody)
	require.Nil(t, err)
	err = json.NewDecoder(resp.Body).Decode(&jsn)
	require.Nil(t, err)
	test.userToken = jsn.Token

	return test
}

func functionName(f func(*testing.T, *testResource)) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	elements := strings.Split(fullName, ".")
	return elements[len(elements)-1]
}

var testFunctions = [...]func(*testing.T, *testResource){
	testGetAppConfig,
	testModifyAppConfig,
	testModifyAppConfigWithValidationFail,
	testGetOptions,
	testModifyOptions,
	testModifyOptionsValidationFail,
	testImportConfig,
	testImportConfigMissingExternal,
	testImportConfigWithMissingGlob,
	testImportConfigWithGlob,
	testAdminUserSetAdmin,
	testNonAdminUserSetAdmin,
	testAdminUserSetEnabled,
	testNonAdminUserSetEnabled,
}

func TestEndpoints(t *testing.T) {
	for _, f := range testFunctions {
		r := setupEndpointTest(t)
		defer r.server.Close()
		t.Run(functionName(f), func(t *testing.T) {
			f(t, r)
		})
	}
}
