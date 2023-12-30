package transactions

import (
	"gglow/iohandler"
	"gglow/settings"
	"gglow/store"
	"testing"
)

func TestCodeHandler(t *testing.T) {
	accessorIn := &settings.Accessor{
		// Driver: "sqlite3",
		// Path:   "../glow.db",
		Driver:   "postgres",
		User:     "dave",
		Password: "football",
		Host:     "localhost",
		Port:     "5432",
		Database: "test",
	}
	// input:

	accessorOut := &settings.Accessor{
		Driver: "code",
		Path:   "../generated_test_postgres",
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

		writeFolder(t, item, items, dataIn, dataOut)
	}
}

func writeFolder(t *testing.T, folder string, items []string, dataIn iohandler.IoHandler,
	dataOut iohandler.OutHandler) {
	err := dataOut.WriteFolder(folder)
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
