package glow

import "testing"

func testGridBase(t *testing.T, grid *Grid, length, rows uint16,
	origin Origin, orientation Orientation) {
	if err := grid.Setup(length, rows, origin, orientation); err != nil {
		t.Fatalf(err.Error())
	}
	if grid.Length != length {
		t.Fatalf("Length expected 4 got %d", grid.Length)
	}
	if grid.Rows != rows {
		t.Fatalf("Rows expected 2 got %d", grid.Rows)
	}
	if grid.Origin != origin {
		t.Fatalf("Origin expected %d got %d", TopLeft, grid.Origin)
	}
	if grid.Orientation != orientation {
		t.Fatalf("Orientation expected %d got %d", Diagonal, grid.Orientation)
	}

}

func testTable(t *testing.T, grid *Grid, table []uint16) {
	var offset uint16
	for i, expected := range table {
		offset = grid.Map(uint16(i))
		if offset != expected {
			t.Fatalf("Map(%d) want %d got %d", i,
				expected, offset)
		}
	}
}

func TestGridTowByTwo(t *testing.T) {
	var grid Grid
	const length uint16 = 4
	const rows uint16 = 2

	testGridBase(t, &grid, length, rows, TopLeft, Diagonal)

	first := grid.GetFirst()
	if first != 1 {
		t.Fatalf("GetFirst expected 1 got %d", first)
	}
	offset := grid.GetOffset()
	if offset != 1 {
		t.Fatalf("GetOffset expected 1 got %d", offset)
	}

	for i := uint16(0); i < 4; i++ {
		offset = grid.Map(i)
		if offset != i {
			t.Fatalf("Map expected %d got %d", i, offset)
		}
	}
}

var fourByNineDiagonal = []uint16{
	0, 1, 9, 2, 10, 18, 3, 11, 19,
	27, 4, 12, 20, 28, 5, 13, 21, 29,
	6, 14, 22, 30, 7, 15, 23, 31, 8,
	16, 24, 32, 17, 25, 33, 26, 34, 35,
}

func TestGridFourByNineDiagonal(t *testing.T) {
	var grid Grid
	const length uint16 = 36
	const rows uint16 = 4

	testGridBase(t, &grid, length, rows, TopLeft, Diagonal)

	if len(fourByNineDiagonal) != 36 {
		t.Fatalf("expected length %d got %d", length, len(fourByNineDiagonal))
	}

	first := grid.GetFirst()
	if first != 6 {
		t.Fatalf("GetFirst expected 1 got %d", first)
	}
	offset := grid.GetOffset()
	if offset != 3 {
		t.Fatalf("GetOffset expected 1 got %d", offset)
	}
	last := grid.GetLast()
	if last != 29 {
		t.Fatalf("GetFirst expected 1 got %d", last)
	}

	testTable(t, &grid, fourByNineDiagonal)
}

func TestGridFourByNineHorizontal(t *testing.T) {
	var grid Grid
	const length uint16 = 36
	const rows uint16 = 4
	var offset uint16

	testGridBase(t, &grid, length, rows, TopLeft, Horizontal)

	for expected := uint16(0); expected < length; expected++ {
		offset = grid.Map(expected)
		if offset != expected {
			t.Fatalf("Map(%d) want %d got %d", expected,
				expected, offset)
		}
	}

	testGridBase(t, &grid, length, rows, BottomRight, Horizontal)
	for i := uint16(0); i < length; i++ {
		offset = grid.Map(i)
		expected := grid.Length - i - 1
		if offset != expected {
			t.Fatalf("Map(%d) want %d got %d", i,
				expected, offset)
		}
	}
}

var fourByNineVertical = []uint16{
	0, 9, 18, 27,
	1, 10, 19, 28,
	2, 11, 20, 29,
	3, 12, 21, 30,
	4, 13, 22, 31,
	5, 14, 23, 32,
	6, 15, 24, 33,
	7, 16, 25, 34,
	8, 17, 26, 35,
}

var fourByNineVerticalBottomLeft = []uint16{
	27, 18, 9, 0,
	28, 19, 10, 1,
	29, 20, 11, 2,
	30, 21, 12, 3,
	31, 22, 13, 4,
	32, 23, 14, 5,
	33, 24, 15, 6,
	34, 25, 16, 7,
	35, 26, 17, 8,
}

var fourByNineVerticalTopRight = []uint16{
	8, 17, 26, 35,
	7, 16, 25, 34,
	6, 15, 24, 33,
	5, 14, 23, 32,
	4, 13, 22, 31,
	3, 12, 21, 30,
	2, 11, 20, 29,
	1, 10, 19, 28,
	0, 9, 18, 27,
}

var fourByNineVerticalBottomRight = []uint16{
	35, 26, 17, 8,
	34, 25, 16, 7,
	33, 24, 15, 6,
	32, 23, 14, 5,
	31, 22, 13, 4,
	30, 21, 12, 3,
	29, 20, 11, 2,
	28, 19, 10, 1,
	27, 18, 9, 0,
}

func TestGridFourByNineVertical(t *testing.T) {
	var grid Grid
	const length uint16 = 36
	const rows uint16 = 4

	testGridBase(t, &grid, length, rows, TopLeft, Vertical)
	testTable(t, &grid, fourByNineVertical)

	testGridBase(t, &grid, length, rows, BottomLeft, Vertical)
	testTable(t, &grid, fourByNineVerticalBottomLeft)

	testGridBase(t, &grid, length, rows, TopRight, Vertical)
	testTable(t, &grid, fourByNineVerticalTopRight)

	testGridBase(t, &grid, length, rows, BottomRight, Vertical)
	testTable(t, &grid, fourByNineVerticalBottomRight)
}
