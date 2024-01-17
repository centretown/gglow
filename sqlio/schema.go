package sqlio

import (
	"fmt"
	"gglow/iohandler"
)

func init() {
	for _, schema := range Schemas {
		schema.ListFolder = schema.listFolder
		schema.SelectFolder = schema.setFolder
		schema.WriteEffect = schema.writeEffect
		schema.ExistsEffect = schema.existsEffect
		schema.ReadEffect = schema.readEffect

		// additions and overrides
		switch schema.Version.Major {
		case 1:
		}
		SchemaMap[schema.Version.ToUint64()] = schema
	}
}

type Schema struct {
	Version   iohandler.Version
	AlterSQL  []string
	CreateSQL []string
	DropSQL   []string

	ListFolder   func(folder string) (query string)
	SelectFolder func(folder string) (query string)
	ExistsEffect func(folder, title string) (query string)
	ReadEffect   func(folder, title string) (query string)
	WriteEffect  func(update bool, items ...string) (query string)

	selectFolderSQL string
	listEffectsSQL  string
	readEffectSQL   string
	existsEffectSQL string
	insertEffectSQL string
	updateEffectSQL string
}

var SchemaMap = make(map[uint64]*Schema)
var Schemas = []*Schema{
	schema_v0,
	schema_v1,
}

func (schema *Schema) listFolder(folder string) string {
	if folder == "" || folder == iohandler.DOTS {
		return schema.selectFolderSQL
	}
	return fmt.Sprintf(schema.listEffectsSQL, folder)
}

func (schema *Schema) setFolder(folder string) string {
	if folder == "" || folder == iohandler.DOTS {
		return schema.selectFolderSQL
	}
	return fmt.Sprintf(schema.listEffectsSQL, folder)
}

func (schema *Schema) writeEffect(update bool, items ...string) string {
	if len(items) < 3 {
		return ""
	}
	folder, title, source := items[0], items[1], items[2]
	if update {
		return fmt.Sprintf(schema.updateEffectSQL, source, folder, title)
	}
	return fmt.Sprintf(schema.insertEffectSQL,
		folder, title, source)
}

func (schema *Schema) existsEffect(folder, title string) string {
	return fmt.Sprintf(schema.existsEffectSQL, folder, title)
}

func (schema *Schema) readEffect(folder, title string) string {
	return fmt.Sprintf(schema.readEffectSQL, folder, title)
}
