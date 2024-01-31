package action

import (
	"fmt"
	"gglow/iohandler"
	"gglow/store"
	"strings"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Method      string
	Input       *iohandler.Accessor
	FilterItems []*FilterItem
	Outputs     []*iohandler.Accessor
	Notes       []string
	Errors      []string
	filter      Filter
}

type ActionView struct {
	Method      string
	Input       *iohandler.AccessorView
	FilterItems []*FilterItem
	Outputs     []*iohandler.AccessorView
	Notes       []string
	Errors      []string
}

func NewAction() *Action {
	a := &Action{
		Input:       &iohandler.Accessor{},
		FilterItems: make([]*FilterItem, 0),
		Outputs:     make([]*iohandler.Accessor, 0),
		Notes:       make([]string, 0),
		Errors:      make([]string, 0),
	}
	return a
}

func (a *Action) NewActionView() string {
	av := &ActionView{
		Method:      a.Method,
		Input:       iohandler.NewAccessorView(a.Input),
		FilterItems: a.FilterItems,
		Notes:       a.Notes,
		Errors:      a.Errors,
		Outputs:     make([]*iohandler.AccessorView, len(a.Outputs)),
	}
	for i := range a.Outputs {
		av.Outputs[i] = iohandler.NewAccessorView(a.Outputs[i])
	}
	buf, _ := yaml.Marshal(av)
	return string(buf)

}

func (a *Action) Process() (err error) {
	switch strings.ToLower(a.Method) {
	case "verify":
		a.Verify()
	case "update", "clone":
		err = a.Copy()
	default:
		err = a.AddError(fmt.Errorf("unknown method %s", a.Method))
	}
	return
}

func (a *Action) Copy() (err error) {
	a.Verify()
	if a.HasErrors() {
		err = fmt.Errorf("action %s has errors", a.Method)
		return
	}
	for _, output := range a.Outputs {
		err = a.updateDatabase(output)
		if err != nil {
			return
		}
	}
	return
}

func (a *Action) Verify() {
	a.verifyInput(a.Input, true)
	for _, output := range a.Outputs {
		a.verifyOutput(output)
	}
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

func (a *Action) updateDatabase(output *iohandler.Accessor) (err error) {
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

	a.AddNote("connecting to input", a.Input.Database)
	dataIn, err = store.NewIoHandler(a.Input)
	if err != nil {
		return a.AddError(err)
	}

	a.AddNote("connecting to output", output.Database)
	dataOut, err = store.NewOutHandler(output)
	if err != nil {
		return a.AddError(err)
	}

	if a.Method == "clone" {
		a.AddNote("create database", output.Database)
		err = dataOut.Create(output.Database)
		if err != nil {
			return a.AddError(err)
		}
	}

	err = a.writeDatabase(dataIn, dataOut)
	if err != nil {
		return a.AddError(err)
	}
	a.AddNote("output written to", output.Database)
	return nil
}

func (a *Action) verifyInput(config *iohandler.Accessor, refresh bool) error {
	handler, err := store.NewIoHandler(config)
	if err != nil {
		return a.AddError(err)
	}
	handler.OnExit()
	return nil
}

func (a *Action) verifyOutput(config *iohandler.Accessor) error {
	_, err := store.NewOutHandler(config)
	if err != nil {
		return a.AddError(err)
	}
	return nil
}

func (action *Action) writeDatabase(dataIn iohandler.IoHandler, dataOut iohandler.OutHandler) error {
	folders, err := dataIn.ListFolders()
	if err != nil {
		err = fmt.Errorf("dataIn ListFolders %s", err)
		return err
	}

	action.filter = NewFilter(action.FilterItems)

	for _, folder := range folders {

		if action.filter.IsSelected(folder) {
			action.AddNote(fmt.Sprintf("add folder %s", folder))
			items, err := dataIn.ListEffects(folder)
			if err != nil {
				err = fmt.Errorf("ReadFolder %s: %v", folder, err)
				return err
			}

			err = action.writeFolder(folder, items, dataIn, dataOut)
			if err != nil {
				err = fmt.Errorf("writeFolder %s: %v", folder, err)
				return err
			}
		}
	}
	return nil
}

func (action *Action) writeFolder(folder string, items []string, dataIn iohandler.IoHandler,
	dataOut iohandler.OutHandler) error {

	err := dataOut.CreateFolder(folder)
	if err != nil {
		return err
	}
	for _, item := range items {
		if !iohandler.IsFolder(item) && action.filter.IsSelected(folder, item) {
			action.AddNote(fmt.Sprintf("add effect %s.%s", folder, item))
			frame, err := dataIn.ReadEffect(folder, item)
			if err != nil {
				return err
			}

			err = dataOut.CreateEffect(folder, item, frame)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
