package datastore

import (
	"sort"
	"testing"
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

func TestLabelQueries(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	testLabels(t, db)
}

func testLabels(t *testing.T, db kolide.Datastore) {
	hosts := []kolide.Host{}
	var host *kolide.Host
	var err error
	for i := 0; i < 10; i++ {
		host, err = db.EnrollHost(string(i), "foo", "", "", 10)
		assert.Nil(t, err, "enrollment should succeed")
		hosts = append(hosts, *host)
	}

	baseTime := time.Now()

	// No queries should be returned before labels or queries added
	queries, err := db.LabelQueriesForHost(host, baseTime)
	assert.Nil(t, err)
	assert.Empty(t, queries)

	// No labels should match
	labels, err := db.LabelsForHost(host.ID)
	assert.Nil(t, err)
	assert.Empty(t, labels)

	labelQueries := []kolide.Query{
		kolide.Query{
			Name:     "query1",
			Query:    "query1",
			Platform: "darwin",
		},
		kolide.Query{
			Name:     "query2",
			Query:    "query2",
			Platform: "darwin",
		},
		kolide.Query{
			Name:     "query3",
			Query:    "query3",
			Platform: "darwin",
		},
		kolide.Query{
			Name:     "query4",
			Query:    "query4",
			Platform: "darwin",
		},
	}

	for _, query := range labelQueries {
		newQuery, err := db.NewQuery(&query)
		assert.Nil(t, err)
		assert.NotZero(t, newQuery.ID)
	}

	// this one should not show up
	_, err = db.NewQuery(&kolide.Query{
		Platform: "not_darwin",
		Query:    "query5",
	})
	assert.Nil(t, err)

	// No queries should be returned before labels added
	queries, err = db.LabelQueriesForHost(host, baseTime)
	assert.Nil(t, err)
	assert.Empty(t, queries)

	newLabels := []kolide.Label{
		// Note these are intentionally out of order
		kolide.Label{
			Name:    "label3",
			QueryID: 3,
		},
		kolide.Label{
			Name:    "label1",
			QueryID: 1,
		},
		kolide.Label{
			Name:    "label2",
			QueryID: 2,
		},
		kolide.Label{
			Name:    "label4",
			QueryID: 4,
		},
	}

	for _, label := range newLabels {
		newLabel, err := db.NewLabel(&label)
		assert.Nil(t, err)
		assert.NotZero(t, newLabel.ID)
	}

	expectQueries := map[string]string{
		"1": "query3",
		"2": "query1",
		"3": "query2",
		"4": "query4",
	}

	host.Platform = "darwin"

	// Now queries should be returned
	queries, err = db.LabelQueriesForHost(host, baseTime)
	assert.Nil(t, err)
	assert.Equal(t, expectQueries, queries)

	// No labels should match with no results yet
	labels, err = db.LabelsForHost(host.ID)
	assert.Nil(t, err)
	assert.Empty(t, labels)

	// Record a query execution
	err = db.RecordLabelQueryExecutions(host, map[string]bool{"1": true}, baseTime)
	assert.Nil(t, err)

	// Use a 10 minute interval, so the query we just added should show up
	queries, err = db.LabelQueriesForHost(host, time.Now().Add(-(10 * time.Minute)))
	assert.Nil(t, err)
	delete(expectQueries, "1")
	assert.Equal(t, expectQueries, queries)

	// Record an old query execution -- Shouldn't change the return
	err = db.RecordLabelQueryExecutions(host, map[string]bool{"2": true}, baseTime.Add(-1*time.Hour))
	assert.Nil(t, err)
	queries, err = db.LabelQueriesForHost(host, time.Now().Add(-(10 * time.Minute)))
	assert.Nil(t, err)
	assert.Equal(t, expectQueries, queries)

	// Record a newer execution for that query and another
	err = db.RecordLabelQueryExecutions(host, map[string]bool{"2": false, "3": true}, baseTime)
	assert.Nil(t, err)

	// Now these should no longer show up in the necessary to run queries
	delete(expectQueries, "2")
	delete(expectQueries, "3")
	queries, err = db.LabelQueriesForHost(host, time.Now().Add(-(10 * time.Minute)))
	assert.Nil(t, err)
	assert.Equal(t, expectQueries, queries)

	// Now the two matching labels should be returned
	labels, err = db.LabelsForHost(host.ID)
	assert.Nil(t, err)
	if assert.Len(t, labels, 2) {
		labelNames := []string{labels[0].Name, labels[1].Name}
		sort.Strings(labelNames)
		assert.Equal(t, "label2", labelNames[0])
		assert.Equal(t, "label3", labelNames[1])
	}

	// A host that hasn't executed any label queries should still be asked
	// to execute those queries
	hosts[0].Platform = "darwin"
	queries, err = db.LabelQueriesForHost(host, time.Now())
	assert.Nil(t, err)
	assert.Len(t, queries, 4)

	// There should still be no labels returned for a host that never
	// executed any label queries
	labels, err = db.LabelsForHost(hosts[0].ID)
	assert.Nil(t, err)
	assert.Empty(t, labels)
}

func TestManagingLabelsOnPacks(t *testing.T) {
	ds := setup(t)
	defer teardown(t, ds)
	testManagingLabelsOnPacks(t, ds)
}

func testManagingLabelsOnPacks(t *testing.T, ds kolide.Datastore) {
	mysqlQuery := &kolide.Query{
		Name:  "MySQL",
		Query: "select pid from processes where name = 'mysqld';",
	}
	mysqlQuery, err := ds.NewQuery(mysqlQuery)
	assert.Nil(t, err)

	osqueryRunningQuery := &kolide.Query{
		Name:  "Is osquery currently running?",
		Query: "select pid from processes where name = 'osqueryd';",
	}
	osqueryRunningQuery, err = ds.NewQuery(osqueryRunningQuery)
	assert.Nil(t, err)

	monitoringPack := &kolide.Pack{
		Name: "monitoring",
	}
	err = ds.NewPack(monitoringPack)
	assert.Nil(t, err)

	mysqlLabel := &kolide.Label{
		Name:    "MySQL Monitoring",
		QueryID: mysqlQuery.ID,
	}
	mysqlLabel, err = ds.NewLabel(mysqlLabel)
	assert.Nil(t, err)

	err = ds.AddLabelToPack(mysqlLabel.ID, monitoringPack.ID)
	assert.Nil(t, err)

	labels, err := ds.GetLabelsForPack(monitoringPack)
	assert.Nil(t, err)
	assert.Len(t, labels, 1)
	assert.Equal(t, "MySQL Monitoring", labels[0].Name)

	osqueryLabel := &kolide.Label{
		Name:    "Osquery Monitoring",
		QueryID: osqueryRunningQuery.ID,
	}
	osqueryLabel, err = ds.NewLabel(osqueryLabel)
	assert.Nil(t, err)

	err = ds.AddLabelToPack(osqueryLabel.ID, monitoringPack.ID)
	assert.Nil(t, err)

	labels, err = ds.GetLabelsForPack(monitoringPack)
	assert.Nil(t, err)
	assert.Len(t, labels, 2)
}
