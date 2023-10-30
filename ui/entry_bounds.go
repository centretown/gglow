package ui

type FloatEntryBounds struct {
	MinVal, MaxVal, OnVal, OffVal float64
}

func NewFloatEntryBounds(min, max, on, off float64) *FloatEntryBounds {
	if min > max {
		min, max = max, min
	}
	eb := &FloatEntryBounds{min, max, on, off}
	return eb
}

type IntEntryBounds struct {
	MinVal, MaxVal, OnVal, OffVal int
}

func NewIntEntryBounds(min, max, on, off int) *IntEntryBounds {
	if min > max {
		min, max = max, min
	}
	return &IntEntryBounds{min, max, on, off}
}
