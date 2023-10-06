package data

type BindKey int

const (
	Length BindKey = iota
	Rows
	Grid
	Chroma
	HueShift
	Scan
	Begin
	End
	BIND_KEY_COUNT
)

var bindKeys = []string{
	"Length",
	"Rows",
	"Grid",
	"Chroma",
	"HueShift",
	"Scan",
	"Begin",
	"End",
}

func (k BindKey) String() string {
	return bindKeys[k]
}
