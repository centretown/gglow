package transactions

import "glow-gui/settings"

type ConfigurationLog struct {
	Configuration *settings.Configuration
	Errors        []string
	Notes         []string
}

func NewConfigurationLog(config *settings.Configuration) *ConfigurationLog {
	cf := &ConfigurationLog{
		Configuration: config,
		Errors:        make([]string, 0),
		Notes:         make([]string, 0),
	}
	return cf
}

func (cf *ConfigurationLog) AddError(err error) {
	cf.Errors = append(cf.Errors, err.Error())
}

func (cf *ConfigurationLog) AddNote(note string) {
	cf.Errors = append(cf.Notes, note)
}

func (cf *ConfigurationLog) HasErrors() bool {
	return len(cf.Errors) > 0
}
