package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

// NewQuery creates a Query
func (d *Datastore) NewQuery(query *kolide.Query) (*kolide.Query, error) {

	sql := `
		INSERT INTO queries  ( name, description, query,
			snapshot, differential, platform, version, ` + "`interval`" + `)
		VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )
	`

	result, err := d.db.Exec(sql, query.Name, query.Description, query.Query, query.Snapshot,
		query.Differential, query.Platform, query.Version, query.Interval)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}

	id, _ := result.LastInsertId()
	query.ID = uint(id)
	return query, nil
}

// SaveQuery saves changes to a Query.
func (d *Datastore) SaveQuery(q *kolide.Query) error {
	sql := `
		UPDATE queries
			SET name = ?, description = ?, query = ?, ` + "`interval`" + ` = ?, snapshot = ?,
			 	differential = ?, platform = ?, version = ?
			WHERE id = ? AND NOT deleted
	`
	_, err := d.db.Exec(sql, q.Name, q.Description, q.Query, q.Interval,
		q.Snapshot, q.Differential, q.Platform, q.Version, q.ID)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// DeleteQuery soft deletes Query identified by Query.ID
func (d *Datastore) DeleteQuery(query *kolide.Query) error {
	query.MarkDeleted(d.clock.Now())
	sql := `
		UPDATE queries
			SET deleted_at = ?, deleted = ?
			WHERE id = ?
	`
	_, err := d.db.Exec(sql, query.DeletedAt, true, query.ID)
	if err != nil {
		return errors.DatabaseError(err)
	}

	return nil
}

// Query returns a single Query identified by id, if such
// exists
func (d *Datastore) Query(id uint) (*kolide.Query, error) {
	sql := `
		SELECT * FROM queries WHERE id = ? AND NOT deleted
	`
	query := &kolide.Query{}
	if err := d.db.Get(query, sql, id); err != nil {
		return nil, errors.DatabaseError(err)
	}

	if err := d.loadPacksForQueries([]*kolide.Query{query}); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return query, nil
}

// ListQueries returns a list of queries with sort order and results limit
// determined by passed in kolide.ListOptions
func (d *Datastore) ListQueries(opt kolide.ListOptions) ([]*kolide.Query, error) {
	sql := `
		SELECT * FROM queries WHERE NOT deleted
	`
	sql = appendListOptionsToSQL(sql, opt)
	results := []*kolide.Query{}

	if err := d.db.Select(&results, sql); err != nil {
		return nil, errors.DatabaseError(err)
	}

	if err := d.loadPacksForQueries(results); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return results, nil

}

// loadPacksForQueries loads the packs associated with the provided queries
func (d *Datastore) loadPacksForQueries(queries []*kolide.Query) error {
	sql := `
		SELECT p.*, pq.query_id AS query_id
		FROM packs p
		JOIN pack_queries pq
			ON p.id = pq.pack_id
		WHERE query_id IN (?)
	`

	// Used to map the results
	id_queries := map[uint]*kolide.Query{}
	// Used for the IN clause
	ids := []uint{}
	for _, q := range queries {
		q.Packs = make([]kolide.Pack, 0)
		ids = append(ids, q.ID)
		id_queries[q.ID] = q
	}

	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return errors.DatabaseError(err)
	}

	rows := []struct {
		QueryID uint `db:"query_id"`
		kolide.Pack
	}{}

	err = d.db.Select(&rows, query, args...)
	if err != nil {
		return errors.DatabaseError(err)
	}

	for _, row := range rows {
		q := id_queries[row.QueryID]
		q.Packs = append(q.Packs, row.Pack)
	}

	return nil
}
