package redis

import (
	"github.com/gomodule/redigo/redis"
	"gitlab.com/micro/api_gateway/storage/repo"
)

type redisRepo struct {
	Rds *redis.Pool
}

func NewRedisRepo(rds *redis.Pool) repo.RedisRepo {
	return &redisRepo{
		Rds: rds,
	}
}

// SetWithTTL
func (r *redisRepo) SetWithTTL(key, value string, seconds int) (err error) {
	conn := r.Rds.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, seconds, value)
	return
}

func (r *redisRepo) Get(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("GET", key)
}
