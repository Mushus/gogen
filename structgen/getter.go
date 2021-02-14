package main

import (
	"bytes"
	"text/template"

	"github.com/Mushus/gogen/aq"
	"github.com/Mushus/gogen/goname"
	"github.com/pkg/errors"
)

type getterGenerator struct {
	structName string
	params     []getterParams
}

type getterParams struct {
	Receiver  string
	Name      string
	FieldName string
	Type      string
}

func newGetterGenerator(structName string) *getterGenerator {
	return &getterGenerator{
		structName: structName,
	}
}

func (g *getterGenerator) collectParams(aqi *aq.Instance, oldGeneratedCode []byte) error {
	s := aqi.Structs().FindOne(aq.StructNameIs(g.structName))
	if !s.Exists() {
		return errors.Errorf("struct %#v not found", g.structName)
	}

	g.params = g.createGetterParams(s.Fields())
	return nil
}

func (g *getterGenerator) createGetterParams(fields aq.Fields) []getterParams {
	receiver := g.structName

	getters := []getterParams{}

	for _, f := range fields {
		if f.IsPublic() {
			continue
		}

		tag := newTagValue(f.Tag().Body().Get("getter"))

		if tag.contains("-") {
			continue
		}

		fieldName := f.Name()
		name := tag.find(func(v string) bool { return v != "-" })
		if name == "" {
			name = goname.UpperCamelCase(fieldName)
		}

		typ := f.Type().GoCode()

		getters = append(getters, getterParams{
			Receiver:  receiver,
			Name:      name,
			FieldName: fieldName,
			Type:      typ,
		})
	}

	return getters
}

const getterTmplStr = `
{{- range . }}
func (r *{{ .Receiver }}) {{ .Name }}() {{ .Type }} {
	if r == nil {
		var v {{ .Type }}
		return v
	}
	return r.{{ .FieldName }}
}

{{ end -}}
`

var getterTmpl = template.Must(template.New("getter").Parse(getterTmplStr))

func (g *getterGenerator) generate() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	br := new(bytes.Buffer)
	if err := getterTmpl.Execute(br, g.params); err != nil {
		return nil, errors.Wrap(err, "cannot execute template")
	}
	return br.Bytes(), nil
}
