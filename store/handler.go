package store

import (
	"fmt"
	"glow-gui/iohandler"
	"glow-gui/settings"
	"glow-gui/sqlio"
)

func makeDSN(config *settings.Configuration) (dsn string) {
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
	return
}

func NewHandler(config *settings.Configuration) (handler iohandler.IoHandler, err error) {
	handler, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}
