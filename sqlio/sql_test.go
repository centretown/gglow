package sqlio

import (
	"fmt"
	"testing"
)

func TestMySql(t *testing.T) {
	testSql(t, driverMYSQL, dsnMYSQL)
}

func TestSqlLite(t *testing.T) {
	// testSql(t, driverSQLLite, dsnSQLLite)
	testSql(t, driverSQLLite, dsnSQLLite)
}

func testCreate(t *testing.T, driver, dsn string) {
	sqlh := NewSqlHandler(driver, dsn)
	defer sqlh.OnExit()

	err := sqlh.Ping()
	if err != nil {
		t.Fatal(err)
	}

	err = sqlh.CreateNewDatabase()
	if err != nil {
		t.Fatal(err)
	}
}

func testSql(t *testing.T, driver, dsn string) {
	sqlh := NewSqlHandler(driver, dsn)

	err := sqlh.Ping()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlh.OnExit()

	list := sqlh.KeyList()
	if err != nil {
		t.Fatal(err)
	}

	showList(t, list)

	err = sqlh.CreateNewFolder("test")
	if err != nil {
		t.Log(err)
	}

	list, err = sqlh.RefreshFolder("..")
	showList(t, list)
	if err != nil {
		t.Fatal(err)
	}
}

func showList(t *testing.T, list []string) {
	for _, item := range list {
		fmt.Println(item)
	}
}
