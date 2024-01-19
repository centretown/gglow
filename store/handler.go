package store

import (
	"fmt"
	"gglow/codeio"
	"gglow/iohandler"
	"gglow/sqlio"
	"path/filepath"
	"strings"
)

func makeDSN(config *iohandler.Accessor) (dsn string) {
	switch config.Driver {
	case iohandler.DRIVER_SQLLITE3:
		if !strings.HasSuffix(config.Path, ".db") {
			config.Path = filepath.Join(config.Path, "glow.db")
		}
		dsn = config.Path
	case iohandler.DRIVER_MYSQL:
		if config.Path == "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				config.User, config.Password, config.Host, config.Port, config.Database)
		} else {
			dsn = config.Path
		}
	case iohandler.DRIVER_POSTGRES:
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

func NewIoHandler(config *iohandler.Accessor) (handler iohandler.IoHandler, err error) {
	handler, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}

func NewOutHandler(config *iohandler.Accessor) (handler iohandler.OutHandler, err error) {
	if config.Driver == iohandler.DRIVER_CODE {
		handler, err = codeio.NewCodeHandler(config.Path)
		return
	}
	handler, err = sqlio.NewSqlHandler(config.Driver, makeDSN(config))
	return
}
