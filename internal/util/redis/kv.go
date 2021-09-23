package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
)

func (r *Redis) Expire(key string, lifeCycleSecond uint) (int, error) {
	return redis.Int(r.Do("Expire", key, lifeCycleSecond))
}

func (r *Redis) Ping() error {
	result, e := redis.String(r.Do("PING"))

	if e != nil {
		return e
	}

	if strings.ToUpper(result) != "PONG" {
		return fmt.Errorf("ping result was not PONG: %s", result)
	}

	return nil
}
