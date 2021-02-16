package main

import (
	"bytes"
	"text/template"

	"github.com/Mushus/gogen/aq"
	"github.com/Mushus/gogen/goname"
	"github.com/jinzhu/inflection"
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
	Getters   []listGetterParams
}

type listFieldParams struct {
	NameUpper       string
	NameUpperPlural string
	Name            string
	Type            string
	Equalable       bool
	Compareable     bool
	Pointer         bool
}

type listGetterParams struct {
	ListFuncName string
	ListType     string
	Name         string
}

func newListGenerator(structName string) *listGenerator {
	return &listGenerator{
		structName: structName,
	}
}

func (g *listGenerator) collectParams(aqi *aq.AQ, oldGeneratedCode []byte) error {
	s := aqi.Structs().Find(aq.StructNameIs(g.structName))
	if !s.Exists() {
		return errors.Errorf("struct %#v not found", g.structName)
	}

	sn := g.structName
	listName := inflection.Plural(sn)
	if sn == listName {
		listName = sn + "List"
	}
	typ := g.structName

	g.params = listParams{
		ListName:  listName,
		Type:      typ,
		TypeUpper: goname.UpperCamelCase(typ),
		Fields:    createListFieldParams(s.Fields()),
		Getters:   createListGetterParams(aqi, listName, typ),
	}
	return nil
}

func createListFieldParams(fields aq.Fields) []listFieldParams {
	list := []listFieldParams{}
	for _, f := range fields {
		getter := newTagValue(f.Tag().Body().Get("getter"))
		if getter.contains("-") {
			continue
		}

		name := f.Name()
		typ := f.Type().UnwrapPtr().GoCode()

		nameUpper := goname.UpperCamelCase(name)
		nameUpperPlural := inflection.Plural(nameUpper)
		if nameUpperPlural == nameUpper {
			nameUpperPlural = nameUpper + "List"
		}

		list = append(list, listFieldParams{
			NameUpper:       nameUpper,
			NameUpperPlural: nameUpperPlural,
			Name:            name,
			Type:            f.Type().GoCode(),
			Equalable:       equalableType[typ],
			Compareable:     compareableType[typ],
			Pointer:         f.Type().IsPtr(),
		})
	}

	return list
}

func createListGetterParams(aqi *aq.AQ, listName string, typ string) []listGetterParams {

	funcs := aqi.Funcs()
	getters := funcs.Filter(func(i int, v *aq.FuncDecl) bool {
		return v.Recv().Type().UnwrapPtr().Name() == typ &&
			v.Params().Count() == 0 &&
			v.Results().Count() == 1
	})

	listGetters := funcs.Filter(func(i int, v *aq.FuncDecl) bool {
		return v.Recv().Type().UnwrapPtr().Name() == listName
	})

	list := []listGetterParams{}
	for _, f := range getters {
		name := f.Name()
		if lgnoreGetters[name] {
			continue
		}

		nameUpper := goname.UpperCamelCase(name)
		listFuncName := inflection.Plural(nameUpper)
		if listFuncName == nameUpper {
			listFuncName = nameUpper + "List"
		}

		if listGetters.Has(func(i int, v *aq.FuncDecl) bool { return v.Name() == listFuncName }) {
			continue
		}

		typ := "[]" + f.Results().First().Type().GoCode()

		list = append(list, listGetterParams{
			Name:         f.Name(),
			ListType:     typ,
			ListFuncName: listFuncName,
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

var lgnoreGetters = map[string]bool{
	"Exists": true,
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

func (r {{ .ListName }}) Compact() {{ .ListName }} {
	l := {{ .ListName }}{}
	for _, v := range r {
		if v == nil {
			l = append(l, v)
		}
	}
	return l
}

func (r {{ .ListName }}) Concat(l {{ .ListName }}) {{ .ListName }} {
	return append(append({{ .ListName }}{}, r...), l...)
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

func (r {{ .ListName }}) Has(f func(i int, v *{{ .Type }}) bool) bool {
	return r.Some(f)
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

{{- $listName := .ListName }}
{{- $typeUpper := .TypeUpper }}
{{- $type := .Type }}
{{- range .Fields }}

func (r {{ $listName }}) {{ .NameUpperPlural }}() []{{ .Type }} {
 	l := make([]{{ .Type }}, 0, len(r))
 	for _, r := range r {
 		l = append(l, r.{{ .Name }})
 	}
 	return l
}
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
{{- range .Getters }}

func (r {{ $listName }}) {{ .ListFuncName }}() {{ .ListType }} {
 	l := make({{ .ListType }}, 0, len(r))
 	for _, r := range r {
 		l = append(l, r.{{ .Name }}())
 	}
 	return l
}
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
