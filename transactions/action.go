package transactions

import (
	"fmt"
	"glow-gui/effects"
	"glow-gui/settings"
	"glow-gui/store"
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

func (a *Action) AddNote(note string) {
	a.Notes = append(a.Notes, note)
}

func (a *Action) AddError(err error) error {
	a.Errors = append(a.Errors, err.Error())
	return err
}

func (a *Action) HasErrors() bool {
	return len(a.Errors) > 0
}

func (action *Action) copyDatabase(output *settings.Configuration) (err error) {
	var dataIn, dataOut effects.IoHandler
	dataIn, err = store.DataSource(action.Input, nil)
	if err != nil {
		return action.AddError(err)
	}
	defer dataIn.OnExit()

	dataOut, err = store.DataSource(output, nil)
	if err != nil {
		return action.AddError(err)
	}

	action.AddNote("database created")

	err = dataOut.CreateNewDatabase()
	if err != nil {
		return action.AddError(err)
	}

	_, err = dataIn.Refresh()
	if err != nil {
		return action.AddError(err)
	}
	action.AddNote("input read")

	err = action.WriteDatabase(dataIn, dataOut)
	if err != nil {
		return action.AddError(err)
	}
	action.AddNote("output written")
	return nil
}

func (action *Action) Copy() (err error) {
	action.Verify()
	if action.HasErrors() {
		err = fmt.Errorf("action %s has errors", action.Method)
		return
	}

	for _, output := range action.Outputs {
		err = action.copyDatabase(output)
		if err != nil {
			return
		}
	}
	return
}

func (action *Action) verifyConfiguration(config *settings.Configuration, refresh bool) error {
	st, err := store.DataSource(config, nil)
	if err != nil {
		return action.AddError(err)
	}
	defer st.OnExit()

	if refresh {
		_, err = st.Refresh()
		if err != nil {
			return action.AddError(err)
		}
	}
	return nil
}

func (action *Action) Verify() {
	action.verifyConfiguration(action.Input, true)
	for _, output := range action.Outputs {
		action.verifyConfiguration(output, false)
	}
}
