package transactions

import "glow-gui/settings"

type Action struct {
	Method  string                    `yaml:"method" json:"method"`
	Input   *settings.Configuration   `yaml:"input" json:"input"`
	Outputs []*settings.Configuration `yaml:"outputs" json:"outputs"`
}
