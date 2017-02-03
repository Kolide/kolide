package service

import (
	"github.com/kolide/kolide/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) GetScheduledQuery(ctx context.Context, id uint) (*kolide.ScheduledQuery, error) {
	return svc.ds.ScheduledQuery(id)
}

func (svc service) GetScheduledQueriesInPack(ctx context.Context, id uint, opts kolide.ListOptions) ([]*kolide.ScheduledQuery, error) {
	return svc.ds.ListScheduledQueriesInPack(id, opts)
}

func (svc service) ScheduleQuery(ctx context.Context, sq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
	return svc.ds.NewScheduledQuery(sq)
}

func (svc service) ModifyScheduledQuery(ctx context.Context, id uint, p *kolide.ScheduledQueryPayload) (*kolide.ScheduledQuery, error) {
	sq := &kolide.ScheduledQuery{
		ID: id,
	}

	if p.PackID != nil {
		sq.PackID = *p.PackID
	}

	if p.QueryID != nil {
		sq.QueryID = *p.QueryID
	}

	if p.Interval != nil {
		sq.Interval = *p.Interval
	}

	sq.Snapshot = p.Snapshot
	sq.Removed = p.Removed
	sq.Platform = p.Platform
	sq.Version = p.Version
	sq.Shard = p.Shard

	return svc.ds.SaveScheduledQuery(sq)
}

func (svc service) DeleteScheduledQuery(ctx context.Context, id uint) error {
	return svc.ds.DeleteScheduledQuery(id)
}
