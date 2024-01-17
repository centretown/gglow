package ui

import (
	"gglow/fyglow/effectio"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
)

type FolderDialog struct {
	*SimpleDialog
}

func NewFolderDialog(effect *effectio.EffectIo, window fyne.Window) *FolderDialog {
	fd := &FolderDialog{}

	fd.SimpleDialog = NewSimpleDialog(effect, window,
		text.AddFolderLabel.String(), text.FolderLabel.String())

	fd.NameEntry.Validator = validation.NewAllStrings(fd.validateFolderName)
	fd.Apply = func() {
		name, _ := fd.Name.Get()
		err := fd.effect.AddFolder(name)
		if err != nil {
			fyne.LogError(name, err)
		}

	}
	return fd
}

func (fd *FolderDialog) validateFolderName(s string) error {
	err := fd.effect.ValidateNewFolderName(s)
	if err != nil {
		fd.ApplyButton.Disable()
		return err
	}
	fd.ApplyButton.Enable()
	return nil
}
