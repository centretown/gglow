package ui

type EntryBoundsFloat struct {
	MinVal, MaxVal, OnVal, OffVal float64
}

func NewEntryBoundsFloat(min, max, on, off float64) *EntryBoundsFloat {
	if min > max {
		min, max = max, min
	}
	eb := &EntryBoundsFloat{min, max, on, off}
	return eb
}

type EntryBoundsInt struct {
	MinVal, MaxVal, OnVal, OffVal int
}

func NewEntryBoundsInt(min, max, on, off int) *EntryBoundsInt {
	if min > max {
		min, max = max, min
	}
	return &EntryBoundsInt{min, max, on, off}
}
