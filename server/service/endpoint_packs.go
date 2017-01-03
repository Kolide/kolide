package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type packResponse struct {
	kolide.Pack
	QueryCount      uint   `json:"query_count"`
	TotalHostsCount uint   `json:"total_hosts_count"`
	HostIDs         []uint `json:"host_ids"`
	LabelIDs        []uint `json:"label_ids"`
}

func packResponseForPack(ctx context.Context, svc kolide.Service, pack kolide.Pack) (*packResponse, error) {
	queries, err := svc.GetScheduledQueriesInPack(ctx, pack.ID, kolide.ListOptions{})
	if err != nil {
		return nil, err
	}
	hosts, err := svc.ListHostsInPack(ctx, pack.ID, kolide.ListOptions{})
	if err != nil {
		return nil, err
	}
	hostIDs := make([]uint, len(hosts), len(hosts))
	for _, host := range hosts {
		hostIDs = append(hostIDs, host.ID)
	}
	labels, err := svc.ListLabelsForPack(ctx, pack.ID)
	labelIDs := make([]uint, len(labels), len(labels))
	for _, label := range labels {
		labelIDs = append(labelIDs, label.ID)
	}
	if err != nil {
		return nil, err
	}
	return &packResponse{
		Pack:            pack,
		QueryCount:      uint(len(queries)),
		TotalHostsCount: uint(len(hosts)),
		HostIDs:         hostIDs,
		LabelIDs:        labelIDs,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Get Pack
////////////////////////////////////////////////////////////////////////////////

type getPackRequest struct {
	ID uint
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

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return getPackResponse{Err: err}, nil
		}

		return getPackResponse{
			Pack: *resp,
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

		resp := listPacksResponse{Packs: make([]packResponse, len(packs), len(packs))}
		for _, pack := range packs {
			packResp, err := packResponseForPack(ctx, svc, *pack)
			if err != nil {
				return getPackResponse{Err: err}, nil
			}
			resp.Packs = append(resp.Packs, *packResp)
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
	Pack packResponse `json:"pack,omitempty"`
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

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return createPackResponse{Err: err}, nil
		}

		return createPackResponse{
			Pack: *resp,
		}, nil
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
	Pack packResponse `json:"pack,omitempty"`
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

		resp, err := packResponseForPack(ctx, svc, *pack)
		if err != nil {
			return modifyPackResponse{Err: err}, nil
		}

		return modifyPackResponse{
			Pack: *resp,
		}, nil
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
