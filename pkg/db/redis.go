package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type RedisOptions struct {
	Host          string
	Port          int
	Address       []string
	Username      string
	Password      string
	Database      int
	MasterName    string
	MinIdleConns  int
	MaxActive     int
	EnableCluster bool
	Timeout       int
}

var (
	// singlePool  pool struct
	singlePool atomic.Value
	// redisUp  redis is up
	redisUp atomic.Value
)

// singleton get redis instance
func singleton() redis.UniversalClient {
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

// connectSingleton
// determine whether to connect to the redis server
func connectSingleton(opts *RedisOptions) (connect bool) {
	if singleton() == nil {

		singlePool.Store(NewRedisClusterPool(opts))
		connect = true
	}

	return true
}

// NewRedisClusterPool create a redis cluster pool.
func NewRedisClusterPool(opts *RedisOptions) redis.UniversalClient {

	poolSize := 500
	if opts.MaxActive > 0 {
		poolSize = opts.MaxActive
	}

	timeout := 5 * time.Second
	if opts.Timeout > 0 {
		timeout = time.Duration(opts.Timeout) * time.Second
	}

	var client redis.UniversalClient
	redisOptions := &redis.UniversalOptions{
		Addrs:        getRedisAddress(opts),
		MasterName:   opts.MasterName,
		Password:     opts.Password,
		DB:           opts.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		MinIdleConns: opts.MinIdleConns,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
	}

	if redisOptions.MasterName != "" {
		client = redis.NewFailoverClient(redisOptions.Failover())
	} else if opts.EnableCluster {
		client = redis.NewClusterClient(redisOptions.Cluster())
	} else {
		client = redis.NewClient(redisOptions.Simple())
	}

	return client
}

// getRedisAddress
// get address of redis
func getRedisAddress(opts *RedisOptions) (address []string) {

	if len(opts.Address) > 0 {
		address = opts.Address
	}

	if len(address) == 0 && opts.Port != 0 {
		address = append(address, net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port)))
	}

	return address
}

type RedisCluster struct {
	KeyPrefix string
}

func clusterConnectionIsOpen(r RedisCluster) bool {
	c := singleton()
	testKey := r.KeyPrefix + "redis-test" + uuid.Must(uuid.NewV4(), nil).String()
	if err := c.Set(context.Background(), testKey, "test-value", time.Second).Err(); err != nil {
		log.Printf("Error trying to set test key: %s", err)
		return false
	}

	if err := c.Get(context.Background(), testKey).Err(); err != nil {
		log.Printf("Error trying to get test key: %s", err)
		return false
	}

	return true
}

// Connect2Redis
// starts a go routine that periodically tries to connect to redis.
func Connect2Redis(ctx context.Context, opts *RedisOptions) {

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	rc := []RedisCluster{{}}

	var ok bool
	for _, c := range rc {

		if !connectSingleton(opts) {
			break
		}

		if !clusterConnectionIsOpen(c) {
			break
		}

		ok = true
	}
	redisUp.Store(ok)

again:
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			for _, c := range rc {

				if !connectSingleton(opts) {
					redisUp.Store(false)
					goto again
				}

				if !clusterConnectionIsOpen(c) {
					redisUp.Store(false)
					goto again
				}
			}
			redisUp.Store(true)
		}
	}
}
