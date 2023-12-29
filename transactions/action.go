package transactions

import (
	"fmt"
	"gglow/iohandler"
	"gglow/settings"
	"gglow/store"
	"strings"
)

const (
	File    = "file"
	SqlLite = "sqlite3"
	MySql   = "mysql"
)

type Action struct {
	Method  string               `yaml:"method" json:"method"`
	Input   *settings.Accessor   `yaml:"input" json:"input"`
	Outputs []*settings.Accessor `yaml:"outputs" json:"outputs"`
	Notes   []string
	Errors  []string
}

func NewAction() *Action {
	a := &Action{
		Input:   &settings.Accessor{},
		Outputs: make([]*settings.Accessor, 0),
		Notes:   make([]string, 0),
		Errors:  make([]string, 0),
	}
	return a
}

func (a *Action) AddNote(notes ...string) {
	if len(notes) == 0 {
		return
	}
	// note := fmt.Sprintf("%v", notes)
	a.Notes = append(a.Notes, notes...)
}

func (a *Action) AddError(err error) error {
	a.Errors = append(a.Errors, err.Error())
	return err
}

func (a *Action) HasErrors() bool {
	return len(a.Errors) > 0
}

func (a *Action) createDatabase(output *settings.Accessor, dataOut iohandler.OutHandler) (err error) {
	a.AddNote("create database...", output.Database)
	err = dataOut.Create(output.Database)
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("created")
	return
}

func (a *Action) connectIn(config *settings.Accessor) (handler iohandler.IoHandler, err error) {
	a.AddNote("connecting...")
	handler, err = store.NewIoHandler(config)
	if err != nil {
		a.AddError(err)
		return
	}
	a.AddNote("connected")
	return
}

func (a *Action) connectOut(config *settings.Accessor) (handler iohandler.OutHandler, err error) {
	a.AddNote("connecting...")
	handler, err = store.NewOutHandler(config)
	if err != nil {
		a.AddError(err)
		return
	}
	a.AddNote("connected")
	return
}

func (a *Action) cloneDatabase(output *settings.Accessor) (err error) {
	var dataIn iohandler.IoHandler
	dataIn, err = a.connectIn(a.Input)
	if err != nil {
		return
	}
	defer dataIn.OnExit()

	var dataOut iohandler.OutHandler
	dataOut, err = a.connectOut(output)
	if err != nil {
		return
	}
	defer dataOut.OnExit()

	err = a.createDatabase(output, dataOut)
	if err != nil {
		return
	}

	_, err = dataIn.RootFolder()
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("input read")

	a.AddNote("write database...")
	err = a.WriteDatabase(dataIn, dataOut)
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("output written")
	return nil
}

func (a *Action) Clone() (err error) {
	a.Verify()
	if a.HasErrors() {
		err = fmt.Errorf("action %s has errors", a.Method)
		return
	}

	for _, output := range a.Outputs {
		err = a.cloneDatabase(output)
		if err != nil {
			return
		}
	}
	return
}

func (a *Action) verifyConfiguration(config *settings.Accessor, refresh bool) error {
	st, err := store.NewIoHandler(config)
	if err != nil {
		return a.AddError(err)
	}
	defer st.OnExit()

	if refresh {
		_, err = st.RootFolder()
		if err != nil {
			return a.AddError(err)
		}
	}
	return nil
}

func (a *Action) Verify() {
	a.verifyConfiguration(a.Input, true)
	for _, output := range a.Outputs {
		a.verifyConfiguration(output, false)
	}
}

func (a *Action) Update() {
}

func (a *Action) Generate() (err error) {
	a.Verify()
	if a.HasErrors() {
		err = fmt.Errorf("action %s has errors", a.Method)
		return
	}

	for _, output := range a.Outputs {
		err = a.cloneDatabase(output)
		if err != nil {
			return
		}
	}
	return
}

func (a *Action) Process() (err error) {
	switch strings.ToLower(a.Method) {
	case "verify":
		a.Verify()
	case "clone":
		err = a.Clone()
	case "update":
		a.Update()
	case "generate":
		err = a.Generate()
	default:
		a.AddError(fmt.Errorf("unknown method %s", a.Method))
	}

	return
}
