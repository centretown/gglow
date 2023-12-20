package settings

type Configuration struct {
	Driver   string `yaml:"driver" json:"driver"`
	Path     string `yaml:"path" json:"path"`
	User     string `yaml:"-" json:"-"`
	Password string `yaml:"-" json:"-"`
	Url      string
	Database string
	Folder   string `yaml:"folder" json:"folder"`
	Effect   string `yaml:"effect" json:"effect"`
}

func NewConfiguration() *Configuration {
	cf := &Configuration{}
	return cf
}
