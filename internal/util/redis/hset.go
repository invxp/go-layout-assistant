package redis

import "github.com/gomodule/redigo/redis"

/*
工具包
Redis功能实现
*/

func (r *Redis) HSet(key, field, value string) (int, error) {
	return redis.Int(r.Do("HSet", key, field, value))
}

func (r *Redis) HGet(key, field string) (string, error) {
	return redis.String(r.Do("HGet", key, field))
}

func (r *Redis) HKeys(key string) ([]string, error) {
	return redis.Strings(r.Do("HKeys", key))
}

func (r *Redis) HVals(key string) ([]string, error) {
	return redis.Strings(r.Do("HVals", key))
}

func (r *Redis) HDel(key, field string) (int, error) {
	return redis.Int(r.Do("HDel", key, field))
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(r.Do("HGetAll", key))
}

func (r *Redis) HIncrBy(key, field string, value int) (int, error) {
	return redis.Int(r.Do("HIncrBy", key, field, value))
}
