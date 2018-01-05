package db

import (
	"fmt"
	"time"

	// to prevent: "Register called twice for driver mysql",change init to main.go
	// _ "github.com/go-sql-driver/mysql"#
	"github.com/go-xorm/xorm"
)

type MySQLConfig struct {
	DriverName   string `yaml:"driver_name"`
	Addr         string
	Name         string
	Username     string
	Password     string
	MaxIdleConns int  `yaml:"max_idle_conns"`
	MaxOpenConns int  `yaml:"max_open_conns"`
	ShowSQL      bool `yaml:"show_sql"`
}

func InitMySQL2Xorm(conf *MySQLConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(conf.DriverName,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FShanghai",
			conf.Username,
			conf.Password,
			conf.Addr,
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
