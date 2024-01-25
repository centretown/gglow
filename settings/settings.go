package settings

type Settings int

const (
	StripColumns Settings = iota
	StripRows
	ContentWidth
	ContentHeight
	ContentView
	GlowThemeVariant
	GlowThemeScale
	AccessFile
)

var settings = []string{
	"strip_columns",
	"strip_rows",
	"content_width",
	"content_height",
	"content_view",
	"theme_variant",
	"theme_scale",
	"access_file",
}

func (s Settings) String() string {
	return settings[s]
}
