package glow

type HSVMask struct {
	H, S, V int
}

func (hm *HSVMask) None() bool { return hm.H|hm.S|hm.V == 0 }

// func (hm *HSVMask) Apply(hsv HSVColor, index, length uint16) HSVColor {
// 	var hue int = int(hsv.Hue) + hm.Hue*int(index)/int(length)
// 	hue += B2I(hue < 0) * int(hue_limit)
// 	var sat int = int(hsv.Saturation) + hm.Saturation*int(index)/int(length)
// 	return HSVColor{
// 		Hue:        uint16(hue),
// 		Saturation: 0,
// 		Value:      0,
// 	}
// }

func (hm *HSVMask) NextHue(hue, index, length uint16) uint16 {
	var next int = int(hue) + hm.H*int(index)/int(length)
	next += B2I(next < 0) * int(hue_limit)
	next %= int(hue_limit)
	return uint16(next)
}

func (hm *HSVMask) NextSaturation(saturation uint8, index, length uint16) uint8 {
	return uint8(int(saturation) + hm.S*int(index)/int(length))
}
func (hm *HSVMask) NextValue(value uint8, index, length uint16) uint8 {
	return uint8(int(value) + hm.V*int(index)/int(length))
}

type HSVGradient struct {
	Mask HSVMask
}
