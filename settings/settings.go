package settings

type Settings int

const (
	StripColumns Settings = iota
	StripRows
	ContentWidth
	ContentHeight
	ContentSplit
	GlowThemeVariant
	GlowThemeScale
	Effect
	EffectPath
	EffectRoute
)

var settings = []string{
	"strip_columns",
	"strip_rows",
	"content_width",
	"content_height",
	"content_split",
	"theme_variant",
	"theme_scale",
	"effect",
	"effect_path",
	"effect_route",
}

func (s Settings) String() string {
	return settings[s]
}
