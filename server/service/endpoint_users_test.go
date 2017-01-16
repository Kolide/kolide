package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testAdminUserSetAdmin(t *testing.T, r *testResource) {
	user, err := r.ds.User("user1")
	require.Nil(t, err)
	assert.False(t, user.Admin)
	inJson := `{"admin":true}`
	buff := bytes.NewBufferString(inJson)
	path := fmt.Sprintf("/api/v1/kolide/users/%d/admin", user.ID)
	req, err := http.NewRequest("PATCH", r.server.URL+path, buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	var actual adminUserResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.Nil(t, err)
	assert.Nil(t, actual.Err)
	user, err = r.ds.User("user1")
	require.Nil(t, err)
	assert.True(t, user.Admin)
}

func testNonAdminUserSetAdmin(t *testing.T, r *testResource) {
	user, err := r.ds.User("user1")
	require.Nil(t, err)
	assert.False(t, user.Admin)
	inJson := `{"admin":true}`
	buff := bytes.NewBufferString(inJson)
	path := fmt.Sprintf("/api/v1/kolide/users/%d/admin", user.ID)
	req, err := http.NewRequest("PATCH", r.server.URL+path, buff)
	require.Nil(t, err)
	// user NOT admin
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	rb := make([]byte, 500)
	resp.Body.Read(rb)
	user, err = r.ds.User("user1")
	require.Nil(t, err)
	assert.False(t, user.Admin)
}

func testAdminUserSetEnabled(t *testing.T, r *testResource) {
	user, err := r.ds.User("user1")
	require.Nil(t, err)
	assert.True(t, user.Enabled)
	inJson := `{"enabled":false}`
	buff := bytes.NewBufferString(inJson)
	path := fmt.Sprintf("/api/v1/kolide/users/%d/enable", user.ID)
	req, err := http.NewRequest("PATCH", r.server.URL+path, buff)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	var actual adminUserResponse
	err = json.NewDecoder(resp.Body).Decode(&actual)
	require.Nil(t, err)
	assert.Nil(t, actual.Err)
	user, err = r.ds.User("user1")
	require.Nil(t, err)
	assert.False(t, user.Enabled)
}

func testNonAdminUserSetEnabled(t *testing.T, r *testResource) {
	user, err := r.ds.User("user1")
	require.Nil(t, err)
	assert.True(t, user.Enabled)
	inJson := `{"enabled":false}`
	buff := bytes.NewBufferString(inJson)
	path := fmt.Sprintf("/api/v1/kolide/users/%d/enable", user.ID)
	req, err := http.NewRequest("PATCH", r.server.URL+path, buff)
	require.Nil(t, err)
	// user NOT admin
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.userToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	rb := make([]byte, 500)
	resp.Body.Read(rb)
	user, err = r.ds.User("user1")
	require.Nil(t, err)
	// shouldn't change
	assert.True(t, user.Enabled)
}
