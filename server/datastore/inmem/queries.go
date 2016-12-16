package inmem

import (
	"sort"

	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewQuery(query *kolide.Query) (*kolide.Query, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	newQuery := *query

	for _, q := range d.queries {
		if query.Name == q.Name {
			return nil, errors.ErrExists
		}
	}

	newQuery.ID = d.nextID(newQuery)
	d.queries[newQuery.ID] = &newQuery

	return &newQuery, nil
}

func (d *Datastore) SaveQuery(query *kolide.Query) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.queries[query.ID]; !ok {
		return errors.ErrNotFound
	}

	d.queries[query.ID] = query
	return nil
}

func (d *Datastore) DeleteQuery(query *kolide.Query) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.queries[query.ID]; !ok {
		return errors.ErrNotFound
	}

	delete(d.queries, query.ID)
	return nil
}

// DeleteQueries (soft) deletes the existing query objects with the provided
// IDs. The number of deleted queries is returned along with any error.
func (d *Datastore) DeleteQueries(ids []uint) (uint, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	deleted := uint(0)
	for _, id := range ids {
		if _, ok := d.queries[id]; ok {
			delete(d.queries, id)
			deleted++
		}
	}

	return deleted, nil
}

func (d *Datastore) getUserNameByID(id uint) string {
	if u, ok := d.users[id]; ok {
		return u.Name
	}
	return ""
}

func (d *Datastore) Query(id uint) (*kolide.Query, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	query, ok := d.queries[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	query.AuthorName = d.getUserNameByID(query.AuthorID)

	if err := d.loadPacksForQueries([]*kolide.Query{query}); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return query, nil
}

func (d *Datastore) ListQueries(opt kolide.ListOptions) ([]*kolide.Query, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	// We need to sort by keys to provide reliable ordering
	keys := []int{}
	for k, _ := range d.queries {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	queries := []*kolide.Query{}
	for _, k := range keys {
		q := d.queries[uint(k)]
		if q.Saved {
			q.AuthorName = d.getUserNameByID(q.AuthorID)
			queries = append(queries, q)
		}
	}

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":           "ID",
			"created_at":   "CreatedAt",
			"updated_at":   "UpdatedAt",
			"name":         "Name",
			"query":        "Query",
			"interval":     "Interval",
			"snapshot":     "Snapshot",
			"differential": "Differential",
			"platform":     "Platform",
			"version":      "Version",
		}
		if err := sortResults(queries, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := d.getLimitOffsetSliceBounds(opt, len(queries))
	queries = queries[low:high]

	if err := d.loadPacksForQueries(queries); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return queries, nil
}

// loadPacksForQueries loads the packs associated with the provided queries
func (d *Datastore) loadPacksForQueries(queries []*kolide.Query) error {
	for _, q := range queries {
		q.Packs = make([]kolide.Pack, 0)
		for _, sq := range d.scheduledQueries {
			if sq.QueryID == q.ID {
				q.Packs = append(q.Packs, *d.packs[sq.PackID])
			}
		}
	}

	return nil
}
