package db

import (
	"fmt"
	"net/url"
	"time"

	// to prevent: "Register called twice for driver mysql",change init to main.go
	// _ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/wenzhenxi/gorsa"
)

type DBConfig struct {
	Key          string `yaml:"key"` // 连接唯一名
	DriverName   string `yaml:"driver_name"`
	Host         string
	Port         int
	Name         string
	Username     string
	Password     string
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	ShowSQL      bool   `yaml:"show_sql"`
	Loc          string `yaml:"loc"`
	Secret       string `yaml:"secret"`
}

func InitMySQL2Xorm(conf *DBConfig) (*xorm.Engine, error) {
	if conf.Secret != "" {
		conf.Name, _ = gorsa.PublicDecrypt(conf.Name, conf.Secret)
		conf.Username, _ = gorsa.PublicDecrypt(conf.Username, conf.Secret)
		conf.Password, _ = gorsa.PublicDecrypt(conf.Password, conf.Secret)
	}

	if conf.Loc == "" {
		conf.Loc = url.QueryEscape("Asia/Shanghai")
	}
	
	engine, err := xorm.NewEngine(conf.DriverName,
		fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s`,
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name,
			conf.Loc))
	if err != nil {
		return nil, err
	}

	if conf.Loc == "UTC" {
		engine.DatabaseTZ = time.UTC
		engine.TZLocation = time.UTC
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

func InitPG2Xorm(conf *DBConfig) (*xorm.Engine, error) {
	if conf.Secret != "" {
		conf.Name, _ = gorsa.PublicDecrypt(conf.Name, conf.Secret)
		conf.Username, _ = gorsa.PublicDecrypt(conf.Username, conf.Secret)
		conf.Password, _ = gorsa.PublicDecrypt(conf.Password, conf.Secret)
	}

	if conf.Loc == "" {
		conf.Loc = url.QueryEscape("Asia/Shanghai")
	}

	engine, err := xorm.NewEngine(conf.DriverName,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name))
	if err != nil {
		return nil, err
	}

	if conf.Loc == "UTC" {
		engine.DatabaseTZ = time.UTC
		engine.TZLocation = time.UTC
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
