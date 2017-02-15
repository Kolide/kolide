package service

import (
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide/server/config"
	"github.com/kolide/kolide/server/datastore/mysql"
	"github.com/kolide/kolide/server/kolide"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestListHosts(t *testing.T) {
	ds, err := mysql.NewTestDB(config.TestConfig().Mysql, clock.NewMockClock())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	hosts, err := svc.ListHosts(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

	_, err = ds.NewHost(&kolide.Host{
		DetailUpdateTime: time.Now(),
		SeenTime:         time.Now(),
		NodeKey:          "1",
		UUID:             "1",
		HostName:         "foo",
	})
	assert.Nil(t, err)

	hosts, err = svc.ListHosts(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 1)
}

func TestGetHost(t *testing.T) {
	ds, err := mysql.NewTestDB(config.TestConfig().Mysql, clock.NewMockClock())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	host, err := ds.NewHost(&kolide.Host{
		DetailUpdateTime: time.Now(),
		SeenTime:         time.Now(),
		NodeKey:          "1",
		UUID:             "1",
		HostName:         "foo",
	})
	assert.Nil(t, err)
	assert.NotZero(t, host.ID)

	hostVerify, err := svc.GetHost(ctx, host.ID)
	assert.Nil(t, err)

	assert.Equal(t, host.ID, hostVerify.ID)
}

func TestDeleteHost(t *testing.T) {
	ds, err := mysql.NewTestDB(config.TestConfig().Mysql, clock.NewMockClock())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	host, err := ds.NewHost(&kolide.Host{
		DetailUpdateTime: time.Now(),
		SeenTime:         time.Now(),
		NodeKey:          "1",
		UUID:             "1",
		HostName:         "foo",
	})
	assert.Nil(t, err)
	assert.NotZero(t, host.ID)

	err = svc.DeleteHost(ctx, host.ID)
	assert.Nil(t, err)

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

}
