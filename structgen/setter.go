package main

import (
	"bytes"
	"text/template"

	"github.com/Mushus/gogen/aq"
	"github.com/Mushus/gogen/goname"
	"github.com/pkg/errors"
)

type setterGenerator struct {
	structName string
	params     []setterParams
}

type setterParams struct {
	Receiver  string
	Name      string
	FieldName string
	Type      string
}

func newSetterGenerator(structName string) *setterGenerator {
	return &setterGenerator{
		structName: structName,
	}
}

func (g *setterGenerator) collectParams(aqi *aq.Instance, oldGeneratedCode []byte) error {
	s := aqi.Structs().Find(aq.StructNameIs(g.structName))
	if !s.Exists() {
		return errors.Errorf("struct %#v not found", g.structName)
	}

	g.params = g.createSetterParams(s.Fields())
	return nil
}

func (g *setterGenerator) createSetterParams(fields aq.Fields) []setterParams {
	receiver := g.structName

	setter := []setterParams{}

	for _, f := range fields {
		if f.IsExported() {
			continue
		}

		tag := newTagValue(f.Tag().Body().Get("setter"))

		if tag.contains("-") {
			continue
		}

		fieldName := f.Name()
		name := tag.find(func(v string) bool { return v != "-" })
		if name == "" {
			name = "Set" + goname.UpperCamelCase(fieldName)
		}

		typ := f.Type().GoCode()

		setter = append(setter, setterParams{
			Receiver:  receiver,
			Name:      name,
			FieldName: fieldName,
			Type:      typ,
		})
	}

	return setter
}

const setterTmplStr = `
{{- range . }}
func (r *{{ .Receiver }}) {{ .Name }}({{ .FieldName }} {{ .Type }})  {
	r.{{ .FieldName }} = {{ .FieldName }}
}

{{ end -}}
`

var setterTmpl = template.Must(template.New("setter").Parse(setterTmplStr))

func (g *setterGenerator) generate() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	br := new(bytes.Buffer)
	if err := setterTmpl.Execute(br, g.params); err != nil {
		return nil, errors.Wrap(err, "cannot execute template")
	}
	return br.Bytes(), nil
}
