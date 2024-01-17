package action

import (
	"gglow/iohandler"
	"gglow/store"
	"testing"
)

func TestCodeHandler(t *testing.T) {
	accessorIn := &iohandler.Accessor{
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

	accessorOut := &iohandler.Accessor{
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
	_, err = dataIn.SetRootCurrent()
	if err != nil {
		t.Fatal(err)
	}

	list := dataIn.ListCurrent()
	for _, item := range list {
		items, err := dataIn.SetCurrentFolder(item)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dataOut.SetCurrentFolder(item)
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
		if !dataIn.IsFolder(folder, item) {
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
