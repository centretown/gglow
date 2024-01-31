package sqlio

import (
	"fmt"
	"gglow/iohandler"
	"strings"
	"testing"
)

func TestSql(t *testing.T) {
	ioh, err := NewSqlHandler("sqlite3", "/home/dave/src/gglow/fyglow/glow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer ioh.OnExit()

	err = ioh.Ping()
	if err != nil {
		t.Fatal(err)
	}

	folders, err := ioh.ListFolders()
	if err != nil {
		t.Fatal(err)
	}
	for _, folder := range folders {
		fmt.Println(folder)
		list, err := ioh.ListKeys(folder)
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			fmt.Println("\t", s)
		}

		listEffects(t, ioh, folder)
	}
}

func listEffects(t *testing.T, ioh iohandler.IoHandler, folder string) {
	list, err := ioh.ListEffects(folder)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		if iohandler.IsFolder(item) {
			fmt.Print(">")
		}
		fmt.Println(item)
	}
	fmt.Println(strings.Repeat("*", 20))
}
