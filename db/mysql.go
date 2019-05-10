package db

import (
	"fmt"
	"time"

	// to prevent: "Register called twice for driver mysql",change init to main.go
	// _ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBConfig struct {
	DriverName   string `yaml:"driver_name"`
	Host         string
	Port         int
	Name         string
	Username     string
	Password     string
	MaxIdleConns int  `yaml:"max_idle_conns"`
	MaxOpenConns int  `yaml:"max_open_conns"`
	ShowSQL      bool `yaml:"show_sql"`
}

func InitMySQL2Xorm(conf *DBConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(conf.DriverName,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Asia%%2FShanghai",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name))
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(conf.MaxOpenConns)
	engine.SetMaxIdleConns(conf.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Hour * 7)
	engine.ShowSQL(conf.ShowSQL)

	if err = engine.Ping(); err != nil {
		return nil, err
	}

	return engine, nil
}
