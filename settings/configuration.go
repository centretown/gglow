package settings

type Configuration struct {
	Driver string `yaml:"driver" json:"driver"`
	Path   string `yaml:"path" json:"path"`
	Folder string `yaml:"folder" json:"folder"`
	Effect string `yaml:"effect" json:"effect"`
}
