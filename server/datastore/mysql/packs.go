package mysql

import (
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

// NewPack creates a new Pack
func (d *Datastore) NewPack(pack *kolide.Pack) (*kolide.Pack, error) {

	sql := `
		INSERT INTO packs ( name, platform )
			VALUES ( ?, ?)
	`

	result, err := d.db.Exec(sql, pack.Name, pack.Platform)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	pack.ID = uint(id)
	return pack, nil
}

// SavePack stores changes to pack
func (d *Datastore) SavePack(pack *kolide.Pack) error {

	sql := `
		UPDATE packs
			SET name = ?, platform = ?
			WHERE id = ? AND NOT deleted
	`

	_, err := d.db.Exec(sql, pack.Name, pack.Platform, pack.ID)
	return err
}

// DeletePack soft deletes a kolide.Pack so that it won't show up in results
func (d *Datastore) DeletePack(pid uint) error {
	sql := `
		UPDATE packs
			SET deleted_at = ?, deleted = TRUE
			WHERE id = ?
	`
	_, err := d.db.Exec(sql, d.clock.Now(), pid)
	return err
}

// Pack fetch kolide.Pack with matching ID
func (d *Datastore) Pack(pid uint) (*kolide.Pack, error) {
	sql := `
		SELECT * FROM packs
			WHERE id = ? AND NOT deleted
	`
	pack := &kolide.Pack{}

	if err := d.db.Get(pack, sql, pid); err != nil {
		return nil, err
	}

	return pack, nil
}

// ListPacks returns all kolide.Pack records limited and sorted by kolide.ListOptions
func (d *Datastore) ListPacks(opt kolide.ListOptions) ([]*kolide.Pack, error) {
	sql := `
		SELECT * FROM packs
			WHERE NOT deleted
	`
	sql = appendListOptionsToSQL(sql, opt)
	packs := []*kolide.Pack{}
	if err := d.db.Select(&packs, sql); err != nil {
		return nil, err
	}
	return packs, nil
}

// AddQueryToPack associates a kolide.Query with a kolide.Pack
func (d *Datastore) AddQueryToPack(qid uint, pid uint) error {
	sql := `
		INSERT INTO pack_queries ( pack_id, query_id)
			VALUES (?, ?)
	`
	_, err := d.db.Exec(sql, pid, qid)
	return err
}

// ListQueriesInPack gets all kolide.Query records associated with a kolide.Pack
func (d *Datastore) ListQueriesInPack(pack *kolide.Pack) ([]*kolide.Query, error) {
	sql := `
	SELECT
	  q.id,
	  q.created_at,
	  q.updated_at,
	  q.name,
	  q.query,
	  q.interval,
	  q.snapshot,
	  q.differential,
	  q.platform,
	  q.version
	FROM
	  queries q
	JOIN
	  pack_queries pq
	ON
	  pq.query_id = q.id
	AND
	  pq.pack_id = ?
	AND NOT q.deleted
	`
	queries := []*kolide.Query{}
	if err := d.db.Select(&queries, sql, pack.ID); err != nil {
		return nil, err
	}
	return queries, nil
}

// RemoveQueryFromPack disassociated a kolide.Query from a kolide.Pack
func (d *Datastore) RemoveQueryFromPack(query *kolide.Query, pack *kolide.Pack) error {
	sql := `
		DELETE FROM pack_queries
			WHERE pack_id = ? AND query_id = ?
	`
	_, err := d.db.Exec(sql, pack.ID, query.ID)
	return err

}

// AddLabelToPack associates a kolide.Label with a kolide.Pack
func (d *Datastore) AddLabelToPack(lid uint, pid uint) error {
	sql := `
		INSERT INTO pack_targets ( pack_id,	type,	target_id )
			VALUES ( ?, ?, ? )
	`
	_, err := d.db.Exec(sql, pid, kolide.TargetLabel, lid)

	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// ListLabelsForPack will return a list of kolide.Label records associated with kolide.Pack
func (d *Datastore) ListLabelsForPack(pack *kolide.Pack) ([]*kolide.Label, error) {
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

	if err := d.db.Select(&labels, sql, kolide.TargetLabel, pack.ID); err != nil {
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
	_, err := d.db.Exec(sql, label.ID, pack.ID)
	return err
}
