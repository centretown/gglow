package ui

import (
	"encoding/json"
	"fmt"
	"glow-gui/glow"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Draggable = (*ColorPatch)(nil)
var _ fyne.Tappable = (*ColorPatch)(nil)
var _ fyne.Focusable = (*ColorPatch)(nil)
var _ fyne.Widget = (*ColorPatch)(nil)
var _ desktop.Mouseable = (*ColorPatch)(nil)
var _ desktop.Hoverable = (*ColorPatch)(nil)
var _ fyne.Shortcutable = (*ColorPatch)(nil)

// var _ desktop.Keyable = (*ColorPatch)(nil)

type ColorPatch struct {
	widget.BaseWidget
	tapped     func() `json:"-"`
	rectangle  *canvas.Rectangle
	background *canvas.Rectangle
	unused     bool
	colorHSV   glow.HSV

	hovered, focused bool
	isDirty          binding.Bool
}

func NewColorPatch(isDirty binding.Bool) (patch *ColorPatch) {
	var hsv glow.HSV
	hsv.FromColor(theme.DisabledColor())
	patch = NewColorPatchWithColor(isDirty, hsv, nil)
	patch.unused = true
	return
}

func NewColorPatchWithColor(isDirty binding.Bool, hsv glow.HSV, tapped func()) *ColorPatch {
	cp := &ColorPatch{
		background: canvas.NewRectangle(theme.ButtonColor()),
		rectangle:  canvas.NewRectangle(hsv.ToRGB()),
		tapped:     tapped,
		isDirty:    isDirty,
	}
	cp.colorHSV = hsv
	cp.ExtendBaseWidget(cp)
	return cp
}

func (cp *ColorPatch) applyPatchTheme() {
	cp.background.FillColor = cp.backgroundColor()
	cp.Refresh()
}

func (cp *ColorPatch) backgroundColor() (c color.Color) {
	switch {
	case cp.focused:
		c = theme.FocusColor()
	case cp.hovered:
		c = theme.HoverColor()
	default:
		c = theme.ButtonColor()
	}
	return
}

func (cp *ColorPatch) copy() string {
	buf, err := json.Marshal(cp.colorHSV)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (cp *ColorPatch) paste(s string) {
	if len(s) < 1 {
		return
	}

	b := []byte(s)
	var hsv glow.HSV
	err := json.Unmarshal(b, &hsv)
	if err != nil {
		return
	}
	cp.isDirty.Set(true)
	cp.SetHSVColor(hsv)
}

// Shortcutable
func (cp *ColorPatch) TypedShortcut(sc fyne.Shortcut) {
	switch sc.ShortcutName() {
	case "Copy":
		p, ok := sc.(*fyne.ShortcutCopy)
		if ok {
			p.Clipboard.SetContent(cp.copy())
		}
	case "Paste":
		p, ok := sc.(*fyne.ShortcutPaste)
		if ok {
			cp.paste(p.Clipboard.Content())
		}
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (cp *ColorPatch) MouseIn(*desktop.MouseEvent) {
	cp.hovered = true
	cp.applyPatchTheme()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (cp *ColorPatch) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (cp *ColorPatch) MouseOut() {
	cp.hovered = false
	cp.applyPatchTheme()
}

// Mouseable
func (cp *ColorPatch) MouseUp(*desktop.MouseEvent) {
}
func (cp *ColorPatch) MouseDown(*desktop.MouseEvent) {
} // Mouseable

// Draggable
func (cp *ColorPatch) Dragged(d *fyne.DragEvent) {
	// pos := d.Position

}

func (cp *ColorPatch) DragEnd() {
	fmt.Println("DragEnd")
} // Draggable

// fyne.Focusable
func (cp *ColorPatch) TypedRune(rune) {
}

func (cp *ColorPatch) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeySpace:
		cp.Tapped(nil)
	case fyne.KeyC:
	}
}

func (cp *ColorPatch) FocusGained() {
	cp.focused = true
	cp.applyPatchTheme()
}

func (cp *ColorPatch) FocusLost() {
	cp.focused = false
	cp.applyPatchTheme()
} // fyne.Focusable

func (cp *ColorPatch) SetTapped(tapped func()) {
	cp.tapped = tapped
}

func (cp *ColorPatch) SetUnused(b bool) {
	cp.unused = b
	cp.setFill(theme.DisabledColor())
}

func (cp *ColorPatch) Unused() bool {
	return cp.unused
}

func (cp *ColorPatch) GetHSVColor() glow.HSV {
	return cp.colorHSV
}

func (cp *ColorPatch) GetColor() color.Color {
	return cp.rectangle.FillColor
}

func (cp *ColorPatch) CopyPatch(source *ColorPatch) {
	cp.unused = source.unused
	cp.colorHSV = source.colorHSV
	if cp.unused {
		cp.SetUnused(true)
	} else {
		cp.setFill(cp.colorHSV.ToRGB())
	}
}

func (cp *ColorPatch) SetHSVColor(hsv glow.HSV) {
	cp.unused = false
	cp.colorHSV = hsv
	cp.setFill(hsv.ToRGB())
}

func (cp *ColorPatch) SetColor(color color.Color) {
	cp.unused = false
	cp.colorHSV.FromColor(color)
	cp.setFill(color)
}

func (cp *ColorPatch) setFill(color color.Color) {
	cp.rectangle.FillColor = color
	cp.rectangle.Refresh()
}

func (cp *ColorPatch) Tapped(_ *fyne.PointEvent) {
	if cp.tapped != nil {
		cp.tapped()
	}
}

func (cp *ColorPatch) TappedSecondary(_ *fyne.PointEvent) {
	cp.requestFocus()
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()

	copyItem := fyne.NewMenuItem("Copy", func() {
		cp.TypedShortcut(&fyne.ShortcutCopy{Clipboard: clipboard})
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		cp.TypedShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	})

	unuseItem := fyne.NewMenuItem("Unuse", func() {
		cp.SetUnused(true)
	})

	menu := fyne.NewMenu("", copyItem, pasteItem, unuseItem)
	c := fyne.CurrentApp().Driver().CanvasForObject(cp)
	popUp := widget.NewPopUpMenu(menu, c)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cp)
	popUp.ShowAtPosition(pos.Add(fyne.Delta{DX: 0, DY: cp.Size().Height}))

}

type patchRenderer struct {
	objects []fyne.CanvasObject
	patch   *ColorPatch
}

func (cp *ColorPatch) requestFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(cp); c != nil {
		c.Focus(cp)
	}
}

func (cp *ColorPatch) CreateRenderer() fyne.WidgetRenderer {
	pr := &patchRenderer{
		objects: []fyne.CanvasObject{cp.rectangle},
		patch:   cp,
	}
	cp.ExtendBaseWidget(cp)
	return pr
}

func (pr *patchRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{pr.patch.background, pr.patch.rectangle}
}

func (pr *patchRenderer) Destroy() {}

func (pr *patchRenderer) Refresh() {
	pr.Layout(pr.patch.Size())
}

func (pr *patchRenderer) Layout(size fyne.Size) {
	pr.patch.background.Resize(size)
	if pr.patch.hovered || pr.patch.focused {
		diff := theme.Padding() * 2
		vec := fyne.Delta{DX: diff / 2, DY: diff / 2}
		rectPos := pr.patch.background.Position().Add(vec)
		pr.patch.rectangle.Move(rectPos)
		size = size.SubtractWidthHeight(diff, diff)
	} else {
		pr.patch.rectangle.Move(pr.patch.background.Position())
	}

	pr.patch.rectangle.Resize(size)
}

func (pr *patchRenderer) MinSize() fyne.Size {
	return fyne.NewSquareSize(theme.IconInlineSize() + theme.Padding())
}
