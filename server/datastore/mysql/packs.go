package mysql

import (
	"database/sql"

	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

// NewPack creates a new Pack
func (d *Datastore) NewPack(pack *kolide.Pack) (*kolide.Pack, error) {

	sql := `
		INSERT INTO packs ( name, description, platform, created_by, disabled )
			VALUES ( ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(sql, pack.Name, pack.Description, pack.Platform, pack.CreatedBy, pack.Disabled)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	id, _ := result.LastInsertId()
	pack.ID = uint(id)
	return pack, nil
}

// SavePack stores changes to pack
func (d *Datastore) SavePack(pack *kolide.Pack) error {

	sql := `
			UPDATE packs
			SET name = ?, platform = ?, disabled = ?, description = ?,
			WHERE id = ? AND NOT deleted
	`

	_, err := d.db.Exec(sql, pack.Name, pack.Platform, pack.Disabled, pack.Description, pack.ID)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// DeletePack soft deletes a kolide.Pack so that it won't show up in results
func (d *Datastore) DeletePack(pid uint) error {
	sql := `
		UPDATE packs
			SET deleted_at = ?, deleted = TRUE
			WHERE id = ?
	`
	_, err := d.db.Exec(sql, d.clock.Now(), pid)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// Pack fetch kolide.Pack with matching ID
func (d *Datastore) Pack(pid uint) (*kolide.Pack, error) {
	sql := `SELECT * FROM packs WHERE id = ? AND NOT deleted`
	pack := &kolide.Pack{}
	if err := d.db.Get(pack, sql, pid); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return pack, nil
}

// ListPacks returns all kolide.Pack records limited and sorted by kolide.ListOptions
func (d *Datastore) ListPacks(opt kolide.ListOptions) ([]*kolide.Pack, error) {
	sql := `SELECT * FROM packs WHERE NOT deleted`
	sql = appendListOptionsToSQL(sql, opt)
	packs := []*kolide.Pack{}
	if err := d.db.Select(&packs, sql); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return packs, nil
}

// AddQueryToPack associates a kolide.Query with a kolide.Pack
func (d *Datastore) AddQueryToPack(qid uint, pid uint, opts kolide.QueryOptions) error {
	sql := `
	    INSERT INTO scheduled_queries (
		    pack_id,
			query_id,
			snapshot,
			differential,
			` + "`interval`" + `,
			platform,
			version,
			shard
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := d.db.Exec(sql,
		pid,
		qid,
		opts.Snapshot,
		opts.Differential,
		opts.Interval,
		opts.Platform,
		opts.Version,
		opts.Shard); err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// ListQueriesInPack gets all kolide.Query records associated with a kolide.Pack
func (d *Datastore) ListQueriesInPack(pack *kolide.Pack) ([]kolide.QueryWithOptions, error) {
	type annotatedQuery struct {
		kolide.UpdateCreateTimestamps
		kolide.DeleteFields
		ID           uint
		Name         string
		Description  string
		Query        string
		AuthorID     uint `db:"author_id"`
		Saved        bool
		Interval     uint
		Snapshot     sql.NullBool
		Differential sql.NullBool
		Platform     sql.NullString
		Version      sql.NullString
		Shard        sql.NullInt64
	}
	sql := `
	SELECT
	  q.*,
	  sq.interval,
	  sq.snapshot,
	  sq.differential,
	  sq.platform,
	  sq.version,
	  sq.shard
	FROM
	  queries q
	JOIN
	  scheduled_queries sq
	ON
	  sq.query_id = q.id
	AND
	  sq.pack_id = ?
	AND NOT q.deleted
	`
	results := []annotatedQuery{}
	if err := d.db.Select(&results, sql, pack.ID); err != nil {
		return nil, errors.DatabaseError(err)
	}

	queries := []kolide.QueryWithOptions{}
	for _, row := range results {
		query := kolide.QueryWithOptions{
			Query: kolide.Query{
				ID:          row.ID,
				Name:        row.Name,
				Description: row.Description,
				Query:       row.Query,
			},
			Options: kolide.QueryOptions{
				Interval: row.Interval,
			},
		}

		if row.Snapshot.Valid {
			query.Options.Snapshot = &row.Snapshot.Bool
		}

		if row.Differential.Valid {
			query.Options.Differential = &row.Differential.Bool
		}

		if row.Platform.Valid {
			query.Options.Platform = &row.Platform.String
		}

		if row.Version.Valid {
			query.Options.Version = &row.Version.String
		}

		if row.Shard.Valid {
			shard := uint(row.Shard.Int64)
			query.Options.Shard = &shard
		}

		queries = append(queries, query)

	}
	return queries, nil
}

// RemoveQueryFromPack disassociated a kolide.Query from a kolide.Pack
func (d *Datastore) RemoveQueryFromPack(query *kolide.Query, pack *kolide.Pack) error {
	sql := `
		DELETE FROM scheduled_queries
			WHERE pack_id = ? AND query_id = ?
	`
	if _, err := d.db.Exec(sql, pack.ID, query.ID); err != nil {
		return errors.DatabaseError(err)
	}

	return nil

}

// AddLabelToPack associates a kolide.Label with a kolide.Pack
func (d *Datastore) AddLabelToPack(lid uint, pid uint) error {
	sql := `
		INSERT INTO pack_targets ( pack_id,	type, target_id )
			VALUES ( ?, ?, ? )
	`
	_, err := d.db.Exec(sql, pid, kolide.TargetLabel, lid)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// ListLabelsForPack will return a list of kolide.Label records associated with kolide.Pack
func (d *Datastore) ListLabelsForPack(pid uint) ([]*kolide.Label, error) {
	sql := `
	SELECT
		l.id,
		l.created_at,
		l.updated_at,
		l.name
	FROM
		labels l
	JOIN
		pack_targets pt
	ON
		pt.target_id = l.id
	WHERE
		pt.type = ?
			AND
		pt.pack_id = ?
	AND NOT l.deleted
	`

	labels := []*kolide.Label{}

	if err := d.db.Select(&labels, sql, kolide.TargetLabel, pid); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return labels, nil
}

// RemoreLabelFromPack will remove the association between a kolide.Label and
// a kolide.Pack
func (d *Datastore) RemoveLabelFromPack(label *kolide.Label, pack *kolide.Pack) error {
	sql := `
		DELETE FROM pack_labels
			WHERE target_id = ? AND pack_id = ?
	`
	if _, err := d.db.Exec(sql, label.ID, pack.ID); err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

func (d *Datastore) ListHostsInPack(pid uint, opt kolide.ListOptions) ([]*kolide.Host, error) {
	sql := `
		SELECT DISTINCT h.*
		FROM hosts h
		JOIN pack_targets pt
		JOIN label_query_executions lqe
		ON (
		  pt.target_id = lqe.label_id
		  AND lqe.host_id = h.id
		  AND lqe.matches
		  AND pt.type = ?
		) OR (
		  pt.target_id = h.id
		  AND pt.type = ?
		)
		WHERE pt.pack_id = ?
	`
	sql = appendListOptionsToSQL(sql, opt)
	hosts := []*kolide.Host{}
	if err := d.db.Select(&hosts, sql, kolide.TargetLabel, kolide.TargetHost, pid); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return hosts, nil
}
