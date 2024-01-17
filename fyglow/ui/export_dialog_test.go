package ui

import (
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"gglow/sqlio"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestExport(t *testing.T) {
	app := test.NewApp()
	window := app.NewWindow("")
	accessor := &iohandler.Accessor{
		Driver: "sqlite3",
		Path:   "/home/dave/src/gglow/fyglow/glow.db",
	}
	ioh, err := sqlio.NewSqlHandler(accessor.Driver, accessor.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer ioh.OnExit()

	effect := effectio.NewEffect(ioh, app.Preferences(), accessor)
	xd := NewExportDialog(effect, window)
	// for _, item := range xd.items {
	// 	fmt.Println(item.Folder)
	// 	for _, effect := range item.Effects {
	// 		fmt.Println("\t", effect)
	// 	}
	// }

	xd.tree.OpenBranch("")
}

func TestExportB(t *testing.T) {
	app := test.NewApp()
	window := app.NewWindow("")
	accessor := &iohandler.Accessor{
		Driver: "sqlite3",
		Path:   "/home/dave/src/gglow/fyglow/glow.db",
	}
	ioh, err := sqlio.NewSqlHandler(accessor.Driver, accessor.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer ioh.OnExit()
	effect := effectio.NewEffect(ioh, app.Preferences(), accessor)
	xd := NewExportDialog(effect, window)
	xd.tree.OpenAllBranches()
	xd.Show()
}
