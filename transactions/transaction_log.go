package transactions

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type TransactionLog struct {
	Logs []*ActionLog
}

func NewTransActionLog() *TransactionLog {
	tl := &TransactionLog{
		Logs: make([]*ActionLog, 0),
	}
	return tl
}

func (tl *TransactionLog) HasErrors() bool {
	for _, al := range tl.Logs {
		if al.HasErrors() {
			return true
		}
	}
	return false
}

func (tl *TransactionLog) AddLog(alog *ActionLog) {
	tl.Logs = append(tl.Logs, alog)
}
func (tl *TransactionLog) ShowLogs() {
	b, err := yaml.Marshal(tl)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
	if tl.HasErrors() {
		fmt.Println("Transaction has errors!")
	}
}
