package redis

import "github.com/gomodule/redigo/redis"

func (r *Redis) LPush(key string, value string) (int, error) {
	return redis.Int(r.Do("LPush", key, value))
}

func (r *Redis) LTrim(key string, start, stop int) (string, error) {
	return redis.String(r.Do("LTrim", key, start, stop))
}

func (r *Redis) LRange(key string, start, stop int) ([]string, error) {
	return redis.Strings(r.Do("LRange", key, start, stop))
}
