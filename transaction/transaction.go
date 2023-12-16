package transaction

import (
	"fmt"
	"glow-gui/settings"
	"glow-gui/store"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Transaction struct {
	Input   *settings.Configuration   `yaml:"input" json:"input"`
	Outputs []*settings.Configuration `yaml:"outputs" json:"outputs"`
}

type ConfigurationLog struct {
	Configuration *settings.Configuration
	Errors        []string
}

func NewConfigurationLog(config *settings.Configuration) *ConfigurationLog {
	cf := &ConfigurationLog{
		Configuration: config,
		Errors:        make([]string, 0),
	}
	return cf
}

func (cf *ConfigurationLog) AddError(err error) {
	cf.Errors = append(cf.Errors, err.Error())
}

type TransactionLog struct {
	Transaction *Transaction
	Logs        []*ConfigurationLog
}

func NewTransactionLog(transaction *Transaction) *TransactionLog {
	tl := &TransactionLog{
		Transaction: transaction,
		Logs:        make([]*ConfigurationLog, 0),
	}
	return tl
}

func (tl *TransactionLog) AddLog(cf *ConfigurationLog) {
	tl.Logs = append(tl.Logs, cf)
}

func VerifyTransactions(transactionFile string) (logs []*TransactionLog, err error) {
	logs = make([]*TransactionLog, 0)
	transactions, err := ReadTransactions(transactionFile)
	if err != nil {
		return
	}

	for _, transaction := range transactions {
		transactionLog := NewTransactionLog(transaction)

		configLog := NewConfigurationLog(transaction.Input)
		config := transaction.Input
		st, err := store.DataSource(config, nil, true)
		if err != nil {
			configLog.AddError(err)
		} else {
			st.OnExit()
		}

		transactionLog.AddLog(configLog)

		for _, config := range transaction.Outputs {
			configLog := NewConfigurationLog(config)
			st, err := store.DataSource(config, nil, false)
			if err != nil {
				configLog.AddError(err)
			} else {
				st.OnExit()
			}
			transactionLog.AddLog(configLog)
		}

		logs = append(logs, transactionLog)
	}

	return
}

func ReadTransactions(transactionFile string) (transactions []*Transaction, err error) {
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

	transactions = make([]*Transaction, 0)
	err = yaml.Unmarshal(buf, &transactions)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func ShowLogs(logs []*TransactionLog) {
	b, err := yaml.Marshal(logs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
