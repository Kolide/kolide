package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kolide/kolide/server/kolide"
)

var (
	// errBadRoute is used for mux errors
	errBadRoute = errors.New("bad route")
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	if e, ok := response.(statuser); ok {
		w.WriteHeader(e.status())
		if e.status() == http.StatusNoContent {
			return nil
		}
	}

	// if redirect, ok := response.(redirecter); ok {
	// 	if redirect.error() == nil {
	// 		w.Header().Set("Location", redirect.redirect())
	// 		// 302 redirect
	// 		w.WriteHeader(http.StatusFound)
	// 	}
	// }

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

// statuser allows response types to implement a custom
// http success status - default is 200 OK
type statuser interface {
	status() int
}

// redirector will redirect response to the given URL, we expose error here
// so we can check for errors before redirecting
type redirecter interface {
	redirect() string
	error() error
}

func idFromRequest(r *http.Request, name string) (uint, error) {
	vars := mux.Vars(r)
	id, ok := vars[name]
	if !ok {
		return 0, errBadRoute
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return uint(uid), nil
}

// default number of items to include per page
const defaultPerPage = 20

// listOptionsFromRequest parses the list options from the request parameters
func listOptionsFromRequest(r *http.Request) (kolide.ListOptions, error) {
	var err error

	pageString := r.URL.Query().Get("page")
	perPageString := r.URL.Query().Get("per_page")
	orderKey := r.URL.Query().Get("order_key")
	orderDirectionString := r.URL.Query().Get("order_direction")

	var page int = 0
	if pageString != "" {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			return kolide.ListOptions{}, errors.New("non-int page value")
		}
		if page < 0 {
			return kolide.ListOptions{}, errors.New("negative page value")
		}
	}

	// We default to 0 for per_page so that not specifying any paging
	// information gets all results
	var perPage int = 0
	if perPageString != "" {
		perPage, err = strconv.Atoi(perPageString)
		if err != nil {
			return kolide.ListOptions{}, errors.New("non-int per_page value")
		}
		if perPage <= 0 {
			return kolide.ListOptions{}, errors.New("invalid per_page value")
		}
	}

	if perPage == 0 && pageString != "" {
		// We explicitly set a non-zero default if a page is specified
		// (because the client probably intended for paging, and
		// leaving the 0 would turn that off)
		perPage = defaultPerPage
	}

	if orderKey == "" && orderDirectionString != "" {
		return kolide.ListOptions{},
			errors.New("order_key must be specified with order_direction")
	}

	var orderDirection kolide.OrderDirection
	switch orderDirectionString {
	case "desc":
		orderDirection = kolide.OrderDescending
	case "asc":
		orderDirection = kolide.OrderAscending
	case "":
		orderDirection = kolide.OrderAscending
	default:
		return kolide.ListOptions{},
			errors.New("unknown order_direction: " + orderDirectionString)

	}

	return kolide.ListOptions{
		Page:           uint(page),
		PerPage:        uint(perPage),
		OrderKey:       orderKey,
		OrderDirection: orderDirection,
	}, nil
}

func decodeNoParamsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
