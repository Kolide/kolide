package mysql

import (
	"fmt"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) NewDistributedQueryCampaign(camp *kolide.DistributedQueryCampaign) (*kolide.DistributedQueryCampaign, error) {

	sqlStatement := `
		INSERT INTO distributed_query_campaigns (
			query_id,
			status,
			user_id
		)
		VALUES(?,?,?)
	`
	result, err := d.db.Exec(sqlStatement, camp.QueryID, camp.Status, camp.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "create distributed query campaign")
	}

	id, _ := result.LastInsertId()
	camp.ID = uint(id)
	return camp, nil
}

func (d *Datastore) DistributedQueryCampaign(id uint) (*kolide.DistributedQueryCampaign, error) {
	sql := `
		SELECT * FROM distributed_query_campaigns WHERE id = ? AND NOT deleted
	`
	campaign := &kolide.DistributedQueryCampaign{}
	if err := d.db.Get(campaign, sql, id); err != nil {
		return nil, errors.Wrap(err, "get DistributedQueryCampaign by ID")
	}

	return campaign, nil
}

func (d *Datastore) SaveDistributedQueryCampaign(camp *kolide.DistributedQueryCampaign) error {
	sqlStatement := `
		UPDATE distributed_query_campaigns SET
			query_id = ?,
			status = ?,
			user_id = ?
		WHERE id = ?
		AND NOT deleted
	`
	_, err := d.db.Exec(sqlStatement, camp.QueryID, camp.Status, camp.UserID, camp.ID)
	if err != nil {
		return errors.Wrap(err, "save distributed query campaign")
	}

	return nil
}

func (d *Datastore) DistributedQueryCampaignTargetIDs(id uint) (hostIDs []uint, labelIDs []uint, err error) {
	sqlStatement := `
		SELECT * FROM distributed_query_campaign_targets WHERE distributed_query_campaign_id = ?
	`
	targets := []kolide.DistributedQueryCampaignTarget{}

	if err = d.db.Select(&targets, sqlStatement, id); err != nil {
		return nil, nil, errors.Wrap(err, "get DistributedQueryCampaign targets by ID")
	}

	hostIDs = []uint{}
	labelIDs = []uint{}
	for _, target := range targets {
		if target.Type == kolide.TargetHost {
			hostIDs = append(hostIDs, target.TargetID)
		} else if target.Type == kolide.TargetLabel {
			labelIDs = append(labelIDs, target.TargetID)
		} else {
			return []uint{}, []uint{}, fmt.Errorf("invalid target type: %d", target.Type)
		}
	}

	return hostIDs, labelIDs, nil
}

func (d *Datastore) NewDistributedQueryCampaignTarget(target *kolide.DistributedQueryCampaignTarget) (*kolide.DistributedQueryCampaignTarget, error) {
	sqlStatement := `
		INSERT into distributed_query_campaign_targets (
			type,
			distributed_query_campaign_id,
			target_id
		)
		VALUES (?,?,?)
	`
	result, err := d.db.Exec(sqlStatement, target.Type, target.DistributedQueryCampaignID, target.TargetID)
	if err != nil {
		return nil, errors.Wrap(err, "create DistributedQueryCampaign target")
	}

	id, _ := result.LastInsertId()
	target.ID = uint(id)
	return target, nil
}

func (d *Datastore) NewDistributedQueryExecution(exec *kolide.DistributedQueryExecution) (*kolide.DistributedQueryExecution, error) {
	sqlStatement := `
		INSERT INTO distributed_query_executions (
			host_id,
			distributed_query_campaign_id,
			status,
			error,
			execution_duration
		) VALUES (?,?,?,?,?)
	`
	result, err := d.db.Exec(sqlStatement, exec.HostID, exec.DistributedQueryCampaignID,
		exec.Status, exec.Error, exec.ExecutionDuration)
	if err != nil {
		return nil, errors.Wrap(err, "create DistributedQueryCampaignExecution")
	}

	id, _ := result.LastInsertId()
	exec.ID = uint(id)

	return exec, nil
}
