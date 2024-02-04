package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"gglow/glow"
	"gglow/iohandler"
	"gglow/text"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var _ iohandler.IoHandler = (*SqlHandler)(nil)

type SqlHandler struct {
	db         *sql.DB
	driver     string
	serializer iohandler.Serializer
	schema     *Schema
}

func NewSqlHandler(driver, dsn string) (*SqlHandler, error) {
	sqlh := &SqlHandler{
		serializer: &iohandler.JsonSerializer{},
		driver:     driver,
		schema:     Schemas[0],
	}
	var err error
	sqlh.db, err = sql.Open(driver, dsn)
	if err != nil {
		iohandler.LogError(text.MsgParseEffectPath.Format(dsn), err)
		return nil, err
	}
	return sqlh, nil
}

func (sqlh *SqlHandler) OnExit() error {
	return sqlh.db.Close()
}

func (sqlh *SqlHandler) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := sqlh.db.PingContext(ctx)
	if err != nil {
		iohandler.LogError("unable to connect to database", err)
	}
	return err
}

func (sqlh *SqlHandler) ReadEffect(folder, title string) (frame *glow.Frame, err error) {
	frame = &glow.Frame{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := sqlh.schema.ReadEffect(folder, title)
	var (
		scanEffect, scanFolder string
		source                 []byte
	)
	row := sqlh.db.QueryRowContext(ctx, query)
	err = row.Scan(&scanFolder, &scanEffect, &source)
	if err != nil {
		iohandler.LogError(fmt.Sprintf("sqlh.ReadEffect db Scan: %s/%s", folder, title), err)
		return
	}

	err = sqlh.serializer.Scan(source, frame)
	if err != nil {
		iohandler.LogError(fmt.Sprintf("sqlh.ReadEffect json Scan: %s/%s", folder, title), err)
		return
	}

	return
}

func (sqlh *SqlHandler) CreateFolder(folder string) error {
	return sqlh.CreateEffect(iohandler.AsFolder(), folder, nil)
}

// func (sqlh *SqlHandler) CreateNewFolder(folder string) (err error) {
// 	err = sqlh.WriteFolder(folder)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = sqlh.SetCurrentFolder(folder)
// 	return err
// }

func (sqlh *SqlHandler) CreateNewEffect(folder, title string, frame *glow.Frame) (err error) {
	err = sqlh.CreateEffect(folder, title, frame)
	return
}

func (sqlh *SqlHandler) CreateEffect(folder, title string, frame *glow.Frame) (err error) {
	return sqlh.writeEffect(false, folder, title, frame)
}

func (sqlh *SqlHandler) UpdateEffect(folder, title string, frame *glow.Frame) (err error) {
	return sqlh.writeEffect(true, folder, title, frame)
}

func (sqlh *SqlHandler) writeEffect(update bool, folder, title string, frame *glow.Frame) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var (
		query  string
		source []byte = make([]byte, 0)
	)

	if frame != nil {
		source, err = sqlh.serializer.Format(frame)
		if err != nil {
			iohandler.LogError("unable to encode frame", err)
			return err
		}
	}

	query = sqlh.schema.WriteEffect(update, folder, title, string(source))
	_, err = sqlh.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to execute query: '%s' %v", query, err)
	}

	return nil
}

func (sqlh *SqlHandler) findEffect(folder, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q := sqlh.schema.ExistsEffect(folder, title)
	row := sqlh.db.QueryRowContext(ctx, q)
	var result string
	err := row.Scan(&result)
	return err
}

func (sqlh *SqlHandler) ListKeys(folder string) (list []iohandler.KeyValue, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list = make([]iohandler.KeyValue, 0)
	var (
		query = sqlh.schema.ListFolder(folder)
		rows  *sql.Rows
	)

	rows, err = sqlh.db.QueryContext(ctx, query)
	if err != nil {
		return list, fmt.Errorf("ListFolder '%s' %v", query, err)
	}
	defer rows.Close()

	var scanFolder, scanTitle string
	for rows.Next() {
		err = rows.Scan(&scanFolder, &scanTitle)
		if err != nil {
			break
		}
		list = append(list, iohandler.KeyValue{Key: scanFolder, Value: scanTitle})
	}
	return
}

func (sqlh *SqlHandler) ListFolders() (ls []string, err error) {
	kvs, err := sqlh.ListKeys(iohandler.AsFolder())
	if err != nil {
		return nil, err
	}
	ls = make([]string, len(kvs))
	for i, v := range kvs {
		ls[i] = v.Key
	}
	return
}

func (sqlh *SqlHandler) ListEffects(folder string) (ls []string, err error) {
	ls = make([]string, 0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	query := sqlh.schema.SelectFolder(folder)
	rows, err := sqlh.db.QueryContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("ReadFolders '%s' %v", query, err)
		return
	}
	defer rows.Close()

	var scanFolder, scanTitle string
	for rows.Next() {
		err = rows.Scan(&scanFolder, &scanTitle)
		if err != nil {
			break
		}
		ls = append(ls, scanFolder)
	}
	return ls, err
}

func (sqlh *SqlHandler) Create(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(name) > 0 {
		// query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", name)
		// _, err := sqlh.db.ExecContext(ctx, query)
		// if err != nil {
		// 	iohandler.LogError("CreateNewDatabase", err)
		// 	return err
		// }
		// query = fmt.Sprintf("USE %s;", name)
		// _, err = sqlh.db.ExecContext(ctx, query)
		// if err != nil {
		// 	iohandler.LogError("USE", err)
		// 	return err
		// }
	}

	for _, query := range sqlh.schema.DropSQL {
		_, err := sqlh.db.ExecContext(ctx, query)
		if err != nil {
			iohandler.LogError("CreateNewDatabase", err)
			return err
		}
	}
	for _, query := range sqlh.schema.CreateSQL {
		_, err := sqlh.db.ExecContext(ctx, query)
		if err != nil {
			iohandler.LogError("CreateNewDatabase", err)
			return err
		}
	}
	return nil
}
