package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-playground/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var rdb redis.UniversalClient

func Connect() {
	if rdb != nil {
		log.Println("❌ Already connected to Redis")
		return
	}

	isClustered := os.Getenv("REDIS_CLUSTERED") == "true"
	if isClustered {
		dialTimeoutVal := time.Duration(10) * time.Second
		readTimeoutVal := time.Duration(10) * time.Second

		nodes := os.Getenv("REDIS_CLUSTER_NODES")
		if nodes == "" {
			logger.Log.Fatal("REDIS_CLUSTERED is true, but REDIS_CLUSTER_NODES is not set")
		}

		nodeAddrs := strings.Split(nodes, ",")

		clusterOptions := &redis.ClusterOptions{
			Addrs:        nodeAddrs,
			ReadOnly:     true, // Equivalent to ReadFrom.REPLICA_PREFERRED
			PoolSize:     10,
			MaxIdleConns: 10,
			MinIdleConns: 2,
			DialTimeout:  dialTimeoutVal,
			ReadTimeout:  readTimeoutVal,
			WriteTimeout: readTimeoutVal,
		}

		rdb = redis.NewClusterClient(clusterOptions)
	} else {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			db = 0 // use default DB
		}

		options := &redis.Options{
			Addr:         os.Getenv("REDIS_HOST"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           db, // use default DB
			PoolSize:     10, // 최대 활성 커넥션 수
			MaxIdleConns: 10, // 최대 idle 커넥션 수
			MinIdleConns: 2,  // 최소 idle 커넥션 수
		}

		rdb = redis.NewClient(options)
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("❌ Failed to connect to Redis: ", err)
	} else {
		log.Println("✅ Connected to redis successfully")
		if isClustered {
			options := rdb.(*redis.ClusterClient).Options()
			log.Println("⚙️ HOST =", strings.Join(options.Addrs, ", "))
		} else {
			options := rdb.(*redis.Client).Options()
			log.Println("⚙️ HOST =", options.Addr)
		}
	}
}

func Get(ctx context.Context, key string) (string, error) {
	if rdb == nil {
		return "", errors.New("redis client is not initialized")
	}
	if ctx == nil {
		return "", errors.New("context cannot be nil")
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil // Key does not exist, return empty string and no error.
		}
		return "", err // Other errors
	}

	return val, nil
}
