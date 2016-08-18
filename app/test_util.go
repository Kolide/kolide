package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type testLogger struct {
	t *testing.T
}

func (t *testLogger) Print(v ...interface{}) {
	t.t.Log(v...)
}

func (t *testLogger) Write(p []byte) (n int, err error) {
	t.t.Log(string(p))
	return len(p), nil
}

func openTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Error opening test DB: %s", err.Error()))
	}

	CreateTables(db)
	if db.Error != nil {
		panic(fmt.Sprintf("Error creating test DB tables: %s", db.Error.Error()))
	}

	// Log using t.Log so that output only shows up if the test fails
	db.SetLogger(&testLogger{t: t})
	db.LogMode(true)

	return db
}

type mockOsqueryStatusLogWriter struct {
	Logs []OsqueryStatusLog
}

func (w *mockOsqueryStatusLogWriter) HandleStatusLog(log OsqueryStatusLog, nodeKey string) error {
	w.Logs = append(w.Logs, log)
	return nil
}

type mockOsqueryResultLogWriter struct {
	Logs []OsqueryResultLog
}

func (w *mockOsqueryResultLogWriter) HandleResultLog(log OsqueryResultLog, nodeKey string) error {
	w.Logs = append(w.Logs, log)
	return nil
}

type IntegrationRequests struct {
	r             *gin.Engine
	db            *gorm.DB
	pool          SMTPConnectionPool
	t             *testing.T
	statusHandler *mockOsqueryStatusLogWriter
	resultHandler *mockOsqueryResultLogWriter
}

func (req *IntegrationRequests) New(t *testing.T) {
	// req.t = t
	//
	// req.db = openTestDB(t)
	// req.pool = newMockSMTPConnectionPool()
	//
	// // Until we have a better solution for first-user onboarding, manually
	// // create an admin
	// _, err := NewUser(req.db, "admin", "foobar", "admin@kolide.co", true, false)
	// if err != nil {
	// 	t.Fatalf("Error opening DB: %s", err.Error())
	// }
	//
	// req.r = CreateServer(req.db, req.pool, &testLogger{t: t})
}

func (req *IntegrationRequests) Login(username, password string, sessionOut *string) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(LoginRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/api/v1/kolide/login", buff)
	request.Header.Set("Content-Type", "application/json")
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return
	}
	*sessionOut = response.Header().Get("Set-Cookie")

	return
}

