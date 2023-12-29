package transactions

import (
	"gglow/iohandler"
	"gglow/settings"
	"gglow/store"
	"testing"
)

func TestCodeHandler(t *testing.T) {
	accessorIn := &settings.Accessor{
		Driver: "sqlite3",
		Path:   "../glow.db",
	}

	accessorOut := &settings.Accessor{
		Driver: "code",
		Path:   "../generated_test",
	}

	dataIn, err := store.NewIoHandler(accessorIn)
	if err != nil {
		t.Fatal(err)
	}
	defer dataIn.OnExit()

	dataOut, err := store.NewOutHandler(accessorOut)
	if err != nil {
		t.Fatal(err)
	}
	defer dataOut.OnExit()

	err = dataOut.Create(accessorOut.Path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dataIn.RootFolder()
	if err != nil {
		t.Fatal(err)
	}

	list := dataIn.ListCurrentFolder()
	for _, item := range list {
		items, err := dataIn.SetFolder(item)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dataOut.SetFolder(item)
		if err != nil {
			t.Fatal(err)
		}

		writeFolder(t, items, dataIn, dataOut)
	}
}

func writeFolder(t *testing.T, items []string, dataIn iohandler.IoHandler,
	dataOut iohandler.OutHandler) {
	err := dataOut.WriteFolder(dataOut.FolderName())
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range items {
		if !dataIn.IsFolder(item) {
			frame, err := dataIn.ReadEffect(item)
			if err != nil {
				t.Fatal(err)
			}

			err = dataOut.WriteEffect(item, frame)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}
