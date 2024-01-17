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

	folders := []string{iohandler.DOTS, "effects", "examples"}
	// for _, folder := range folders {
	// 	readFolder(t, ioh, folder)
	// }
	for _, folder := range folders {
		listFolder(t, ioh, folder)
	}
}

func readFolder(t *testing.T, ioh iohandler.IoHandler, folder string) {
	list, err := ioh.ReadFolder(folder)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		if ioh.IsFolder(folder, item) {
			fmt.Print(">")
		}
		fmt.Println(item)
	}
	fmt.Println(strings.Repeat("*", 20))
}

func listFolder(t *testing.T, ioh iohandler.IoHandler, folder string) {

	list, err := ioh.ListFolder(folder)

	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		if ioh.IsFolder(item.Key, item.Value) {
			fmt.Print(">")
		}
		fmt.Println(item.Key, item.Value)
	}
	fmt.Println(strings.Repeat("*", 20))
}
