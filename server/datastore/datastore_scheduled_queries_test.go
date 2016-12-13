package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/kolide/kolide-ose/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testScheduledQuery(t *testing.T, ds kolide.Datastore) {
	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, true)

	query, err := ds.ScheduledQuery(sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)
}

func testDeleteScheduledQuery(t *testing.T, ds kolide.Datastore) {
	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, true)

	query, err := ds.ScheduledQuery(sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)

	err = ds.DeleteScheduledQuery(sq1.ID)
	require.Nil(t, err)

	_, err = ds.ScheduledQuery(sq1.ID)
	require.NotNil(t, err)
}

func testListScheduledQueriesInPack(t *testing.T, ds kolide.Datastore) {
	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID)
	q2 := test.NewQuery(t, ds, "bar", "select * from time;", u1.ID)
	p1 := test.NewPack(t, ds, "baz")

	test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, true)

	queries, err := ds.ListScheduledQueriesInPack(p1.ID, kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, queries, 1)
	assert.Equal(t, uint(60), queries[0].Interval)

	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, false, true)
	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, true, false)

	queries, err = ds.ListScheduledQueriesInPack(p1.ID, kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, queries, 3)
}

func testSaveScheduledQuery(t *testing.T, ds kolide.Datastore) {
	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, true)

	query, err := ds.ScheduledQuery(sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)

	query.Interval = uint(120)
	query, err = ds.SaveScheduledQuery(query)
	require.Nil(t, err)
	assert.Equal(t, uint(120), query.Interval)

	queryVerify, err := ds.ScheduledQuery(sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(120), queryVerify.Interval)
}
