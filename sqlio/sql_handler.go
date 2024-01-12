package sqlio

import (
	"context"
	"database/sql"
	"fmt"
	"gglow/glow"
	"gglow/iohandler"
	"gglow/resources"
	"time"

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
	keyMap     map[string]interface{}
	folders    []string
	foldersMap map[string]interface{}
	driver     string
	serializer iohandler.Serializer
	schema     *Schema
}

func NewSqlHandler(driver, dsn string) (*SqlHandler, error) {
	sqlh := &SqlHandler{
		folder:     iohandler.Dots,
		keyList:    make([]string, 0),
		keyMap:     make(map[string]interface{}),
		folders:    make([]string, 0),
		foldersMap: make(map[string]interface{}),
		serializer: &iohandler.JsonSerializer{},
		driver:     driver,
		schema:     Schemas[0],
	}
	var err error
	sqlh.db, err = sql.Open(driver, dsn)
	if err != nil {
		iohandler.LogError(resources.MsgParseEffectPath.Format(dsn), err)
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

func (sqlh *SqlHandler) SetRootCurrent() ([]string, error) {
	return sqlh.SetCurrentFolder(iohandler.Dots)
}

func (sqlh *SqlHandler) IsRootFolder() bool {
	return sqlh.folder == iohandler.Dots
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

func (sqlh *SqlHandler) ReadEffect(title string) (*glow.Frame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q := sqlh.schema.ReadEffect(sqlh.folder, title)
	var folder, name string
	var source []byte
	row := sqlh.db.QueryRowContext(ctx, q)
	err := row.Scan(&folder, &name, &source)
	if err != nil {
		iohandler.LogError("sqlh.ReadEffect", err)
		return nil, err
	}

	frame := &glow.Frame{}
	err = sqlh.serializer.Scan(source, frame)
	if err != nil {
		iohandler.LogError("unable to decode frame", err)
		return nil, err
	}

	sqlh.folder = folder
	sqlh.title = title
	return frame, nil
}

func (sqlh *SqlHandler) IsFolder(title string) bool {
	return title == iohandler.Dots || sqlh.folder == iohandler.Dots
}

func (sqlh *SqlHandler) WriteFolder(folder string) error {
	sqlh.folder = folder
	return sqlh.WriteEffect(iohandler.Dots, nil)
}

func (sqlh *SqlHandler) ValidateNewFolderName(title string) error {
	err := iohandler.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = sqlh.isDuplicateFolder(title)
	return err
}
func (sqlh *SqlHandler) ValidateNewEffectName(title string) error {
	err := iohandler.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = sqlh.isDuplicate(title)
	return err
}

func (sqlh *SqlHandler) isDuplicateFolder(folder string) error {
	err := sqlh.findEffect(folder, iohandler.Dots)
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
	err = sqlh.WriteFolder(folder)
	if err != nil {
		return err
	}

	_, err = sqlh.SetCurrentFolder(folder)
	return err
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
			iohandler.LogError("unable to encode frame", err)
			return err
		}
		_, update = sqlh.keyMap[title]
	}

	query = sqlh.schema.WriteEffect(update, sqlh.folder, title, string(source))
	_, err = sqlh.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to execute query: '%s' %v", query, err)
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
	q := sqlh.schema.ExistsEffect(folder, title)
	row := sqlh.db.QueryRowContext(ctx, q)
	var result string
	err := row.Scan(&result)
	return err
}

func (sqlh *SqlHandler) ReadFolder(folder string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	query := sqlh.schema.SelectFolder(folder)
	rows, err := sqlh.db.QueryContext(ctx, query)
	if err != nil {
		return sqlh.keyList, fmt.Errorf("ReadFolders '%s' %v", query, err)
	}
	defer rows.Close()

	var title string
	ls := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			break
		}
		ls = append(ls, title)
	}
	return ls, err
}

func (sqlh *SqlHandler) SetCurrentFolder(folder string) ([]string, error) {
	ls, err := sqlh.ReadFolder(folder)
	if err != nil {
		return sqlh.keyList, err
	}

	sqlh.keyList = ls
	sqlh.folder = folder
	sqlh.keyMap = make(map[string]interface{})
	for _, s := range ls {
		sqlh.keyMap[s] = false
	}

	return sqlh.keyList, err
}

func (sqlh *SqlHandler) ListCurrent() []string {
	return sqlh.keyList
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
