package service

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeNewDecoratorRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var dec newDecoratorRequest
	err := json.NewDecoder(req.Body).Decode(&dec)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func decodeDeleteDecoratorRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	id, err := idFromRequest(req, "id")
	if err != nil {
		return nil, err
	}
	return deleteDecoratorRequest{ID: id}, nil
}
