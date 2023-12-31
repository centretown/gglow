package transactions

import (
	"fmt"
	"gglow/iohandler"
	"gglow/store"
	"strings"
)

type Action struct {
	Method    string
	Input     *iohandler.Accessor
	Filters   []Filter
	Outputs   []*iohandler.Accessor
	Notes     []string
	Errors    []string
	filterMap map[string]map[string]bool
}

func NewAction() *Action {
	a := &Action{
		Input:   &iohandler.Accessor{},
		Filters: make([]Filter, 0),
		Outputs: make([]*iohandler.Accessor, 0),
		Notes:   make([]string, 0),
		Errors:  make([]string, 0),
	}
	return a
}

func (a *Action) AddNote(notes ...string) {
	var note string
	last := len(notes) - 1
	for i, s := range notes {
		if i < last {
			s += " "
		}
		note += s
	}
	a.Notes = append(a.Notes, note)
}

func (a *Action) AddError(err error) error {
	a.Errors = append(a.Errors, err.Error())
	return err
}

func (a *Action) HasErrors() bool {
	return len(a.Errors) > 0
}

func (a *Action) createDatabase(output *iohandler.Accessor, dataOut iohandler.OutHandler) (err error) {
	a.AddNote("create database", output.Database)
	err = dataOut.Create(output.Database)
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("created", output.Database)
	return
}

func (a *Action) connectIn(config *iohandler.Accessor) (handler iohandler.IoHandler, err error) {
	a.AddNote("connecting to input", config.Database)
	handler, err = store.NewIoHandler(config)
	if err != nil {
		a.AddError(err)
		return
	}
	a.AddNote("connected to", config.Database)
	return
}

func (a *Action) connectOut(config *iohandler.Accessor) (handler iohandler.OutHandler, err error) {
	a.AddNote("connecting to output", config.Database)
	handler, err = store.NewOutHandler(config)
	if err != nil {
		a.AddError(err)
		return
	}
	a.AddNote("connected to", config.Database)
	return
}

func (a *Action) cloneDatabase(output *iohandler.Accessor) (err error) {
	var dataIn iohandler.IoHandler
	var dataOut iohandler.OutHandler

	onExit := func() {
		if err != nil {
			a.AddError(err)
		}

		if dataIn != nil {
			err = dataIn.OnExit()
			if err != nil {
				a.AddError(err)
			}
		}

		if dataOut != nil {
			err = dataOut.OnExit()
			if err != nil {
				a.AddError(err)
			}
		}
	}

	defer onExit()

	dataIn, err = a.connectIn(a.Input)
	if err != nil {
		return
	}

	dataOut, err = a.connectOut(output)
	if err != nil {
		return
	}

	err = a.createDatabase(output, dataOut)
	if err != nil {
		return
	}

	_, err = dataIn.RootFolder()
	if err != nil {
		return
	}
	a.AddNote("input read from", a.Input.Database)

	a.AddNote("write database")
	err = a.WriteDatabase(dataIn, dataOut)
	if err != nil {
		return
	}
	a.AddNote("output written to", output.Database)
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

func (a *Action) verifyConfiguration(config *iohandler.Accessor, refresh bool) error {
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

func (a *Action) verifyOutConfiguration(config *iohandler.Accessor) error {
	_, err := store.NewOutHandler(config)
	if err != nil {
		return a.AddError(err)
	}
	return nil
}

func (a *Action) Verify() {
	a.verifyConfiguration(a.Input, true)
	for _, output := range a.Outputs {
		a.verifyOutConfiguration(output)
	}
}

func (a *Action) Update() {
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
		err = a.Clone()
	default:
		err = a.AddError(fmt.Errorf("unknown method %s", a.Method))
	}

	return
}
