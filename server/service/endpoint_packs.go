package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Get Pack
////////////////////////////////////////////////////////////////////////////////

type getPackRequest struct {
	ID uint
}

type packResponse struct {
	kolide.Pack
	QueryCount uint `json:"query_count"`
	HostCount  uint `json:"host_count"`
}

type getPackResponse struct {
	Pack packResponse `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r getPackResponse) error() error { return r.Err }

func makeGetPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPackRequest)

		pack, err := svc.GetPack(ctx, req.ID)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		queries, err := svc.ListQueriesInPack(ctx, pack.ID)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		hosts, err := svc.ListHostsInPack(ctx, pack.ID, kolide.ListOptions{})
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		return getPackResponse{
			Pack: packResponse{
				Pack:       *pack,
				QueryCount: uint(len(queries)),
				HostCount:  uint(len(hosts)),
			},
		}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// List Packs
////////////////////////////////////////////////////////////////////////////////

type listPacksRequest struct {
	ListOptions kolide.ListOptions
}

type listPacksResponse struct {
	Packs []packResponse `json:"packs"`
	Err   error          `json:"error,omitempty"`
}

func (r listPacksResponse) error() error { return r.Err }

func makeListPacksEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listPacksRequest)
		packs, err := svc.ListPacks(ctx, req.ListOptions)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		resp := listPacksResponse{Packs: []packResponse{}}
		for _, pack := range packs {
			queries, err := svc.ListQueriesInPack(ctx, pack.ID)
			if err != nil {
				return getPackResponse{Err: err}, nil
			}
			hosts, err := svc.ListHostsInPack(ctx, pack.ID, kolide.ListOptions{})
			if err != nil {
				return getPackResponse{Err: err}, nil
			}
			resp.Packs = append(resp.Packs, packResponse{
				Pack:       *pack,
				QueryCount: uint(len(queries)),
				HostCount:  uint(len(hosts)),
			})
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Create Pack
////////////////////////////////////////////////////////////////////////////////

type createPackRequest struct {
	payload kolide.PackPayload
}

type createPackResponse struct {
	Pack *kolide.Pack `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r createPackResponse) error() error { return r.Err }

func makeCreatePackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPackRequest)
		pack, err := svc.NewPack(ctx, req.payload)
		if err != nil {
			return createPackResponse{Err: err}, nil
		}
		return createPackResponse{pack, nil}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Pack
////////////////////////////////////////////////////////////////////////////////

type modifyPackRequest struct {
	ID      uint
	payload kolide.PackPayload
}

type modifyPackResponse struct {
	Pack *kolide.Pack `json:"pack,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r modifyPackResponse) error() error { return r.Err }

func makeModifyPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyPackRequest)
		pack, err := svc.ModifyPack(ctx, req.ID, req.payload)
		if err != nil {
			return modifyPackResponse{Err: err}, nil
		}
		return modifyPackResponse{pack, nil}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Pack
////////////////////////////////////////////////////////////////////////////////

type deletePackRequest struct {
	ID uint
}

type deletePackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deletePackResponse) error() error { return r.Err }

func makeDeletePackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePackRequest)
		err := svc.DeletePack(ctx, req.ID)
		if err != nil {
			return deletePackResponse{Err: err}, nil
		}
		return deletePackResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Add Query To Pack
////////////////////////////////////////////////////////////////////////////////

type addQueryToPackRequest struct {
	QueryID uint
	PackID  uint
	Options kolide.QueryOptions
}

type addQueryToPackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r addQueryToPackResponse) error() error { return r.Err }

func makeAddQueryToPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addQueryToPackRequest)
		err := svc.AddQueryToPack(ctx, req.QueryID, req.PackID, req.Options)
		if err != nil {
			return addQueryToPackResponse{Err: err}, nil
		}
		return addQueryToPackResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Queries In Pack
////////////////////////////////////////////////////////////////////////////////

type getQueriesInPackRequest struct {
	ID uint
}

type getQueriesInPackResponse struct {
	Queries []kolide.QueryWithOptions `json:"queries"`
	Err     error                     `json:"error,omitempty"`
}

func (r getQueriesInPackResponse) error() error { return r.Err }

func makeGetQueriesInPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getQueriesInPackRequest)
		queries, err := svc.ListQueriesInPack(ctx, req.ID)
		if err != nil {
			return getQueriesInPackResponse{Err: err}, nil
		}

		var resp getQueriesInPackResponse
		for _, query := range queries {
			resp.Queries = append(resp.Queries, query)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Query From Pack
////////////////////////////////////////////////////////////////////////////////

type deleteQueryFromPackRequest struct {
	QueryID uint
	PackID  uint
}

type deleteQueryFromPackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteQueryFromPackResponse) error() error { return r.Err }

func makeDeleteQueryFromPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteQueryFromPackRequest)
		err := svc.RemoveQueryFromPack(ctx, req.QueryID, req.PackID)
		if err != nil {
			return deleteQueryFromPackResponse{Err: err}, nil
		}
		return deleteQueryFromPackResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Add Label To Pack
////////////////////////////////////////////////////////////////////////////////

type addLabelToPackRequest struct {
	PackID  uint
	LabelID uint
}

type addLabelToPackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r addLabelToPackResponse) error() error { return r.Err }

func makeAddLabelToPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addLabelToPackRequest)
		err := svc.AddLabelToPack(ctx, req.LabelID, req.PackID)
		if err != nil {
			return addLabelToPackResponse{Err: err}, nil
		}
		return addLabelToPackResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Labels For Pack
////////////////////////////////////////////////////////////////////////////////

type getLabelsForPackRequest struct {
	PackID uint
}

type getLabelsForPackResponse struct {
	Labels []kolide.Label `json:"labels"`
	Err    error          `json:"error,omitempty"`
}

func (r getLabelsForPackResponse) error() error { return r.Err }

func makeGetLabelsForPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getLabelsForPackRequest)
		labels, err := svc.ListLabelsForPack(ctx, req.PackID)
		if err != nil {
			return getLabelsForPackResponse{Err: err}, nil
		}

		var resp getLabelsForPackResponse
		for _, label := range labels {
			resp.Labels = append(resp.Labels, *label)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Label From Pack
////////////////////////////////////////////////////////////////////////////////

type deleteLabelFromPackRequest struct {
	LabelID uint
	PackID  uint
}

type deleteLabelFromPackResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteLabelFromPackResponse) error() error { return r.Err }

func makeDeleteLabelFromPackEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteLabelFromPackRequest)
		err := svc.RemoveLabelFromPack(ctx, req.LabelID, req.PackID)
		if err != nil {
			return deleteLabelFromPackResponse{Err: err}, nil
		}
		return deleteLabelFromPackResponse{}, nil
	}
}
