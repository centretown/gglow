package glow

// The compiler currently only optimizes this form.
//
// See issue 6011.
//
// go build -gcflags='-l -v' tb.go
// go tool objdump -S -s Bool2Int tb
// code is optimized to:  return i
//
//	asm: MOVZX AL, AX
//	     RET
//
// use for branch free conditional
func Bool2Int(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
