package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Get Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type getScheduledQueryRequest struct {
	ID uint
}

type scheduledQueryResponse struct {
	kolide.ScheduledQuery
}

type getScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r getScheduledQueryResponse) error() error { return r.Err }

func makeGetScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueryRequest)

		sq, err := svc.GetScheduledQuery(ctx, req.ID)
		if err != nil {
			return getScheduledQueryResponse{Err: err}, nil
		}

		return getScheduledQueryResponse{
			Scheduled: scheduledQueryResponse{
				ScheduledQuery: *sq,
			},
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Scheduled Queries In Pack
////////////////////////////////////////////////////////////////////////////////

type getScheduledQueriesInPackRequest struct {
	ID          uint
	ListOptions kolide.ListOptions
}

type getScheduledQueriesInPackResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r getScheduledQueriesInPackResponse) error() error { return r.Err }

func makeGetScheduledQueriesInPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getScheduledQueriesInPackRequest)
		resp := getScheduledQueriesInPackResponse{Scheduled: []scheduledQueryResponse{}}

		queries, err := svc.GetScheduledQueriesInPack(ctx, req.ID, req.ListOptions)
		if err != nil {
			return getScheduledQueriesInPackResponse{Err: err}, nil
		}

		for _, q := range queries {
			resp.Scheduled = append(resp.Scheduled, scheduledQueryResponse{
				ScheduledQuery: *q,
			})
		}

		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Schedule Queries
////////////////////////////////////////////////////////////////////////////////

type scheduleQuerySubmission struct {
	PackID       uint    `json:"pack_id"`
	QueryIDs     []uint  `json:"query_ids"`
	Interval     uint    `json:"interval"`
	Snapshot     *bool   `json:"snapshot"`
	Differential *bool   `json:"differential"`
	Platform     *string `json:"platform"`
	Version      *string `json:"version"`
	Shard        *uint   `json:"shard"`
}

type scheduleQueriesRequest struct {
	Options []scheduleQuerySubmission `json:"options"`
}

type scheduleQueriesResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r scheduleQueriesResponse) error() error { return r.Err }

func makeScheduleQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(scheduleQueriesRequest)
		resp := getScheduledQueriesInPackResponse{Scheduled: []scheduledQueryResponse{}}

		for _, submission := range req.Options {
			for _, queryID := range submission.QueryIDs {
				scheduled, err := svc.ScheduleQuery(ctx, &kolide.ScheduledQuery{
					PackID:       submission.PackID,
					QueryID:      queryID,
					Interval:     submission.Interval,
					Snapshot:     submission.Snapshot,
					Differential: submission.Differential,
					Platform:     submission.Platform,
					Version:      submission.Version,
					Shard:        submission.Shard,
				})
				if err != nil {
					return scheduleQueriesResponse{Err: err}, nil
				}
				resp.Scheduled = append(resp.Scheduled, scheduledQueryResponse{
					ScheduledQuery: *scheduled,
				})
			}
		}

		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type modifyScheduledQueryRequest struct {
	ID      uint
	payload *kolide.ScheduledQuery
}

type modifyScheduledQueryResponse struct {
	Scheduled scheduledQueryResponse `json:"scheduled,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r modifyScheduledQueryResponse) error() error { return r.Err }

func makeModifyScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyScheduledQueryRequest)

		sq := req.payload
		sq.ID = req.ID

		sq, err := svc.ModifyScheduledQuery(ctx, sq)
		if err != nil {
			return modifyScheduledQueryResponse{Err: err}, nil
		}

		return modifyScheduledQueryResponse{
			Scheduled: scheduledQueryResponse{
				ScheduledQuery: *sq,
			},
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Scheduled Query
////////////////////////////////////////////////////////////////////////////////

type deleteScheduledQueryRequest struct {
	ID uint
}

type deleteScheduledQueryResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteScheduledQueryResponse) error() error { return r.Err }

func makeDeleteScheduledQueryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteScheduledQueryRequest)

		err := svc.DeleteScheduledQuery(ctx, req.ID)
		if err != nil {
			return deleteScheduledQueryResponse{Err: err}, nil
		}

		return deleteScheduledQueryResponse{}, nil
	}
}
