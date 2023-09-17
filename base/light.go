package glow

import "image/color"

type Light interface {
	Get(uint16) color.RGBA
	Set(uint16, color.RGBA)
}
