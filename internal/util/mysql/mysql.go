package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/*
工具包
MySQL数据库封装-待完善
*/

type MySQL struct {
	db *sql.DB
}

func New(host, username, password, database string, maxConnIdleTimeMinute, maxConnLifeTimeMinute, maxOpenConnections, maxIdleConnections uint, config map[string]string) (*MySQL, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?", username, password, host, database)

	for k, v := range config {
		dataSource += fmt.Sprintf("%s=%s&", k, v)
	}

	dataSource = dataSource[:len(dataSource)-1]

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	mysql := &MySQL{db: db}

	mysql.db.SetConnMaxIdleTime(time.Duration(maxConnIdleTimeMinute) * time.Minute)
	mysql.db.SetConnMaxLifetime(time.Duration(maxConnLifeTimeMinute) * time.Minute)

	mysql.db.SetMaxOpenConns(int(maxOpenConnections))
	mysql.db.SetMaxIdleConns(int(maxIdleConnections))

	err = db.Ping()

	return mysql, err
}

func (mysql *MySQL) Close() error {
	return mysql.db.Close()
}
