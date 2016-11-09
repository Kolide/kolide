package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

// NewPack creates a new Pack
func (d *Datastore) NewPack(pack *kolide.Pack) (*kolide.Pack, error) {
	pack.MarkAsCreated(d.clock.Now())

	sql := `
		INSERT INTO packs (created_at, updated_at, name, platform )
			VALUES ( ?, ?, ?, ?)
	`

	result, err := d.db.Exec(sql, pack.CreatedAt, pack.UpdatedAt, pack.Name, pack.Platform)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	pack.ID = uint(id)
	return pack, nil
}

func (d *Datastore) SavePack(pack *kolide.Pack) error {
	panic("not implemented")
}

func (d *Datastore) DeletePack(pid uint) error {
	panic("not implemented")
}

func (d *Datastore) Pack(pid uint) (*kolide.Pack, error) {
	panic("not implemented")
}

func (d *Datastore) ListPacks(opt kolide.ListOptions) ([]*kolide.Pack, error) {
	panic("not implemented")
}

func (d *Datastore) AddQueryToPack(qid uint, pid uint) error {
	panic("not implemented")
}

func (d *Datastore) ListQueriesInPack(pack *kolide.Pack) ([]*kolide.Query, error) {
	panic("not implemented")
}

func (d *Datastore) RemoveQueryFromPack(query *kolide.Query, pack *kolide.Pack) error {
	panic("not implemented")
}

func (d *Datastore) AddLabelToPack(lid uint, pid uint) error {
	panic("not implemented")
}

func (d *Datastore) ListLabelsForPack(pack *kolide.Pack) ([]*kolide.Label, error) {
	panic("not implemented")
}

func (d *Datastore) RemoveLabelFromPack(label *kolide.Label, pack *kolide.Pack) error {
	panic("not implemented")
}
