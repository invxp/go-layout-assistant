package assistant

import (
	"fmt"
	"github.com/invxp/go-layout-assistant/internal/util/log"
	"github.com/invxp/go-layout-assistant/internal/util/mysql"
	"github.com/invxp/go-layout-assistant/internal/util/redis"
	"reflect"
	"strings"
)

//loadLogger 加载日志服务, 带自动轮转
func (assistant *Assistant) loadLogger() error {
	if !assistant.conf.Log.Enable {
		return nil
	}

	var err error
	assistant.logger, err = log.New(
		assistant.conf.Log.Path,
		fmt.Sprintf("%s.logf", strings.ToLower(reflect.TypeOf(*assistant).Name())),
		assistant.conf.Log.MaxAgeHours,
		assistant.conf.Log.MaxRotationMegabytes)

	return err
}

//loadMySQL 加载MySQL, 暂未封装, 预留
func (assistant *Assistant) loadMySQL(config map[string]string) error {
	if !assistant.conf.MySQL.Enable {
		return nil
	}

	var err error
	assistant.mysql, err = mysql.New(
		assistant.conf.MySQL.Host,
		assistant.conf.MySQL.Username,
		assistant.conf.MySQL.Password,
		assistant.conf.MySQL.Database,
		assistant.conf.MySQL.MaxConnIdleTimeMinute,
		assistant.conf.MySQL.MaxConnLifeTimeMinute,
		assistant.conf.MySQL.MaxOpenConnections,
		assistant.conf.MySQL.MaxIdleConnections,
		config)

	return err
}

//loadRedis 加载Redis, 封装了一些接口, 待持续补充(Go用Redis比较多)
func (assistant *Assistant) loadRedis() error {
	if !assistant.conf.Redis.Enable {
		return nil
	}

	var err error
	assistant.redis, err = redis.New(
		assistant.conf.Redis.Host,
		assistant.conf.Redis.Password,
		assistant.conf.Redis.Database,
		assistant.conf.Redis.MaxIdle,
		assistant.conf.Redis.MaxActive,
		assistant.conf.Redis.MaxConnTimeoutSecond,
		assistant.conf.Redis.MaxConnIdleTimeMinute)

	return err
}
