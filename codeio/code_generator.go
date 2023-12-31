package codeio

import (
	"gglow/iohandler"
	"os"
	"strings"
	"text/template"
)

type CodeGenerator struct {
	file *os.File
}

func (gen *CodeGenerator) Open(name string) (err error) {
	gen.file, err = os.Create(name)
	return
}

func (gen *CodeGenerator) Close() error {
	return gen.file.Close()
}

func (gen *CodeGenerator) makeList(folders []*iohandler.EffectItems) []iohandler.EffectItem {
	var gen_list = make([]iohandler.EffectItem, 0)
	for _, folder := range folders {
		for _, s := range folder.List {
			gen_list = append(gen_list,
				iohandler.EffectItem{
					Title:    MakeTitle(folder.Title, s.Title),
					Constant: MakeConstant(folder.Title, s.Title),
					Frame:    s.Frame})
		}
	}
	return gen_list
}

var _ iohandler.Generator = (*HeaderGenerator)(nil)

type HeaderGenerator struct {
	CodeGenerator
}

func NewHeaderGenerator() *HeaderGenerator {
	hg := &HeaderGenerator{}
	return hg
}

const templHeader = `
// CAUTION GENERATED FILE
#pragma once
#include "Frame.h"
namespace glow {
enum CATALOG_INDEX : uint8_t {
{{range .}}{{printf "%s,\n" .Constant}}{{end}}{{printf "FRAME_COUNT"}}};
extern const char *catalog_names[FRAME_COUNT];
extern Frame catalog[FRAME_COUNT];
extern Frame &from_catalog(CATALOG_INDEX index);
extern const char *catalog_name(CATALOG_INDEX index);
} // namespace glow
`

func (hg *HeaderGenerator) Write(folders []*iohandler.EffectItems) (err error) {
	t := template.Must(template.New("header").Parse(templHeader))
	var gen_list = hg.makeList(folders)
	err = t.Execute(hg.CodeGenerator.file, gen_list)
	return
}

var _ iohandler.Generator = (*SourceGenerator)(nil)

type SourceGenerator struct {
	CodeGenerator
}

func NewSourceGenerator() *SourceGenerator {
	sg := &SourceGenerator{}
	return sg
}

const templSource = `
// CAUTION GENERATED FILE
#include "catalog.h"
namespace glow {
const char *catalog_names[FRAME_COUNT] = {
{{range .}}{{printf "\"%s\",\n" .Title}}{{end}}};
Frame catalog[FRAME_COUNT]={
{{range .}}{{.Frame.MakeCode }}{{end}}};
Frame &from_catalog(CATALOG_INDEX index){return catalog[index%FRAME_COUNT];}
const char *catalog_name(CATALOG_INDEX index){return catalog_names[index%FRAME_COUNT];}
} // namespace glow
`

func (sg *SourceGenerator) Write(folders []*iohandler.EffectItems) (err error) {
	t := template.Must(template.New("source").Parse(templSource))
	var gen_list = sg.makeList(folders)
	err = t.Execute(sg.CodeGenerator.file, gen_list)
	return
}

var _ iohandler.Generator = (*EffectGenerator)(nil)

type EffectGenerator struct {
	CodeGenerator
}

func NewEffectGenerator() *EffectGenerator {
	eg := &EffectGenerator{}
	return eg
}

const templEffect = `
{{range .}}
- addressable_lambda: 
    name: "{{.Title}}"
    update_interval: 16ms
    lambda: |-
      #include "glow/catalog.h"
      static glow::Frame frame(glow::from_catalog(glow::{{.Constant}}));
      if (initial_run) {
        frame.setup(it.size(), 4, {{.Frame.Interval}});
      }
      if (frame.is_ready()) {
        frame.spin(it);
      }
{{end}}
`

func (eg *EffectGenerator) Write(folders []*iohandler.EffectItems) (err error) {
	t := template.Must(template.New("effect").Parse(templEffect))
	var gen_list = eg.makeList(folders)
	err = t.Execute(eg.CodeGenerator.file, gen_list)
	return
}

func MakeConstant(folder, title string) (s string) {
	return strings.ToUpper(strings.ReplaceAll(folder+" "+title, " ", "_"))
}

func MakeTitle(folder, title string) (s string) {
	return strings.ReplaceAll(folder+"_"+title, " ", "_")
}

// "\"" + folder.Title + ":" + s.Title + "\""
