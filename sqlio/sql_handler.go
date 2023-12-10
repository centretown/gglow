package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	dsn = "dave:football@tcp(192.168.40.1:3306)/glow"
)

type SqlHandler struct {
	Folder  string
	db      *sql.DB
	keyList []string
	keyMap  map[string]bool
}

func NewSqlHandler() *SqlHandler {
	var err error
	sqlh := &SqlHandler{
		Folder:  "effects",
		keyList: make([]string, 0),
	}

	sqlh.db, err = sql.Open("mysql", dsn)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(dsn), err)
		os.Exit(1)
	}

	return sqlh
}

func (sqlh *SqlHandler) OnExit() {
	sqlh.db.Close()
}

func (sqlh *SqlHandler) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := sqlh.db.PingContext(ctx)
	if err != nil {
		fyne.LogError("unable to connect to database", err)
	}
	return err
}

func (sqlh *SqlHandler) ReadEffect(title string) (*glow.Frame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := fmt.Sprintf("SELECT * FROM effects AS p WHERE p.folder = '%s' AND p.title = '%s'",
		sqlh.Folder, title)

	var folder, name string
	var source []byte
	row := sqlh.db.QueryRowContext(ctx, q)
	err := row.Scan(&folder, &name, &source)
	if err != nil {
		fyne.LogError("unable to query database", err)
		return nil, err
	}

	frame := &glow.Frame{}
	err = yaml.Unmarshal(source, frame)
	if err != nil {
		fyne.LogError("unable to decode frame", err)
		return nil, err
	}

	return frame, nil
}

func (sqlh *SqlHandler) ValidateNewFolderName(title string) error {
	return nil
}
func (sqlh *SqlHandler) ValidateNewEffectName(title string) error {
	return nil
}
func (sqlh *SqlHandler) CreateNewFolder(title string) error {
	return nil
}

func (sqlh *SqlHandler) isDuplicate(title string) error {
	_, found := sqlh.keyMap[title]
	if found {
		return fmt.Errorf("%s %s", title, resources.MsgDuplicate.String())
	}
	return nil
}

func (sqlh *SqlHandler) CreateNewEffect(title string, frame *glow.Frame) error {
	err := sqlh.isDuplicate(title)
	if err != nil {
		return err
	}
	return sqlh.WriteEffect(title, frame)
}

func (sqlh *SqlHandler) WriteEffect(title string, frame *glow.Frame) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		query  string
		source []byte
	)

	source, err := yaml.Marshal(frame)
	if err != nil {
		fyne.LogError("unable to encode frame", err)
		return err
	}

	_, update := sqlh.keyMap[title]
	if update {
		query = fmt.Sprintf("UPDATE effects SET effect = '%s' WHERE folder = '%s' AND title = '%s'",
			string(source), sqlh.Folder, title)
		fmt.Println("update")
	} else {
		query = fmt.Sprintf("INSERT INTO effects (folder, title, effect) VALUES('%s', '%s', '%s')",
			sqlh.Folder, title, string(source))
		fmt.Println("insert")
	}

	_, err = sqlh.db.ExecContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlh.keyList = append(sqlh.keyList, title)
	sqlh.keyMap[title] = true
	return nil
}

func (sqlh *SqlHandler) RefreshKeys(folder string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := fmt.Sprintf("SELECT p.title FROM effects AS p WHERE p.folder = '%s' ORDER BY p.folder",
		folder)

	rows, err := sqlh.db.QueryContext(ctx, q)
	if err != nil {
		fyne.LogError("unable to execute search query", err)
		return sqlh.keyList, err
	}

	sqlh.Folder = folder
	sqlh.keyList = make([]string, 0)
	sqlh.keyMap = make(map[string]bool)
	var title string

	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			break
		}
		sqlh.keyList = append(sqlh.keyList, title)
		sqlh.keyMap[title] = true
	}

	return sqlh.keyList, err
}

func (sqlh *SqlHandler) KeyList() []string {
	return sqlh.keyList
}
