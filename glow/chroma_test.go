package glow

import "testing"

func testChromaBase(t *testing.T, chroma *Chroma,
	length uint16, source HSV, target HSV, hueShift int16) {

	if err := chroma.Setup(length, source, target, hueShift); err != nil {
		t.Fatalf(err.Error())
	}

	if chroma.Length != length {
		t.Fatalf("chroma.length got %d want %d", chroma.Length, length)
	}

	if len(chroma.Colors) != 2 {
		t.Fatalf("chroma.Colord.length got %d want 2", len(chroma.Colors))
	}

	if source != chroma.Colors[0] {
		t.Fatalf("chroma.Colors[0] got %v want %v", chroma.Colors[0], source)
	}

	if target != chroma.Colors[1] {
		t.Fatalf("chroma.Colors[1] got %v want %v", chroma.Colors[1], target)
	}

	if chroma.HueShift != hueShift {
		t.Fatalf("chroma.hueShift got %d want %d", chroma.HueShift, hueShift)
	}
}

func testChromaColors(t *testing.T) {

	var chroma Chroma
	type colorList []HSV
	colorsLists := []colorList{
		{{180, 1, 1}},
		{{180, 1, 1}, {160, 1, 1}},
		{{180, 1, 1}, {160, 1, 1}, {140, 1, 1}},
		{{180, 1, 1}, {160, 1, 1}, {140, 1, 1}, {120, 1, 1}},
		{{180, 1, 1}, {160, 1, 1}, {140, 1, 1}, {120, 1, 1}, {100, 1, 1}},
	}

	colorMap := func(index uint16) {
		t.Log(index)
		c := chroma.Map(index)
		t.Log(c)
	}

	for i, list := range colorsLists {
		chroma.Colors = list
		t.Log("index", i, list)
		chroma.SetupLength(100, -1)
		colorMap(0)
		colorMap(33)
		colorMap(50)
		colorMap(99)
	}

}

func TestChroma(t *testing.T) {
	var chroma Chroma
	testChromaBase(t, &chroma, 10, HSV{0, 1, 1}, HSV{180, 1, 1}, 1)
	testChromaColors(t)
}
