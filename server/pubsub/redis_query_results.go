package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kolide/fleet/server/kolide"
	"github.com/mna/redisc"
	"github.com/pkg/errors"
)

type redisQueryResults struct {
	// connection cluster
	cluster *redisc.Cluster
}

var _ kolide.QueryResultStore = &redisQueryResults{}

// NewRedisCluster creates a RedisCluster connection pool
func NewRedisCluster(nodes []string, password string) *redisc.Cluster {
	return &redisc.Cluster{
		StartupNodes: nodes,
		CreatePool: func(server string, options ...redis.DialOption) (*redis.Pool, error) {
			return NewRedisPool(server, password, options...), nil
		},
	}
}

// NewRedisPool creates a Redis connection pool using the provided server
// address and password.
func NewRedisPool(server, password string, options ...redis.DialOption) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, options...)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// NewRedisQueryResults creats a new Redis implementation of the
// QueryResultStore interface using the provided Redis connection pool.
func NewRedisQueryResults(cluster *redisc.Cluster) *redisQueryResults {
	return &redisQueryResults{cluster: cluster}
}

func pubSubForID(id uint) string {
	return fmt.Sprintf("results_%d", id)
}

func (r *redisQueryResults) WriteResult(result kolide.DistributedQueryResult) error {
	conn := r.cluster.Get()
	defer conn.Close()

	channelName := pubSubForID(result.DistributedQueryCampaignID)

	jsonVal, err := json.Marshal(&result)
	if err != nil {
		return errors.Wrap(err, "marshalling JSON for result")
	}

	n, err := redis.Int(conn.Do("PUBLISH", channelName, string(jsonVal)))
	if err != nil {
		return errors.Wrap(err, "PUBLISH failed to channel "+channelName)
	}
	if n == 0 {
		return noSubscriberError{channelName}
	}

	return nil
}

// receiveMessages runs in a goroutine, forwarding messages from the Pub/Sub
// connection over the provided channel. This effectively allows a select
// statement to run on conn.Receive() (by running on the channel that is being
// fed by this function)
func receiveMessages(conn *redis.PubSubConn, outChan chan<- interface{}) {
	defer func() {
		close(outChan)
	}()

	for {
		msg := conn.Receive()
		outChan <- msg
		switch msg := msg.(type) {
		case error:
			// If an error occurred (i.e. connection was closed),
			// then we should exit
			return
		case redis.Subscription:
			// If the subscription count is 0, the ReadChannel call
			// that invoked this goroutine has unsubscribed, and we
			// can exit
			if msg.Count == 0 {
				return
			}
		}
	}
}

func (r *redisQueryResults) ReadChannel(ctx context.Context, query kolide.DistributedQueryCampaign) (<-chan interface{}, error) {
	outChannel := make(chan interface{})

	conn := redis.PubSubConn{Conn: r.cluster.Get()}

	pubSubName := pubSubForID(query.ID)
	conn.Subscribe(pubSubName)

	msgChannel := make(chan interface{})
	// Run a separate goroutine feeding redis messages into
	// msgChannel
	go receiveMessages(&conn, msgChannel)

	go func() {
		defer close(outChannel)
		defer conn.Close()

		for {
			// Loop reading messages from conn.Receive() (via
			// msgChannel) until the context is cancelled.
			select {
			case msg, ok := <-msgChannel:
				if !ok {
					return
				}
				switch msg := msg.(type) {
				case redis.Message:
					var res kolide.DistributedQueryResult
					err := json.Unmarshal(msg.Data, &res)
					if err != nil {
						outChannel <- err
					}
					outChannel <- res
				case error:
					outChannel <- errors.Wrap(msg, "reading from redis")
				}

			case <-ctx.Done():
				conn.Unsubscribe()

			}
		}

	}()
	return outChannel, nil
}

// HealthCheck verifies that the redis backend can be pinged, returning an error
// otherwise.
func (r *redisQueryResults) HealthCheck() error {
	conn := r.cluster.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err
}
