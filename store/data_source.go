package store

import (
	"fmt"
	"glow-gui/glowio"
	"glow-gui/settings"
	"glow-gui/sqlio"
)

func makeDSN(config *settings.Configuration) (dsn string) {
	fmt.Println("config", config)
	switch config.Driver {
	case "sqlite3":
		dsn = config.Path
	case "mysql":
		if config.Path == "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				config.User, config.Password, config.Host, config.Port, config.Database)
		} else {
			dsn = config.Path
		}
	case "postgres":
		if config.Path == "" {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				config.Host, config.Port, config.User, config.Password, config.Database)
		} else {
			dsn = config.Path
		}
	default:
	}
	fmt.Println("dsn", dsn)
	return
}

func NewIoHandler(config *settings.Configuration) (store glowio.IoHandler, err error) {

	if config.Driver == "sqlite" {
		config.Driver = "sqlite3"
	}

	store, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}
