package main

import (
	"bytes"
	"text/template"

	"github.com/Mushus/gogen/aq"
	"github.com/Mushus/gogen/goname"
	"github.com/pkg/errors"
)

type listGenerator struct {
	structName string
	params     listParams
}

type listParams struct {
	ListName  string
	Type      string
	TypeUpper string
	Fields    []listFieldParams
}

type listFieldParams struct {
	NameUpper   string
	Name        string
	Type        string
	Equalable   bool
	Compareable bool
	Pointer     bool
}

func newListGenerator(structName string) *listGenerator {
	return &listGenerator{
		structName: structName,
	}
}

func (g *listGenerator) collectParams(aqi *aq.Instance, oldGeneratedCode []byte) error {
	s := aqi.Structs().FindOne(aq.StructNameIs(g.structName))
	if !s.Exists() {
		return errors.Errorf("struct %#v not found", g.structName)
	}

	listName := g.structName + "List"
	typ := g.structName

	g.params = listParams{
		ListName:  listName,
		Type:      typ,
		TypeUpper: goname.UpperCamelCase(typ),
		Fields:    createListFieldParams(s.Fields()),
	}
	return nil
}

func createListFieldParams(fields aq.Fields) []listFieldParams {
	list := []listFieldParams{}
	for _, f := range fields {
		name := f.Name()
		typ := f.Type().Name()
		list = append(list, listFieldParams{
			NameUpper:   goname.UpperCamelCase(name),
			Name:        name,
			Type:        typ,
			Equalable:   equalableType[typ],
			Compareable: compareableType[typ],
			Pointer:     f.Type().IsPtr(),
		})
	}

	return list
}

var equalableType = map[string]bool{
	"bool":       true,
	"string":     true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"byte":       true,
	"rune":       true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
}

var compareableType = map[string]bool{
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"byte":       true,
	"rune":       true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
}

const listTmplStr = `
type {{ .ListName }} []*{{ .Type }}

func (r {{ .ListName }}) Chunk(size int) []{{ .ListName }} {
	list := []{{ .ListName }}{}
	chunk := {{.ListName}}{}
	for _, v := range r {
		chunk := append(chunk, v)
		if len(chunk) >= size {
			list = append(list, chunk)
			chunk = {{.ListName}}{}
		}
	}
	if len(chunk) > 0 {
		list = append(list, chunk)
	}
	return list
}

func (r {{ .ListName }}) Concat(list {{ .ListName }}) {{ .ListName }} {
	return append(append({{ .ListName }}{}, r...), list...)
}

func (r {{ .ListName }}) Copy() {{ .ListName }} {
	dist := make({{ .ListName }}, len(r))
	copy(dist, r)
	return dist
}

func (r {{ .ListName }}) Count() int {
	return len(r)
}

func (r {{ .ListName }}) Each(f func(i int, v *{{ .Type }})) {
	for i, v := range r {
		f(i, v)
	}
}

func (r {{ .ListName }}) Exists() bool {
	return r != nil && len(r) > 0
}

func (r {{ .ListName }}) Every(f func(i int, v *{{ .Type }}) bool) bool {
	for i, v := range r {
		if !f(i, v) {
			return false
		}
	}
	return true
}

func (r {{ .ListName }}) Filter(funcs ...func(i int, v *{{ .Type }}) bool) {{ .ListName }} {
	list := {{ .ListName }}{}
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		list = append(list, v)
	}
	return list
}

func (r {{ .ListName }}) Find(funcs ...func(i int, v *{{ .Type }}) bool) *{{ .Type }} {
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		return v
	}
	return nil
}

func (r {{ .ListName }}) FindIndex(funcs ...func(i int, v *{{ .Type }}) bool) int {
L:
	for i, v := range r {
		for _, f := range funcs {
			if !f(i, v) {
				continue L
			}
		}
		return i
	}
	return -1
}

func (r {{ .ListName }}) First() *{{ .Type }} {
	if len(r) == 0 {
		return nil
	}
	return r[0]
}

func (r {{ .ListName }}) ForPage(pageNo int, size int) {{ .ListName }} {
	rLen := len(r)
	list := make({{ .ListName }}, 0, size)
	for i, k := pageNo * size, 0; i < rLen && k < size; i, k = i+1, k+1 {
		list = append(list, r[i])
	}
	return list
}

func (r {{ .ListName }}) Get(i int) *{{ .Type }} {
	if 0 <= i && i < len(r) {
		return r[i]
	}
	return nil
}

func (r {{ .ListName }}) IsEmpty() bool {
	return len(r) == 0
}

func (r {{ .ListName }}) IsNotEmpty() bool {
	return len(r) > 0
}

func (r {{ .ListName }}) Last() *{{ .Type }} {
	if len(r) == 0 {
		return nil
	}
	return r[len(r) - 1]
}

func (r {{ .ListName }}) Reverse() {{ .ListName }} {
	list := make({{ .ListName }}, 0, len(r))
	for i := len(r) - 1; i >= 0; i-- {
		list = append(list, r[i])
	}
	return list
}

func (r {{ .ListName }}) Some(f func(i int, v *{{ .Type }}) bool) bool {
	for i, v := range r {
		if f(i, v) {
			return true
		}
	}
	return false
}

func (r {{ .ListName }}) Take(size int) {{ .ListName }} {
	if len(r) > size {
		return r
	}
	return r[:size]
}

{{- $typeUpper := .TypeUpper }}
{{- $type := .Type }}
{{- range .Fields }}
{{- if .Equalable}}

func {{ $typeUpper }}{{ .NameUpper }}Is(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} == value
	}
}

func {{ $typeUpper }}{{ .NameUpper }}IsNot(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} != value
	}
}
{{- end }}
{{- if .Compareable}}

func {{ $typeUpper }}{{ .NameUpper }}GT(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} > value
	}
}

func {{ $typeUpper }}{{ .NameUpper }}GE(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} >= value
	}
}

func {{ $typeUpper }}{{ .NameUpper }}LT(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} < value
	}
}

func {{ $typeUpper }}{{ .NameUpper }}LE(value {{ .Type }}) func(i int, v *{{ $type }}) bool {
	return func(i int, v *{{ $type }}) bool {
		return {{ if .Pointer }}*{{ end }}v.{{ .Name }} <= value
	}
}
{{- end }}
{{- end }}
`

var listTmpl = template.Must(template.New("list").Parse(listTmplStr))

func (g *listGenerator) generate() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	br := new(bytes.Buffer)
	if err := listTmpl.Execute(br, g.params); err != nil {
		return nil, errors.Wrap(err, "cannot execute template")
	}
	return br.Bytes(), nil
}
