package effectio

import "fyne.io/fyne/v2/data/binding"

func Alert(x binding.Int) {
	a, _ := x.Get()
	x.Set(a + 1)
}
