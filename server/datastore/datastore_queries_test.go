package datastore

import (
	"fmt"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/patrickmn/sortutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDeleteQuery(t *testing.T, ds kolide.Datastore) {
	query := &kolide.Query{
		Name:     "foo",
		Query:    "bar",
		Interval: 123,
	}
	query, err := ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotEqual(t, query.ID, 0)

	err = ds.DeleteQuery(query)
	assert.Nil(t, err)

	assert.NotEqual(t, query.ID, 0)
	_, err = ds.Query(query.ID)
	assert.NotNil(t, err)
}

func testSaveQuery(t *testing.T, ds kolide.Datastore) {
	query := &kolide.Query{
		Name:  "foo",
		Query: "bar",
	}
	query, err := ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, query.ID)

	query.Query = "baz"
	err = ds.SaveQuery(query)

	assert.Nil(t, err)

	queryVerify, err := ds.Query(query.ID)
	assert.Nil(t, err)
	assert.Equal(t, "baz", queryVerify.Query)
}

func testListQuery(t *testing.T, ds kolide.Datastore) {
	for i := 0; i < 10; i++ {
		_, err := ds.NewQuery(&kolide.Query{
			Name:  fmt.Sprintf("name%02d", i),
			Query: fmt.Sprintf("query%02d", i),
		})
		assert.Nil(t, err)
	}

	opts := kolide.ListOptions{}
	results, err := ds.ListQueries(opts)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(results))
}

func checkPacks(t *testing.T, expected []kolide.Pack, actual []kolide.Pack) {
	sortutil.AscByField(expected, "ID")
	sortutil.AscByField(actual, "ID")
	assert.Equal(t, expected, actual)
}

func testLoadPacksForQueries(t *testing.T, ds kolide.Datastore) {
	q1 := newQuery(t, ds, "q1", "select * from time")
	q2 := newQuery(t, ds, "q2", "select * from osquery_info")

	p1 := newPack(t, ds, "p1")
	p2 := newPack(t, ds, "p2")
	p3 := newPack(t, ds, "p3")

	var err error

	addQueryToPack(t, ds, q1.ID, p2.ID)

	err = ds.LoadPacksForQueries([]*kolide.Query{q1, q2})
	require.Nil(t, err)
	checkPacks(t, []kolide.Pack{*p2}, q1.Packs)
	checkPacks(t, []kolide.Pack{}, q2.Packs)

	addQueryToPack(t, ds, q2.ID, p1.ID)
	addQueryToPack(t, ds, q2.ID, p3.ID)

	err = ds.LoadPacksForQueries([]*kolide.Query{q1, q2})
	require.Nil(t, err)
	checkPacks(t, []kolide.Pack{*p2}, q1.Packs)
	checkPacks(t, []kolide.Pack{*p1, *p3}, q2.Packs)

	addQueryToPack(t, ds, q1.ID, p3.ID)

	err = ds.LoadPacksForQueries([]*kolide.Query{q1, q2})
	require.Nil(t, err)
	checkPacks(t, []kolide.Pack{*p2, *p3}, q1.Packs)
	checkPacks(t, []kolide.Pack{*p1, *p3}, q2.Packs)
}
