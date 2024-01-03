package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"gglow/fyio"
	"gglow/glow"
	"gglow/iohandler"
	"gglow/resources"
	"log"
	"time"

	"fyne.io/fyne/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var _ iohandler.IoHandler = (*SqlHandler)(nil)

type SqlHandler struct {
	folder string
	title  string
	db     *sql.DB

	keyList    []string
	keyMap     map[string]bool
	driver     string
	serializer iohandler.Serializer
	schema     *Schema
}

func NewSqlHandler(driver, dsn string) (*SqlHandler, error) {
	sqlh := &SqlHandler{
		folder:     fyio.Dots,
		keyList:    make([]string, 0),
		keyMap:     make(map[string]bool),
		serializer: &iohandler.JsonSerializer{},
		driver:     driver,
		schema:     Schemas[0],
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

func (sqlh *SqlHandler) OnExit() error {
	return sqlh.db.Close()
}

func (sqlh *SqlHandler) RootFolder() ([]string, error) {
	return sqlh.SetFolder(fyio.Dots)
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
	q := fmt.Sprintf(sqlh.schema.ReadEffect, sqlh.folder, title)
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
	return title == fyio.Dots || sqlh.folder == fyio.Dots
}

func (sqlh *SqlHandler) WriteFolder(folder string) error {
	sqlh.folder = folder
	return sqlh.WriteEffect(fyio.Dots, nil)
}

func (sqlh *SqlHandler) ValidateNewFolderName(title string) error {
	err := fyio.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = sqlh.isDuplicateFolder(title)
	return err
}
func (sqlh *SqlHandler) ValidateNewEffectName(title string) error {
	err := fyio.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = sqlh.isDuplicate(title)
	return err
}

func (sqlh *SqlHandler) isDuplicateFolder(folder string) error {
	err := sqlh.findEffect(folder, fyio.Dots)
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

func (sqlh *SqlHandler) CreateNewEffect(title string, frame *glow.Frame) (err error) {
	err = sqlh.isDuplicate(title)
	if err != nil {
		return
	}

	err = sqlh.WriteEffect(title, frame)
	if err == nil {
		sqlh.title = title
	}
	return
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
		query = fmt.Sprintf(sqlh.schema.UpdateEffect,
			string(source), sqlh.folder, title)
	} else {
		query = fmt.Sprintf(sqlh.schema.InsertEffect,
			sqlh.folder, title, string(source))
	}

	_, err = sqlh.db.ExecContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if !update {
		sqlh.keyList = append(sqlh.keyList, title)
		sqlh.keyMap[title] = false
		//refresh the list and make this current
	}

	return nil
}

func (sqlh *SqlHandler) findEffect(folder, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q := fmt.Sprintf(sqlh.schema.FindEffect, folder, title)
	row := sqlh.db.QueryRowContext(ctx, q)
	var result string
	err := row.Scan(&result)
	return err
}

func (sqlh *SqlHandler) SetFolder(folder string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var query string
	if folder == "" || folder == fyio.Dots {
		query = sqlh.schema.Folder
	} else {
		query = fmt.Sprintf(sqlh.schema.ListEffects, folder)
	}

	rows, err := sqlh.db.QueryContext(ctx, query)
	if err != nil {
		fyne.LogError("unable to execute search query", err)
		return sqlh.keyList, err
	}
	defer rows.Close()

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

func (sqlh *SqlHandler) ListCurrentFolder() []string {
	return sqlh.keyList
}

func (sqlh *SqlHandler) Create(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(name) > 0 {
		// query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", name)
		// _, err := sqlh.db.ExecContext(ctx, query)
		// if err != nil {
		// 	fyne.LogError("CreateNewDatabase", err)
		// 	return err
		// }
		// query = fmt.Sprintf("USE %s;", name)
		// _, err = sqlh.db.ExecContext(ctx, query)
		// if err != nil {
		// 	fyne.LogError("USE", err)
		// 	return err
		// }
	}

	for _, query := range sqlh.schema.Create {
		_, err := sqlh.db.ExecContext(ctx, query)
		if err != nil {
			fyne.LogError("CreateNewDatabase", err)
			return err
		}
	}
	return nil
}
