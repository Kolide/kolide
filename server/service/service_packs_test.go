package service

import (
	"testing"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/datastore/inmem"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestListPacks(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	queries, err := svc.ListPacks(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 0)

	_, err = ds.NewPack(&kolide.Pack{
		Name: "foo",
	})
	assert.Nil(t, err)

	queries, err = svc.ListPacks(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}

func TestGetPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	packVerify, err := svc.GetPack(ctx, pack.ID)
	assert.Nil(t, err)

	assert.Equal(t, pack.ID, packVerify.ID)
}

func TestNewPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	name := "foo"
	_, err = svc.NewPack(ctx, kolide.PackPayload{
		Name: &name,
	})

	assert.Nil(t, err)

	queries, err := ds.ListPacks(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}

func TestModifyPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	newName := "bar"
	packVerify, err := svc.ModifyPack(ctx, pack.ID, kolide.PackPayload{
		Name: &newName,
	})
	assert.Nil(t, err)

	assert.Equal(t, pack.ID, packVerify.ID)
	assert.Equal(t, "bar", packVerify.Name)
}

func TestDeletePack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	err = svc.DeletePack(ctx, pack.ID)
	assert.Nil(t, err)

	queries, err := ds.ListPacks(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 0)

}

func TestAddQueryToPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	query := &kolide.Query{
		Name:  "bar",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	queries, err := ds.ListQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 0)

	err = svc.AddQueryToPack(ctx, query.ID, pack.ID, kolide.QueryOptions{})
	assert.Nil(t, err)

	queries, err = ds.ListQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}

func TestGetQueriesInPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	query := &kolide.Query{
		Name:  "bar",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	err = ds.AddQueryToPack(query.ID, pack.ID, kolide.QueryOptions{})
	assert.Nil(t, err)

	queries, err := svc.ListQueriesInPack(ctx, pack.ID)
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}

func TestRemoveQueryFromPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	pack := &kolide.Pack{
		Name: "foo",
	}
	_, err = ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotZero(t, pack.ID)

	query := &kolide.Query{
		Name:  "bar",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	err = ds.AddQueryToPack(query.ID, pack.ID, kolide.QueryOptions{})
	assert.Nil(t, err)

	queries, err := ds.ListQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 1)

	err = svc.RemoveQueryFromPack(ctx, query.ID, pack.ID)
	assert.Nil(t, err)

	queries, err = ds.ListQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 0)
}
