package settings

import (
	"flag"
)

type Configuration struct {
	Method string `yaml:"method" json:"method"`
	Path   string `yaml:"path" json:"path"`
	Folder string `yaml:"folder" json:"folder"`
	Effect string `yaml:"effect" json:"effect"`
}

const (
	methodDefault = "sqlite3"
	methodUsage   = "storage method (sqlite3, mysql, file)"
	pathUsage     = "path to data"
	pathDefault   = ""
	folderUsage   = "folder to access"
	folderDefault = ""
	effectUsage   = "effect to read"
	effectDefault = ""
)

func ParseCommandLine(config *Configuration) {
	flag.StringVar(&config.Method, "s", methodDefault, methodUsage+" (short form)")
	flag.StringVar(&config.Method, "storage", methodDefault, methodUsage)
	flag.StringVar(&config.Path, "p", pathDefault, pathUsage+" (short form)")
	flag.StringVar(&config.Path, "path", pathDefault, pathUsage)
	flag.StringVar(&config.Folder, "f", folderDefault, folderUsage+" (short form)")
	flag.StringVar(&config.Folder, "folder", folderDefault, folderUsage)
	flag.StringVar(&config.Effect, "e", effectDefault, effectUsage+" (short form)")
	flag.StringVar(&config.Effect, "effect", effectDefault, effectUsage)
}
