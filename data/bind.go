package data

import (
	"glow-gui/res"

	"fyne.io/fyne/v2/data/binding"
)

func getUntyped(src string, binder binding.Untyped,
	msgid res.MessageID) (face interface{}, err error) {

	face, err = binder.Get()
	if err != nil {
		msgid.Log(src, err)
	}
	return
}

func setUntyped(binder binding.Untyped, face interface{},
	msgid res.MessageID) (err error) {

	err = binder.Set(face)
	if err != nil {
		msgid.Log("model", err)
	}
	return
}

func setUntypedList(binder binding.UntypedList,
	list []interface{}, msgid res.MessageID) (err error) {

	err = binder.Set(list)
	if err != nil {
		res.MsgSetLayerList.Log("model", err)
	}
	return
}
