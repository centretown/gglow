package transactions

import (
	"fmt"
	"glow-gui/glowio"
	"glow-gui/settings"
	"glow-gui/store"
	"strings"
)

const (
	File    = "file"
	SqlLite = "sqlite3"
	MySql   = "mysql"
)

type Action struct {
	Method  string                    `yaml:"method" json:"method"`
	Input   *settings.Configuration   `yaml:"input" json:"input"`
	Outputs []*settings.Configuration `yaml:"outputs" json:"outputs"`
	Notes   []string
	Errors  []string
}

func NewAction() *Action {
	a := &Action{
		Input:   &settings.Configuration{},
		Outputs: make([]*settings.Configuration, 0),
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

func (a *Action) createDatabase(output *settings.Configuration, dataOut glowio.IoHandler) (err error) {
	a.AddNote("create database...", output.Database)
	err = dataOut.CreateNewDatabase(output.Database)
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("created")
	return
}

func (a *Action) connectDatabase(config *settings.Configuration) (handler glowio.IoHandler, err error) {
	a.AddNote("connecting...")
	handler, err = store.NewIoHandler(config)
	if err != nil {
		a.AddError(err)
		return
	}
	a.AddNote("connected")
	return
}

func (a *Action) cloneDatabase(output *settings.Configuration) (err error) {
	var dataIn, dataOut glowio.IoHandler
	dataIn, err = a.connectDatabase(a.Input)
	if err != nil {
		return
	}
	defer dataIn.OnExit()

	dataOut, err = a.connectDatabase(output)
	if err != nil {
		return
	}
	defer dataOut.OnExit()

	err = a.createDatabase(output, dataOut)
	if err != nil {
		return
	}

	_, err = dataIn.Refresh()
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

func (a *Action) verifyConfiguration(config *settings.Configuration, refresh bool) error {
	st, err := store.NewIoHandler(config)
	if err != nil {
		return a.AddError(err)
	}
	defer st.OnExit()

	if refresh {
		_, err = st.Refresh()
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

func (a *Action) Generate() {
}

func (a *Action) Process() {
	switch strings.ToLower(a.Method) {
	case "verify":
		a.Verify()
	case "clone":
		a.Clone()
	case "update":
		a.Update()
	case "generate":
		a.Generate()
	default:
		a.AddError(fmt.Errorf("unknown method %s", a.Method))
	}
}
