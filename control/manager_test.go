package control

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	// err := store.Setup()
	// if err != nil {
	// 	t.Fatalf(err.Error())
	// }

	// app := test.NewApp()
	// window := app.NewWindow("test")
	// model := NewModel()

	// // sliderHueShift := widget.NewSliderWithData(-100, 100, model.Fields.HueShift)
	// // sliderScan := widget.NewSliderWithData(-100, 100, model.Fields.Scan)
	// selectOrigin := widget.NewSelect(resources.OriginLabels, func(s string) {})
	// selectOrigin.OnChanged = func(s string) {
	// 	model.Fields.Origin.Set(selectOrigin.SelectedIndex())
	// }
	// model.Fields.Origin.AddListener(binding.NewDataListener(func() {
	// 	index, _ := model.Fields.Origin.Get()
	// 	selectOrigin.SetSelectedIndex(index)
	// }))

	// vbox := container.NewVBox(sliderHueShift, sliderScan, selectOrigin)
	// window.SetContent(vbox)

	// lookUpList := store.LookUpList()

	// doLogAndTest := func() {
	// 	t.Log()
	// 	hueShift, _ := model.Fields.HueShift.Get()
	// 	t.Log(resources.HueShiftLabel, hueShift)
	// 	t.Log("hue slider.Value", sliderHueShift.Value)
	// 	if hueShift != sliderHueShift.Value {
	// 		t.Errorf("hueShift %f doesn't match slider value %f",
	// 			hueShift, sliderHueShift.Value)
	// 	}

	// 	scan, _ := model.Fields.Scan.Get()
	// 	t.Log(resources.ScanLabel, scan)
	// 	t.Log("scan slider.Value", sliderScan.Value)
	// 	if scan != sliderScan.Value {
	// 		t.Errorf("scan %f doesn't match slider value %f",
	// 			hueShift, sliderScan.Value)
	// 	}

	// 	origin, _ := model.Fields.Origin.Get()
	// 	t.Log(resources.OriginLabel, origin)
	// 	t.Log("origin selection", selectOrigin.SelectedIndex())
	// }

	// for _, s := range lookUpList {
	// 	t.Log()
	// 	t.Log(s)
	// 	err = model.LoadFrame(s)
	// 	if err != nil {
	// 		t.Fatalf(err.Error())
	// 	}

	// 	for i := 0; i < model.LayerList.Length(); i++ {
	// 		model.SetCurrentLayer(i)
	// 		time.Sleep(time.Duration(time.Duration.Nanoseconds(1)))
	// 		doLogAndTest()
	// 	}
	// }
}
