package data

import (
	"glow-gui/resources"

	"fyne.io/fyne/v2/data/binding"
)

func setUntyped(binder binding.Untyped, face interface{},
	msgid resources.MessageID) (err error) {

	err = binder.Set(face)
	if err != nil {
		msgid.Log("model", err)
	}
	return
}

func setUntypedList(binder binding.UntypedList,
	list []interface{}, msgid resources.MessageID) (err error) {

	err = binder.Set(list)
	if err != nil {
		resources.MsgSetLayerList.Log("model", err)
	}
	return
}
