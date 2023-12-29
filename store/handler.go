package store

import (
	"fmt"
	"gglow/codeio"
	"gglow/iohandler"
	"gglow/settings"
	"gglow/sqlio"
)

func makeDSN(config *settings.Accessor) (dsn string) {
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

func NewIoHandler(config *settings.Accessor) (handler iohandler.IoHandler, err error) {
	handler, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}

func NewOutHandler(config *settings.Accessor) (handler iohandler.OutHandler, err error) {
	if config.Driver == "code" {
		handler, err = codeio.NewCodeHandler(config.Path)
		return
	}
	handler, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}
