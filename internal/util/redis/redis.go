package redis

import (
	"fmt"
	"github.com/FZambia/sentinel"
	"github.com/gomodule/redigo/redis"
	"time"
)

/*
工具包
Redis库，支持Sentinel
*/

type Redis struct {
	p *redis.Pool
}

func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	conn := r.p.Get()

	if conn == nil {
		return nil, fmt.Errorf("redis connection was nil")
	}

	defer func() {
		_ = conn.Close()
	}()

	if conn.Err() != nil {
		return nil, fmt.Errorf("redis connection status error: %v", conn.Err())
	}

	return conn.Do(command, args...)
}

func connectToRedis(host, password string, sentinels []string, sentinelName string, database, maxIdle, maxActive, maxConnTimeoutSecond, maxIdleTimeoutMinute uint) *redis.Pool {
	if len(sentinels) > 0 {
		st := &sentinel.Sentinel{
			Addrs:      sentinels,
			MasterName: sentinelName,
			Dial: func(addr string) (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
		}
		return &redis.Pool{
			MaxIdle:     int(maxIdle),
			MaxActive:   int(maxActive),
			Wait:        true,
			IdleTimeout: time.Duration(maxIdleTimeoutMinute) * time.Minute,
			Dial: func() (redis.Conn, error) {
				masterAddr, err := st.MasterAddr()
				if err != nil {
					return nil, err
				}
				pool, err := redis.Dial("tcp", masterAddr, redis.DialPassword(password), redis.DialDatabase(int(database)), redis.DialConnectTimeout(time.Duration(maxConnTimeoutSecond)*time.Second))
				if err != nil {
					return nil, err
				}
				return pool, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
				return nil
			},
		}
	} else {
		return &redis.Pool{
			MaxIdle:     int(maxIdle),
			MaxActive:   int(maxActive),
			Wait:        true,
			IdleTimeout: time.Duration(maxIdleTimeoutMinute) * time.Minute,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", host, redis.DialPassword(password), redis.DialDatabase(int(database)), redis.DialConnectTimeout(time.Duration(maxConnTimeoutSecond)*time.Second))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
				return nil
			},
		}
	}
}

//New 连接单个Redis
func New(host, password string, database, maxIdle, maxActive, maxConnTimeoutSecond, maxIdleTimeoutMinute uint) (*Redis, error) {
	rds := &Redis{}

	rds.p = connectToRedis(host, password, nil, "", database, maxIdle, maxActive, maxConnTimeoutSecond, maxIdleTimeoutMinute)

	if err := rds.Ping(); err != nil {
		return nil, err
	}

	return rds, nil
}

//Sentinel 连接RedisSentinel
func Sentinel(sentinelName, password string, sentinels []string, database, maxIdle, maxActive, maxConnTimeoutSecond, maxIdleTimeoutMinute uint) (*Redis, error) {
	rds := &Redis{}

	rds.p = connectToRedis("", password, sentinels, sentinelName, database, maxIdle, maxActive, maxConnTimeoutSecond, maxIdleTimeoutMinute)

	if err := rds.Ping(); err != nil {
		return nil, err
	}

	return rds, nil
}

func (r *Redis) Close() error {
	return r.p.Close()
}
