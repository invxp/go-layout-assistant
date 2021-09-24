package application

import (
	"fmt"
	"github.com/invxp/go-layout-assistant/internal/util/log"
	"github.com/invxp/go-layout-assistant/internal/util/mysql"
	"github.com/invxp/go-layout-assistant/internal/util/redis"
	"reflect"
	"strings"
)

//loadLogger 加载日志服务, 带自动轮转
func (application *Application) loadLogger() error {
	if !application.conf.Log.Enable {
		return nil
	}

	var err error
	application.logger, err = log.New(
		application.conf.Log.Path,
		fmt.Sprintf("%s.log", strings.ToLower(reflect.TypeOf(*application).Name())),
		application.conf.Log.MaxAgeHours,
		application.conf.Log.MaxRotationMegabytes)

	return err
}

//loadMySQL 加载MySQL, 暂未封装, 预留
func (application *Application) loadMySQL(config map[string]string) error {
	if !application.conf.MySQL.Enable {
		return nil
	}

	var err error
	application.mysql, err = mysql.New(
		application.conf.MySQL.Host,
		application.conf.MySQL.Username,
		application.conf.MySQL.Password,
		application.conf.MySQL.Database,
		application.conf.MySQL.MaxConnIdleTimeMinute,
		application.conf.MySQL.MaxConnLifeTimeMinute,
		application.conf.MySQL.MaxOpenConnections,
		application.conf.MySQL.MaxIdleConnections,
		config)

	return err
}

//loadRedis 加载Redis, 封装了一些接口, 待持续补充(Go用Redis比较多)
func (application *Application) loadRedis() error {
	if !application.conf.Redis.Enable {
		return nil
	}

	var err error
	application.redis, err = redis.New(
		application.conf.Redis.Host,
		application.conf.Redis.Password,
		application.conf.Redis.Database,
		application.conf.Redis.MaxIdle,
		application.conf.Redis.MaxActive,
		application.conf.Redis.MaxConnTimeoutSecond,
		application.conf.Redis.MaxConnIdleTimeMinute)

	return err
}