func (req *IntegrationRequests) CreateUser(username, password, email string, admin, reset bool, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(CreateUserRequestBody{
		Username:           username,
		Password:           password,
		Email:              email,
		Admin:              admin,
		NeedsPasswordReset: reset,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PUT", "/api/v1/kolide/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) GetUser(username, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(GetUserRequestBody{
		Username: username,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/api/v1/kolide/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) ModifyUser(username, name, email, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(ModifyUserRequestBody{
		Username: username,
		Name:     name,
		Email:    email,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/api/v1/kolide/user", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) ChangePassword(username, currentPassword, newPassword, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(ChangePasswordRequestBody{
		Username:          username,
		CurrentPassword:   currentPassword,
		NewPassword:       newPassword,
		NewPasswordConfim: newPassword,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/api/v1/kolide/user/password", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
	}

	return &responseBody
}

func (req *IntegrationRequests) GetUserSessionInfo(username, session string) *GetInfoAboutSessionsForUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(GetInfoAboutSessionsForUserRequestBody{
		Username: username,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/api/v1/kolide/user/sessions", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetInfoAboutSessionsForUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) DeleteUserSessions(username, session string) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(GetInfoAboutSessionsForUserRequestBody{
		Username: username,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("DELETE", "/api/v1/kolide/user/sessions", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return
	}

	return
}

func (req *IntegrationRequests) GetSessionInfo(sessionKey, session string) *SessionInfoResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(GetInfoAboutSessionRequestBody{
		SessionKey: sessionKey,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/api/v1/kolide/session", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody SessionInfoResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) DeleteSession(sessionID uint, session string) {
	response := httptest.NewRecorder()
	body, err := json.Marshal(DeleteSessionRequestBody{
		SessionID: sessionID,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("DELETE", "/api/v1/kolide/session", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return
	}

	return
}

func (req *IntegrationRequests) SetAdminState(username string, admin bool, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(SetUserAdminStateRequestBody{
		Username: username,
		Admin:    admin,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/api/v1/kolide/user/admin", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
	}

	return &responseBody
}

func (req *IntegrationRequests) SetEnabledState(username string, enabled bool, session string) *GetUserResponseBody {
	response := httptest.NewRecorder()
	body, err := json.Marshal(SetUserEnabledStateRequestBody{
		Username: username,
		Enabled:  enabled,
	})
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("PATCH", "/api/v1/kolide/user/enabled", buff)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", session)
	req.r.ServeHTTP(response, request)

	if response.Code != 200 {
		req.t.Fatalf("Response code: %d", response.Code)
		return nil
	}

	var responseBody GetUserResponseBody
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		req.t.Fatal(err.Error())
		return nil
	}

	return &responseBody
}

func (req *IntegrationRequests) CheckUser(username, email, name string, admin, reset, enabled bool) {
	var user User
	err := req.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		req.t.Fatal(err.Error())
		return
	}
	if user.Email != email {
		req.t.Fatalf("user's email was not set in the DB: %s", user.Email)
	}
	if user.Admin != admin {
		req.t.Fatal("user admin settings don't match")
	}
	if user.NeedsPasswordReset != reset {
		req.t.Fatal("user reset settings don't match")
	}
	if user.Enabled != enabled {
		req.t.Fatal("user enabled settings don't match")
	}
	if user.Name != name {
		req.t.Fatalf("user names don't match: %s and %s", user.Name, name)
	}
	return
}

func (req *IntegrationRequests) GetAndCheckUser(username, session string) {
	resp := req.GetUser(username, session)
	req.CheckUser(username, resp.Email, resp.Name, resp.Admin, resp.NeedsPasswordReset, resp.Enabled)
}

func (req *IntegrationRequests) CreateAndCheckUser(username, password, email, name string, admin, reset bool, session string) {
	resp := req.CreateUser(username, password, email, admin, reset, session)
	req.CheckUser(username, email, name, admin, reset, resp.Enabled)
}

func (req *IntegrationRequests) ModifyAndCheckUser(username, email, name string, admin, reset bool, session string) {
	resp := req.ModifyUser(username, name, email, session)
	req.CheckUser(username, email, name, admin, reset, resp.Enabled)
}

func (req *IntegrationRequests) SetEnabledStateAndCheckUser(username string, enabled bool, session string) {
	resp := req.SetEnabledState(username, enabled, session)
	req.CheckUser(username, resp.Email, resp.Name, resp.Admin, resp.NeedsPasswordReset, enabled)
}

func (req *IntegrationRequests) SetAdminStateAndCheckUser(username string, admin bool, session string) {
	resp := req.SetAdminState(username, admin, session)
	req.CheckUser(username, resp.Email, resp.Name, admin, resp.NeedsPasswordReset, resp.Enabled)
}

func (req *IntegrationRequests) EnrollHost(enrollSecret, hostIdentifier string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	body, err := json.Marshal(OsqueryEnrollPostBody{
		EnrollSecret:   enrollSecret,
		HostIdentifier: hostIdentifier,
	})

	if err != nil {
		req.t.Fatal(err.Error())
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	request, _ := http.NewRequest("POST", "/api/v1/osquery/enroll", buff)
	request.Header.Set("Content-Type", "application/json")
	req.r.ServeHTTP(response, request)

	return response
}

func (req *IntegrationRequests) OsqueryLog(nodeKey, logType string, data *json.RawMessage) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	body, err := json.Marshal(OsqueryLogPostBody{
		NodeKey: nodeKey,
		LogType: logType,
		Data:    data,
	})

	if err != nil {
		req.t.Fatal(err.Error())
	}

	buff := new(bytes.Buffer)
	buff.Write(body)
	req.t.Log(buff.String())

	request, _ := http.NewRequest("POST", "/api/v1/osquery/log", buff)
	request.Header.Set("Content-Type", "application/json")
	req.r.ServeHTTP(response, request)

	return response
}
