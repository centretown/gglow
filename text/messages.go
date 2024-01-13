package text

import (
	"fmt"
)

type MessageID int

const (
	MsgGetEffectLookup MessageID = iota
	MsgGetEffectLoad
	MsgGetFrame
	MsgGetLayer
	MsgGetTitle
	MsgSetFrame
	MsgSetLayer
	MsgSetTitle
	MsgSetLayerList
	MsgSetupStore
	MsgParseEffectPath
	MsgPathNotFolder
	MsgNoList
	MsgDuplicate
	MsgRequired
	MsgNotWritable
	MsgFirstUpper
	MsgAlphaNumeric
	MsgNotFound
	MsgListEmpty
)

var invalidMessage = "unknown Message ID"

var messages = []string{
	"unable to lookup effect",
	"unable to load effect",
	"unable to get frame",
	"unable to get layer",

	"unable to get title",
	"unable to set frame",
	"unable to set layer",
	"unable to set title",
	"unable to set layer list",

	"unable to setup storage",
	"unable to parse effect path",
	"effect path not a folder",
	"unable to list folder",
	"name already used",
	"name required",
	"name not writeable",
	"uppercase letter required",
	"must be letter or number",

	"not found", "list empty",
}

func (id MessageID) String() string {
	if int(id) < len(messages) {
		return messages[id]
	}
	return invalidMessage
}

func (id MessageID) Format(tag string) string {
	return fmt.Sprintf("%s %s", id, tag)
}
