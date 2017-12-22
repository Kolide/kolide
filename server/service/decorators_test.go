package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

func setupDecoratorTest(r *testResource) {
	decs := []kolide.Decorator{
		kolide.Decorator{
			Type:  kolide.DecoratorLoad,
			Query: "select something from foo;",
		},
		kolide.Decorator{
			Type:  kolide.DecoratorLoad,
			Query: "select bar from foo;",
		},
		kolide.Decorator{
			Type:  kolide.DecoratorAlways,
			Query: "select x from y;",
		},
		kolide.Decorator{
			Type:     kolide.DecoratorInterval,
			Query:    "select name from system_info;",
			Interval: 3600,
		},
	}
	for _, d := range decs {
		r.ds.NewDecorator(&d)
	}
}

func testModifyDecorator(t *testing.T, r *testResource) {
	dec := &kolide.Decorator{
		Name:  "foo",
		Type:  kolide.DecoratorLoad,
		Query: "select foo from bar;",
	}
	dec, err := r.ds.NewDecorator(dec)
	require.Nil(t, err)
	buffer := bytes.NewBufferString(`{
	"payload": {
      "type": "always",
      "name": "bar",
			"query": "select baz from boom;"
		}
	}`)
	req, err := http.NewRequest("PATCH", r.server.URL+fmt.Sprintf("/api/v1/kolide/decorators/%d", dec.ID), buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var decResp decoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&decResp)
	require.Nil(t, err)
	require.NotNil(t, decResp.Decorator)
	assert.Equal(t, "select baz from boom;", decResp.Decorator.Query)
	assert.Equal(t, kolide.DecoratorAlways, decResp.Decorator.Type)
	assert.Equal(t, "bar", decResp.Decorator.Name)
}

// This test verifies that we can submit the same payload twice without
// raising an error
func testModifyDecoratorNoChanges(t *testing.T, r *testResource) {
	dec := &kolide.Decorator{
		Type:  kolide.DecoratorLoad,
		Query: "select foo from bar;",
	}
	dec, err := r.ds.NewDecorator(dec)
	require.Nil(t, err)
	buffer := bytes.NewBufferString(`{
	"payload": {
	    "type": "load",
			"query": "select foo from bar;"
		}
	}`)
	req, err := http.NewRequest("PATCH", r.server.URL+fmt.Sprintf("/api/v1/kolide/decorators/%d", dec.ID), buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)

	var decResp decoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&decResp)
	require.Nil(t, err)
	require.NotNil(t, decResp.Decorator)
	assert.Equal(t, "select foo from bar;", decResp.Decorator.Query)
	assert.Equal(t, kolide.DecoratorLoad, decResp.Decorator.Type)
}

func testListDecorator(t *testing.T, r *testResource) {
	setupDecoratorTest(r)
	req, err := http.NewRequest("GET", r.server.URL+"/api/v1/kolide/decorators", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var decs listDecoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&decs)
	require.Nil(t, err)

	assert.Len(t, decs.Decorators, 4)
}

func testNewDecorator(t *testing.T, r *testResource) {
	buffer := bytes.NewBufferString(
		`{
		 "payload": {
			"type": "load",
			"query": "select x from y;"
			}
		}`)
	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var dec decoratorResponse
	err = json.NewDecoder(resp.Body).Decode(&dec)
	require.Nil(t, err)
	require.NotNil(t, dec.Decorator)
	assert.Equal(t, kolide.DecoratorLoad, dec.Decorator.Type)
	assert.Equal(t, "select x from y;", dec.Decorator.Query)
}

// invalid json
func testNewDecoratorFailType(t *testing.T, r *testResource) {
	buffer := bytes.NewBufferString(
		`{
		 "payload": {
			"type": "zip",
			"query": "select x from y;"
			}
		}`)

	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var errStruct mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&errStruct)
	require.Nil(t, err)
	require.Len(t, errStruct.Errors, 1)
	assert.Equal(t, "invalid value, must be load, always, or interval", errStruct.Errors[0].Reason)
}

func testNewDecoratorFailValidation(t *testing.T, r *testResource) {
	buffer := bytes.NewBufferString(
		`{
			"payload": {
				"type": "interval",
				"query": "select x from y;",
				"interval": 3601
			}
		}`)

	req, err := http.NewRequest("POST", r.server.URL+"/api/v1/kolide/decorators", buffer)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var errStruct mockValidationError
	err = json.NewDecoder(resp.Body).Decode(&errStruct)
	require.Nil(t, err)
	require.Len(t, errStruct.Errors, 1)
	assert.Equal(t, "must be divisible by 60", errStruct.Errors[0].Reason)
}

