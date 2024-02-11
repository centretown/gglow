package glow

// The compiler currently only optimizes this form.
//
// See issue 6011.
//
// go build -gcflags='-l -v' tb.go
// go tool objdump -S -s B2I tb
// code is optimized to:  return i
//
//	asm: MOVZX AL, AX
//	     RET
//
// use for branch free conditional
// hopefully this gets inlined
func B2I(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
