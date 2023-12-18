package glow

import (
	"fmt"
	"math"
)

type Orientation uint16

const (
	Horizontal Orientation = iota
	Vertical
	Diagonal
	// Centred
	ORIENTATION_COUNT
)

type Origin uint16

const (
	TopLeft Origin = iota
	TopRight
	BottomLeft
	BottomRight
	ORIGIN_COUNT
)

type Grid struct {
	Length      uint16      `yaml:"length" json:"length"`
	Rows        uint16      `yaml:"rows" json:"rows"`
	Origin      Origin      `yaml:"origin" json:"origin"`
	Orientation Orientation `yaml:"orientation" json:"orientation"`

	columns   uint16
	firstEdge uint16
	lastEdge  uint16
	centre    uint16
}

func (grid *Grid) GetFirst() uint16 {
	return grid.firstEdge
}
func (grid *Grid) GetLast() uint16 {
	return grid.lastEdge
}
func (grid *Grid) GetOffset() uint16 {
	return grid.centre
}

func (grid *Grid) Setup(length, rows uint16,
	origin Origin, orientation Orientation) error {
	grid.Length = length
	grid.Rows = rows
	grid.Origin = origin
	grid.Orientation = orientation
	return grid.Validate()
}

func (grid *Grid) SetupLength(length, rows uint16) error {
	grid.Length = length
	grid.Rows = rows
	return grid.Validate()
}

func (grid *Grid) Validate() error {
	if grid.Length == 0 {
		return fmt.Errorf("Grid.Setup zero length")
	}

	if grid.Rows == 0 {
		grid.Rows = 1
	}

	grid.columns = grid.Length / grid.Rows
	grid.setupDiagonal()
	return nil
}

func (grid *Grid) setupDiagonal() {
	grid.firstEdge = 0
	lesser := min(grid.Rows, grid.columns)
	for i := uint16(0); i < lesser; i++ {
		grid.firstEdge += i
	}
	grid.centre = lesser - 1
	grid.lastEdge = grid.firstEdge +
		(grid.columns-lesser)*grid.Rows + grid.Rows - 1
}

func (grid *Grid) Map(index uint16) uint16 {
	var offset uint16 = index
	switch grid.Orientation {
	case Diagonal:
		offset = grid.mapDiagonal(index)
	case Vertical:
		offset = grid.mapColumns(index)
	}
	return grid.mapToOrigin(offset)
}

func (grid *Grid) mapColumns(index uint16) uint16 {
	quot := index / grid.Rows
	rem := index % grid.Rows
	return rem*grid.columns + quot
}

func (grid *Grid) mapToOrigin(offset uint16) uint16 {
	if grid.Origin == TopLeft {
		return offset
	}

	if grid.Origin == BottomRight {
		return grid.Length - offset - 1
	}

	quot := offset / grid.columns
	rem := offset % grid.columns

	if grid.Origin == BottomLeft {
		return (grid.Rows-quot-1)*grid.columns + rem
	}

	if grid.Origin == TopRight {
		return quot*grid.columns + grid.columns - rem - 1
	}

	return offset
}

func (grid *Grid) mapDiagonal(index uint16) uint16 {
	if grid.columns < 3 {
		return index
	}

	if index < grid.firstEdge {
		return grid.mapDiagonalTop(index)
	}

	if index <= grid.lastEdge {
		return grid.mapDiagonalMiddle(index)
	}

	return grid.mapDiagonalBottom(index)
}

func (grid *Grid) mapDiagonalTop(index uint16) uint16 {
	var offset, start uint16
	for start < index {
		offset++
		start += offset
	}

	if start == index {
		return offset
	}

	start -= offset
	offset--
	return offset + (index-start)*(grid.columns-1)
}

func (grid *Grid) mapDiagonalMiddle(index uint16) uint16 {
	i := index - grid.firstEdge
	quot := i / grid.Rows
	rem := i % grid.Rows
	return grid.centre + quot + rem*(grid.columns-1)
}

func (grid *Grid) mapDiagonalBottom(index uint16) uint16 {
	offset := grid.columns*2 - 1
	start := grid.lastEdge + 1
	increment := grid.Rows

	for start < index {
		offset += grid.columns
		increment--
		start += increment
	}

	if start == index {
		return offset
	}

	start -= increment
	offset -= grid.columns
	return offset + (index-start)*(grid.columns-1)
}

func (grid *Grid) AdjustBounds(bound float32) uint16 {
	scaled := uint16(math.Round(float64(bound)))
	if grid.Orientation == Horizontal {
		return scaled / grid.columns * grid.columns
	}
	return scaled / grid.Rows * grid.Rows
}

func (grid *Grid) MakeCode() string {
	s := fmt.Sprintf("{%d,%d,%d,%d}",
		grid.Length, grid.Rows, grid.Origin, grid.Orientation)
	return s
}
