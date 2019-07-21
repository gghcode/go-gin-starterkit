package db

import (
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/go-redis/redis"
)

// RedisConn can access redis
type RedisConn interface {
	Client() *redis.Client
}

type redisConn struct {
	client *redis.Client
}

func (conn *redisConn) Client() *redis.Client {
	return conn.client
}

// NewRedisConn return new connection of redis
func NewRedisConn(conf config.Configuration) RedisConn {
	conn := redisConn{
		client: redis.NewClient(&redis.Options{
			Addr: conf.Redis.Addr,
		}),
	}

	return &conn
}
