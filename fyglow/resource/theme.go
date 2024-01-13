package resource

import (
	"gglow/settings"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

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
	glowScale = preferences.FloatWithFallback(settings.GlowThemeScale.String(),
		ScaleDefault)
	glowVariant = fyne.ThemeVariant(preferences.IntWithFallback(settings.GlowThemeVariant.String(),
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
			c := color.RGBA{230, 230, 230, 255}
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
