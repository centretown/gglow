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

func (tr *Transaction) Process() (err error) {
	for _, action := range tr.Actions {
		action.Process()
	}
	return
}
