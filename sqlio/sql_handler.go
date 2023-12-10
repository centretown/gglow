package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"os"
	"time"

	"fyne.io/fyne/v2"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn = "dave:football@tcp(192.168.40.1:3306)/glow"
)

type SqlHandler struct {
	Folder  string
	dsn     string
	pool    *sql.DB
	keyList []string
}

func NewSqlHandler() *SqlHandler {
	var err error
	sqlh := &SqlHandler{
		keyList: make([]string, 0),
	}

	sqlh.dsn = dsn
	sqlh.pool, err = sql.Open("mysql", sqlh.dsn)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(sqlh.dsn), err)
		os.Exit(1)
	}

	return sqlh
}

func (sqlh *SqlHandler) OnExit() {
	sqlh.pool.Close()
}

func (sqlh *SqlHandler) Ping() (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = sqlh.pool.PingContext(ctx)
	if err != nil {
		fyne.LogError("unable to connect to database", err)
	}
	return err
}

func (sqlh *SqlHandler) ReadEffect(title string) (*glow.Frame, error) {
	return nil, nil
}

func (sqlh *SqlHandler) WriteEffect(title string, frame *glow.Frame) error {
	return nil
}

func (sqlh *SqlHandler) RefreshKeys(folder string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := fmt.Sprintf("SELECT p.title FROM effects AS p WHERE p.folder = '%s' ORDER BY p.title",
		folder)

	rows, err := sqlh.pool.QueryContext(ctx, q)

	if err != nil {
		fyne.LogError("unable to execute search query", err)
		return sqlh.keyList, err
	}

	sqlh.Folder = folder
	sqlh.keyList = make([]string, 0)
	var title string

	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			break
		}
		sqlh.keyList = append(sqlh.keyList, title)
	}

	return sqlh.keyList, err
}

func (sqlh *SqlHandler) KeyList() []string {
	return sqlh.keyList
}
