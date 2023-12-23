package effects

// import (
// 	"fmt"
// 	"gglow/fileio"
// 	"gglow/glow"
// 	"testing"

// 	"fyne.io/fyne/v2/test"
// )

// func showStatus(model *Model) string {
// 	frame := model.GetFrame()
// 	return fmt.Sprintf("EffectName: %v\n IsDirty: %v\n CanUndo: %v\n Interval: %v",
// 		model.EffectName(), model.IsDirty(), model.CanUndo(), frame.Interval)
// }

// func TestModel(t *testing.T) {
// 	var err error
// 	app := test.NewApp()
// 	store := fileio.NewStore(app.Preferences())
// 	model := NewModel(store)

// 	t.Log("Nothing", showStatus(model))

// 	title := model.KeyList()[0]
// 	err = model.ReadEffect(title)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	model.WindowHasContent = true
// 	t.Log("ReadEffect", showStatus(model))

// 	incrementInterval := func(fr *glow.Frame) {
// 		fr.Interval++
// 		model.SetDirty()
// 		t.Log("incrementInterval", showStatus(model))
// 	}

// 	model.AddSaveAction(incrementInterval)

// 	err = model.WriteEffect()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("WriteEffect", showStatus(model))

// 	for index := 0; model.CanUndo(); index++ {
// 		model.UndoEffect()
// 		t.Logf("undo %v \n%v",
// 			index, showStatus(model))
// 	}
// }