func testDeleteDecorator(t *testing.T, r *testResource) {
	setupDecoratorTest(r)
	req, err := http.NewRequest("DELETE", r.server.URL+"/api/v1/kolide/decorators/1", nil)
	require.Nil(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.adminToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decs, _ := r.ds.ListDecorators()
	assert.Len(t, decs, 3)
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

var dtPtr = func(t kolide.DecoratorType) *kolide.DecoratorType { return &t }

func TestDecoratorValidation(t *testing.T) {
	ds := mock.Store{}
	ds.DecoratorFunc = func(id uint) (*kolide.Decorator, error) {
		return &kolide.Decorator{
			ID:    1,
			Query: "select x from y;",
			Type:  kolide.DecoratorAlways,
		}, nil
	}
	ds.SaveDecoratorFunc = func(dec *kolide.Decorator, opts ...kolide.OptionalArg) error {
		return nil
	}
	svc := &service{
		ds: &ds,
	}
	validator := validationMiddleware{
		Service: svc,
		ds:      &ds,
	}

	payload := kolide.DecoratorPayload{
		ID:            uint(1),
		DecoratorType: dtPtr(kolide.DecoratorInterval),
		Interval:      uintPtr(3600),
	}

	dec, err := validator.ModifyDecorator(context.Background(), payload)
	require.Nil(t, err)
	assert.Equal(t, kolide.DecoratorInterval, dec.Type)
	assert.Equal(t, uint(3600), dec.Interval)
}

func TestDecoratorValidationIntervalMissing(t *testing.T) {
	ds := mock.Store{}
	ds.DecoratorFunc = func(id uint) (*kolide.Decorator, error) {
		return &kolide.Decorator{
			ID:    1,
			Query: "select x from y;",
			Type:  kolide.DecoratorAlways,
		}, nil
	}
	ds.SaveDecoratorFunc = func(dec *kolide.Decorator, opts ...kolide.OptionalArg) error {
		return nil
	}
	svc := &service{
		ds: &ds,
	}
	validator := validationMiddleware{
		Service: svc,
		ds:      &ds,
	}

	payload := kolide.DecoratorPayload{
		ID:            uint(1),
		DecoratorType: dtPtr(kolide.DecoratorInterval),
	}

	_, err := validator.ModifyDecorator(context.Background(), payload)
	require.NotNil(t, err)
	r, ok := err.(*invalidArgumentError)
	require.True(t, ok)
	assert.Equal(t, "missing required argument", (*r)[0].reason)
}

func TestDecoratorValidationIntervalSameType(t *testing.T) {
	ds := mock.Store{}
	ds.DecoratorFunc = func(id uint) (*kolide.Decorator, error) {
		return &kolide.Decorator{
			ID:       1,
			Query:    "select x from y;",
			Type:     kolide.DecoratorInterval,
			Interval: 600,
		}, nil
	}
	ds.SaveDecoratorFunc = func(dec *kolide.Decorator, opts ...kolide.OptionalArg) error {
		return nil
	}
	svc := &service{
		ds: &ds,
	}
	validator := validationMiddleware{
		Service: svc,
		ds:      &ds,
	}

	payload := kolide.DecoratorPayload{
		ID:            uint(1),
		DecoratorType: dtPtr(kolide.DecoratorInterval),
		Interval:      uintPtr(1200),
	}

	dec, err := validator.ModifyDecorator(context.Background(), payload)
	require.Nil(t, err)
	assert.Equal(t, uint(1200), dec.Interval)
}

func TestDecoratorValidationIntervalInvalid(t *testing.T) {
	ds := mock.Store{}
	ds.DecoratorFunc = func(id uint) (*kolide.Decorator, error) {
		return &kolide.Decorator{
			ID:       1,
			Query:    "select x from y;",
			Type:     kolide.DecoratorInterval,
			Interval: 600,
		}, nil
	}
	ds.SaveDecoratorFunc = func(dec *kolide.Decorator, opts ...kolide.OptionalArg) error {
		return nil
	}
	svc := &service{
		ds: &ds,
	}
	validator := validationMiddleware{
		Service: svc,
		ds:      &ds,
	}

	payload := kolide.DecoratorPayload{
		ID:            uint(1),
		DecoratorType: dtPtr(kolide.DecoratorInterval),
		Interval:      uintPtr(1203),
	}

	_, err := validator.ModifyDecorator(context.Background(), payload)
	require.NotNil(t, err)
	r, ok := err.(*invalidArgumentError)
	require.True(t, ok)
	assert.Equal(t, "value must be divisible by 60", (*r)[0].reason)
}
