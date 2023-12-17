package transactions

import (
	"fmt"
	"glow-gui/settings"
	"glow-gui/store"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Transaction struct {
	Actions []*Action
}

func NewTransaction() *Transaction {
	ta := &Transaction{
		Actions: make([]*Action, 0),
	}
	return ta
}

func (transaction *Transaction) Verify() (logs *TransactionLog) {
	logs = NewTransActionLog()

	verifyConfig := func(config *settings.Configuration) *ConfigurationLog {
		configLog := NewConfigurationLog(config)
		st, err := store.DataSource(config, nil, true)
		if err != nil {
			configLog.AddError(err)
		} else {
			st.OnExit()
		}
		return configLog
	}

	for _, transaction := range transaction.Actions {
		transactionLog := NewActionLog(transaction)
		configLog := verifyConfig(transaction.Input)
		transactionLog.AddLog(configLog)

		for _, config := range transaction.Outputs {
			configLog := NewConfigurationLog(config)
			transactionLog.AddLog(configLog)
		}
		logs.AddLog(transactionLog)
	}

	return
}

func ReadTransaction(transactionFile string) (transaction *Transaction, err error) {
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

	transaction = NewTransaction()
	err = yaml.Unmarshal(buf, transaction)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
