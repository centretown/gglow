package history

// import (
// 	"gglow/glow"
// 	"testing"

// 	"fyne.io/fyne/v2/test"
// )

// var test_intervals = []uint32{
// 	48, 100, 200, 16, 24,
// }

// func TestHistory(t *testing.T) {
// 	var err error
// 	app := test.NewApp()
// 	store := NewStore(app.Preferences())
// 	store.EffectName = "TestHistoryEffect"
// 	history := store.history

// 	current := &glow.Frame{}
// 	current.Setup(100, 10)
// 	for _, interval := range test_intervals {
// 		current.SetInterval(interval)
// 		err = history.Add(store.route, store.EffectName, current)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	t.Log(current.Interval, store.route, store.EffectName)

// 	hasHistory := history.HasHistory(store.route, store.EffectName)
// 	if !hasHistory {
// 		t.Fatal("no history")
// 	}

// 	index := len(test_intervals)
// 	for hasHistory {
// 		current, err = history.RestorePrevious(store.route, store.EffectName)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		index--
// 		if index < 0 {
// 			t.Fatal("index < 0", index)
// 		}

// 		expected := test_intervals[index]
// 		t.Log(expected, current.Interval, store.route, store.EffectName)
// 		if expected != current.Interval {
// 			t.Fatal("intervals don't match", expected, current.Interval)
// 		}

// 		hasHistory = history.HasHistory(store.route, store.EffectName)
// 	}
// }
