package action

import (
	"fmt"
	"gglow/sqlio"
	"testing"
)

func TestFilters(t *testing.T) {
	ioh, err := sqlio.NewSqlHandler("sqlite3", "/home/dave/src/gglow/fyglow/glow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer ioh.OnExit()

	list, err := BuildFilterItems(ioh)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range list {
		fmt.Println(item.Folder)
		for _, effect := range item.Effects {
			fmt.Println("\t", effect)
		}
	}
}
