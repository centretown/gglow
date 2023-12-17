package transactions

type ActionLog struct {
	Action *Action
	Notes  []string
	Logs   []*ConfigurationLog
}

func NewActionLog(transaction *Action) *ActionLog {
	tl := &ActionLog{
		Action: transaction,
		Logs:   make([]*ConfigurationLog, 0),
		Notes:  make([]string, 0),
	}
	return tl
}

func (al *ActionLog) AddLog(cf *ConfigurationLog) {
	al.Logs = append(al.Logs, cf)
}

func (al *ActionLog) AddNote(note string) {
	al.Notes = append(al.Notes, note)
}

func (al *ActionLog) HasErrors() bool {
	for _, cf := range al.Logs {
		if cf.HasErrors() {
			return true
		}
	}
	return false
}
