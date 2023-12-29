package codeio

import (
	"os"
	"strings"
	"text/template"
)

type Generator struct {
	file *os.File
}

func (gen *Generator) Open(name string) (err error) {
	gen.file, err = os.Create(name)
	return
}

func (gen *Generator) Close() error {
	return gen.file.Close()
}

func (gen *Generator) makeList(folders []*FolderList) []EffectItem {
	var gen_list = make([]EffectItem, 0)
	for _, folder := range folders {
		for _, s := range folder.List {
			gen_list = append(gen_list,
				EffectItem{
					Title:    MakeTitle(folder.Title, s.Title),
					Constant: MakeConstant(folder.Title, s.Title),
					Frame:    s.Frame})
		}
	}
	return gen_list
}

type HeaderGenerator struct {
	Generator
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

func (hg *HeaderGenerator) Write(folders []*FolderList) (err error) {
	t := template.Must(template.New("header").Parse(templHeader))
	var gen_list = hg.makeList(folders)
	err = t.Execute(hg.Generator.file, gen_list)
	return
}

type SourceGenerator struct {
	Generator
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

func (sg *SourceGenerator) Write(folders []*FolderList) (err error) {
	t := template.Must(template.New("source").Parse(templSource))
	var gen_list = sg.makeList(folders)
	err = t.Execute(sg.Generator.file, gen_list)
	return
}

type EffectGenerator struct {
	Generator
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

func (eg *EffectGenerator) Write(folders []*FolderList) (err error) {
	t := template.Must(template.New("effect").Parse(templEffect))
	var gen_list = eg.makeList(folders)
	err = t.Execute(eg.Generator.file, gen_list)
	return
}

func MakeConstant(folder, title string) (s string) {
	return strings.ToUpper(strings.ReplaceAll(folder+" "+title, " ", "_"))
}

func MakeTitle(folder, title string) (s string) {
	return strings.ReplaceAll(folder+"_"+title, " ", "_")
}

// "\"" + folder.Title + ":" + s.Title + "\""
