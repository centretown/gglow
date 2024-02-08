package glow

import "image/color"

type Light interface {
	Get(uint16) color.NRGBA
	Set(uint16, color.NRGBA)
	Refresh()
}
