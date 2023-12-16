package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"glow-gui/effects"
	"glow-gui/glow"
	"glow-gui/resources"
	"log"
	"time"

	"fyne.io/fyne/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var _ effects.IoHandler = (*SqlHandler)(nil)

// const (
// 	dsnMYSQL      = "dave:football@tcp(192.168.40.1:3306)/test"
// 	dsnSQLLite    = "./glow.db"
// 	driverMYSQL   = "mysql"
// 	driverSQLLite = "sqlite3"
// )

type SqlHandler struct {
	folder string
	title  string
	db     *sql.DB

	keyList    []string
	keyMap     map[string]bool
	driver     string
	serializer effects.Serializer
}

// func NewMySqlHandler() *SqlHandler {
// 	return NewSqlHandler(driverMYSQL, dsnMYSQL)
// }

// func NewSqlLiteHandler() *SqlHandler {
// 	return NewSqlHandler(driverSQLLite, dsnSQLLite)
// }

func NewSqlHandler(driver, dsn string) (*SqlHandler, error) {
	sqlh := &SqlHandler{
		folder:     effects.Dots,
		keyList:    make([]string, 0),
		keyMap:     make(map[string]bool),
		serializer: &effects.JsonSerializer{},
		driver:     driver,
	}

	var err error
	sqlh.db, err = sql.Open(driver, dsn)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(dsn), err)
		return nil, err
	}
	return sqlh, nil
}

func (sqlh *SqlHandler) FolderName() string {
	return sqlh.folder
}

func (sqlh *SqlHandler) EffectName() string {
	return sqlh.title
}

func (sqlh *SqlHandler) OnExit() {
	sqlh.db.Close()
}

func (sqlh *SqlHandler) Refresh() ([]string, error) {
	return sqlh.RefreshFolder(effects.Dots)
}

func (sqlh *SqlHandler) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := sqlh.db.PingContext(ctx)
	if err != nil {
		fyne.LogError("unable to connect to database", err)
	}
	return err
}

func (sqlh *SqlHandler) ReadEffect(title string) (*glow.Frame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q := fmt.Sprintf("SELECT * FROM effects WHERE folder = '%s' AND title = '%s'",
		sqlh.folder, title)
	var folder, name string
	var source []byte
	row := sqlh.db.QueryRowContext(ctx, q)
	err := row.Scan(&folder, &name, &source)
	if err != nil {
		fyne.LogError("unable to query database", err)
		return nil, err
	}

	frame := &glow.Frame{}
	err = sqlh.serializer.Scan(source, frame)
	if err != nil {
		fyne.LogError("unable to decode frame", err)
		return nil, err
	}

	sqlh.folder = folder
	sqlh.title = title
	return frame, nil
}

func (sqlh *SqlHandler) IsFolder(title string) bool {
	return title == effects.Dots || sqlh.folder == effects.Dots
}

func (sqlh *SqlHandler) WriteFolder(folder string) error {
	sqlh.folder = folder
	return sqlh.WriteEffect(effects.Dots, nil)
}

func (sqlh *SqlHandler) ValidateNewFolderName(title string) error {
	err := effects.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = sqlh.isDuplicateFolder(title)
	return err
}
func (sqlh *SqlHandler) ValidateNewEffectName(title string) error {
	err := effects.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = sqlh.isDuplicate(title)
	return err
}

func (sqlh *SqlHandler) isDuplicateFolder(folder string) error {
	err := sqlh.findEffect(folder, effects.Dots)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}
	return fmt.Errorf(resources.MsgDuplicate.String())
}

func (sqlh *SqlHandler) CreateNewFolder(folder string) error {
	err := sqlh.isDuplicateFolder(folder)
	if err != nil {
		return err
	}
	return sqlh.WriteFolder(folder)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var (
		query  string
		source []byte = make([]byte, 0)
		err    error
		update bool
	)

	if frame != nil {
		source, err = sqlh.serializer.Format(frame)
		if err != nil {
			fyne.LogError("unable to encode frame", err)
			return err
		}
		_, update = sqlh.keyMap[title]
	}

	if update {
		query = fmt.Sprintf("UPDATE effects SET effect = '%s' WHERE folder = '%s' AND title = '%s'",
			string(source), sqlh.folder, title)
	} else {
		query = fmt.Sprintf("INSERT INTO effects (folder, title, effect) VALUES('%s', '%s', '%s')",
			sqlh.folder, title, string(source))
	}

	_, err = sqlh.db.ExecContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sqlh.keyList = append(sqlh.keyList, title)
	sqlh.keyMap[title] = false
	return nil
}

func (sqlh *SqlHandler) findEffect(folder, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q := fmt.Sprintf("SELECT title FROM effects WHERE folder = '%s' AND title = '%s';",
		folder, title)
	row := sqlh.db.QueryRowContext(ctx, q)
	var result string
	err := row.Scan(&result)
	return err
}

func (sqlh *SqlHandler) RefreshFolder(folder string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var query string
	if folder == "" || folder == effects.Dots {
		query = "SELECT folder FROM folders;"
	} else {
		query = fmt.Sprintf("SELECT title FROM effects WHERE folder = '%s' ORDER BY folder;",
			folder)
	}

	rows, err := sqlh.db.QueryContext(ctx, query)
	if err != nil {
		fyne.LogError("unable to execute search query", err)
		return sqlh.keyList, err
	}

	sqlh.folder = folder
	sqlh.keyList = make([]string, 0)
	sqlh.keyMap = make(map[string]bool)
	var title string

	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			break
		}
		sqlh.keyList = append(sqlh.keyList, title)
		sqlh.keyMap[title] = false
	}

	return sqlh.keyList, err
}

func (sqlh *SqlHandler) KeyList() []string {
	return sqlh.keyList
}

func (sqlh *SqlHandler) CreateNewDatabase() error {
	var sql_create = []string{
		"DROP VIEW IF EXISTS palettes;",
		"DROP VIEW IF EXISTS folders;",
		"DROP TABLE IF EXISTS effects;",
		"DROP TABLE IF EXISTS colors;",
		`CREATE TABLE effects (
folder VARCHAR(80) NOT NULL,
title VARCHAR(80) NOT NULL,
effect TEXT,
PRIMARY KEY (folder, title)
);`,
		"CREATE INDEX effect_title ON effects (title);",
		`CREATE VIEW folders(folder, title) AS
SELECT folder, title
FROM effects
WHERE title = '..'
ORDER BY folder;
`,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, query := range sql_create {
		_, err := sqlh.db.ExecContext(ctx, query)
		if err != nil {
			fyne.LogError("CreateNewDatabase", err)
			return err
		}
	}
	return nil
}