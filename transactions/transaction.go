package transactions

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Transaction struct {
	Actions []*Action
}

func NewTransaction() *Transaction {
	tr := &Transaction{
		Actions: make([]*Action, 0),
	}
	return tr
}

func ReadTransaction(transactionFile string) (tr *Transaction, err error) {
	var file *os.File
	file, err = os.Open(transactionFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	tr = NewTransaction()
	err = yaml.Unmarshal(buf, tr)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (tr *Transaction) ShowLogs() {
	b, err := yaml.Marshal(tr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	if tr.HasErrors() {
		fmt.Println("Transaction has errors")
	}
}

func (tl *Transaction) HasErrors() bool {
	for _, a := range tl.Actions {
		if a.HasErrors() {
			return true
		}
	}
	return false
}

// func (tr *Transaction) verifyConfig(config *settings.Configuration, refresh bool) {
// 	st, err := store.DataSource(config, nil)
// 	if err != nil {
// 		config.AddError(err)
// 		return
// 	}
// 	defer st.OnExit()

// 	if refresh {
// 		_, err = st.Refresh()
// 		if err != nil {
// 			config.AddError(err)
// 		}
// 	}
// }

// func (tr *Transaction) Verify() {
// 	for _, action := range tr.Actions {
// 		tr.verifyConfig(action.Input, true)

// 		for _, output := range action.Outputs {
// 			tr.verifyConfig(output, false)
// 		}
// 	}
// }

// func (tr *Transaction) copyDatabase(action *Action, output *settings.Configuration) (err error) {
// 	var dataIn, dataOut effects.IoHandler
// 	dataIn, err = store.DataSource(action.Input, nil)
// 	if err != nil {
// 		action.Input.AddError(err)
// 		return err
// 	}
// 	action.AddNote("source action.Input")
// 	defer dataIn.OnExit()

// 	dataOut, err = store.DataSource(output, nil)
// 	if err != nil {
// 		output.AddError(err)
// 	}

// 	err = dataOut.CreateNewDatabase()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = dataIn.Refresh()
// 	if err != nil {
// 		return err
// 	}

// 	err = WriteDatabase(dataIn, dataOut)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (tr *Transaction) Process() (err error) {
	for _, action := range tr.Actions {
		switch action.Method {
		case "verify":
			action.Verify()
		case "copy":
			action.Copy()
		}
	}
	return
}
