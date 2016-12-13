package service

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func decodeCreatePackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeModifyPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req modifyPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeDeletePackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req deletePackRequest
	req.ID = id
	return req, nil
}

func decodeGetPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getPackRequest
	req.ID = id
	return req, nil
}

func decodeListPacksRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listPacksRequest{ListOptions: opt}, nil
}

func decodeAddQueryToPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addQueryToPackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	qid, err := idFromRequest(r, "qid")
	if err != nil {
		return nil, err
	}
	pid, err := idFromRequest(r, "pid")
	if err != nil {
		return nil, err
	}

	req.QueryID = qid
	req.PackID = pid

	return req, nil
}

func decodeGetQueriesInPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	var req getQueriesInPackRequest
	req.ID = id
	return req, nil
}

func decodeDeleteQueryFromPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	qid, err := idFromRequest(r, "qid")
	if err != nil {
		return nil, err
	}
	pid, err := idFromRequest(r, "pid")
	if err != nil {
		return nil, err
	}
	var req deleteQueryFromPackRequest
	req.PackID = pid
	req.QueryID = qid
	return req, nil
}

func decodeAddLabelToPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	lid, err := idFromRequest(r, "lid")
	if err != nil {
		return nil, err
	}
	pid, err := idFromRequest(r, "pid")
	if err != nil {
		return nil, err
	}
	return addLabelToPackRequest{
		PackID:  pid,
		LabelID: lid,
	}, nil
}

func decodeGetLabelsForPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	pid, err := idFromRequest(r, "pid")
	if err != nil {
		return nil, err
	}
	var req getLabelsForPackRequest
	req.PackID = pid
	return req, nil
}

func decodeDeleteLabelFromPackRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	lid, err := idFromRequest(r, "lid")
	if err != nil {
		return nil, err
	}
	pid, err := idFromRequest(r, "pid")
	if err != nil {
		return nil, err
	}
	var req deleteLabelFromPackRequest
	req.PackID = pid
	req.LabelID = lid
	return req, nil
}
