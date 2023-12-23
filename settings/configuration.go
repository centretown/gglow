package settings

type Configuration struct {
	Driver   string
	Path     string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Folder   string

	Effect string
}

func NewConfiguration() *Configuration {
	cf := &Configuration{}
	return cf
}
