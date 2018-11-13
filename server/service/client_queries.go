package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

// ApplyQueries sends the list of Queries to be applied (upserted) to the
// Fleet instance.
func (c *Client) ApplyQueries(specs []*kolide.QuerySpec) error {
	req := applyQuerySpecsRequest{Specs: specs}
	response, err := c.AuthenticatedDo("POST", "/api/v1/kolide/spec/queries", req)
	if err != nil {
		return errors.Wrap(err, "POST /api/v1/kolide/spec/queries")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf(
			"apply queries received status %d %s",
			response.StatusCode,
			extractServerErrorText(response.Body),
		)
	}

	var responseBody applyQuerySpecsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return errors.Wrap(err, "decode apply query spec response")
	}

	if responseBody.Err != nil {
		return errors.Errorf("apply query spec: %s", responseBody.Err)
	}

	return nil
}

// GetQuery retrieves the list of all Queries.
func (c *Client) GetQuery(name string) (*kolide.QuerySpec, error) {
	u := &url.URL{Path: name}
	verb, encodedPath := "GET", path.Join("/api/v1/kolide/spec/queries/", u.String())
	response, err := c.AuthenticatedDo(verb, encodedPath, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %s", verb, encodedPath)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, notFoundErr{}
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"get query received status %d %s",
			response.StatusCode,
			extractServerErrorText(response.Body),
		)
	}

	var responseBody getQuerySpecResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "decode get query spec response")
	}

	if responseBody.Err != nil {
		return nil, errors.Errorf("get query spec: %s", responseBody.Err)
	}

	return responseBody.Spec, nil
}

// GetQueries retrieves the list of all Queries.
func (c *Client) GetQueries() ([]*kolide.QuerySpec, error) {
	response, err := c.AuthenticatedDo("GET", "/api/v1/kolide/spec/queries", nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET /api/v1/kolide/spec/queries")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"get queries received status %d %s",
			response.StatusCode,
			extractServerErrorText(response.Body),
		)
	}

	var responseBody getQuerySpecsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "decode get query spec response")
	}

	if responseBody.Err != nil {
		return nil, errors.Errorf("get query spec: %s", responseBody.Err)
	}

	return responseBody.Specs, nil
}

// DeleteQuery deletes the query with the matching name.
func (c *Client) DeleteQuery(name string) error {
	u := &url.URL{Path: name}
	verb, encodedPath := "DELETE", path.Join("/api/v1/kolide/queries/", u.String())
	response, err := c.AuthenticatedDo(verb, encodedPath, nil)
	if err != nil {
		return errors.Wrapf(err, "%s %s", verb, encodedPath)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		return notFoundErr{}
	}
	if response.StatusCode != http.StatusOK {
		return errors.Errorf(
			"delete query received status %d %s",
			response.StatusCode,
			extractServerErrorText(response.Body),
		)
	}

	var responseBody deleteQueryResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return errors.Wrap(err, "decode get query spec response")
	}

	if responseBody.Err != nil {
		return errors.Errorf("get query spec: %s", responseBody.Err)
	}

	return nil
}
