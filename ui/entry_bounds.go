package ui

type EntryBounds struct {
	MinVal, MaxVal, OnVal, OffVal float64
}

func NewEntryBounds(min, max, on, off float64) *EntryBounds {
	if min > max {
		min, max = max, min
	}
	eb := &EntryBounds{min, max, on, off}
	return eb
}
