package resources

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

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
}

func (s Settings) String() string {
	return settings[s]
}

const (
	StripColumnsDefault int               = 9
	StripRowsDefault    int               = 4
	VariantDefault      fyne.ThemeVariant = theme.VariantDark
	ScaleDefault        float64           = 1
)

var (
	glowVariant fyne.ThemeVariant = VariantDefault
	glowScale   float64           = ScaleDefault
)

type GlowTheme struct{}

func (m GlowTheme) GetVariant() fyne.ThemeVariant {
	return glowVariant
}

func NewGlowTheme(preferences fyne.Preferences) *GlowTheme {
	glowScale = preferences.FloatWithFallback(GlowThemeScale.String(),
		ScaleDefault)
	glowVariant = fyne.ThemeVariant(preferences.IntWithFallback(GlowThemeVariant.String(),
		int(VariantDefault)))
	if glowVariant == theme.VariantDark {
		LoadIcons("dark")
	} else {
		LoadIcons("light")
	}

	return &GlowTheme{}
}

const LightStripBackground fyne.ThemeColorName = "LightStripBackground"

func (m GlowTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == LightStripBackground {
		if variant == theme.VariantLight {
			c := color.RGBA{32, 48, 96, 255}
			return c
		}
		c := color.RGBA{24, 12, 8, 255}
		return c
	}
	return theme.DefaultTheme().Color(name, glowVariant)
}

func (m GlowTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m GlowTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// if name == theme.IconNameHome {
	// 	fyne.NewStaticResource("myHome", homeBytes)
	// }

	return theme.DefaultTheme().Icon(name)
}

func (m GlowTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) * float32(glowScale)
}
