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
	StripInterval
	GlowThemeVariant
	GlowThemeScale
)

var settings = []string{
	"StripColumns",
	"StripRows",
	"StripInterval",
	"ThemeVariant",
	"ThemeScale",
}

func (s Settings) String() string {
	return settings[s]
}

const (
	StripColumnsDefault  float64           = 9
	StripRowsDefault     float64           = 4
	StripIntervalDefault float64           = 32
	VariantDefault       fyne.ThemeVariant = theme.VariantDark
	ScaleDefault         float64           = 1
)

var (
	// stripColumns  float64           = 9
	// stripRows     float64           = 4
	// stripInterval float64           = 32
	glowVariant fyne.ThemeVariant = VariantDefault
	glowScale   float64           = ScaleDefault
)

type GlowTheme struct {
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

func (m GlowTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// if name == theme.ColorNameBackground {
	// 	if variant == theme.VariantLight {
	// 		return color.White
	// 	}
	// 	return color.Black
	// }
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
