package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Mushus/gogen/aq"
	"github.com/Mushus/gogen/goname"
	"github.com/pkg/errors"
)

type constructGenerator struct {
	structName string
	params     constructParams
}

type constructParams struct {
	ConstructorName    string
	StructName         string
	Fields             []constructField
	ConstructTemplate  string
	HasConstructMethod bool
}

func newConstructGenerator(structName string) *constructGenerator {
	return &constructGenerator{
		structName: structName,
	}
}

func (g *constructGenerator) collectParams(aqi *aq.AQ, oldGeneratedCode []byte) error {
	if g == nil {
		return nil
	}

	{
		s := aqi.Structs().Find(aq.StructNameIs(g.structName))
		if !s.Exists() {
			return errors.Errorf("struct %#v not found", g.structName)
		}

		g.params.StructName = g.structName

		UpperStructName := goname.UpperCamelCase(g.structName)
		g.params.ConstructorName = fmt.Sprintf("New%s", UpperStructName)

		g.params.Fields = createConstructField(s.Fields())
	}

	{
		constructLogic := defaultConstructLogic
		if b := getMarkedCode(oldGeneratedCode, constructLogicStartMark, constructLogicEndMark); b != nil {
			constructLogic = string(b)
		}

		g.params.ConstructTemplate = constructLogicStartMark + constructLogic + constructLogicEndMark
	}

	{
		g.params.HasConstructMethod = aqi.
			Funcs().
			Some(func(i int, v *aq.FuncDecl) bool {
				return v.Name() == "construct" &&
					v.Recv().Type().Name() == g.structName &&
					v.Recv().Type().IsPtr()
			})
	}

	return nil
}

const constructTmplStr = `
func {{ .ConstructorName }}(
{{- range .Fields }}
	{{ .Name }} {{ .Type }},
{{- end }}
) *{{ .StructName }} {
	c := &{{ .StructName }} {
{{- range .Fields }}
		{{ .Name }}: {{ .Name }},
{{- end }}
	}

{{- if .HasConstructMethod }}
	c.construct()
{{- end }}

	return c
}
`

var constructTmpl = template.Must(template.New("construct").Parse(constructTmplStr))

const constructLogicStartMark = "// ===== construct ====="
const constructLogicEndMark = "// ===== construct ====="

const defaultConstructLogic = `
	
// write constructor logic here

return c
`

func (g *constructGenerator) generate() ([]byte, error) {
	if g == nil {
		return nil, nil
	}

	br := new(bytes.Buffer)
	if err := constructTmpl.Execute(br, g.params); err != nil {
		return nil, errors.Wrap(err, "cannot execute template")
	}
	return br.Bytes(), nil
}

type constructField struct {
	Name string
	Type string
}

func createConstructField(fields aq.Fields) []constructField {
	r := []constructField{}
	for _, f := range fields {
		tag := newTagValue(f.Tag().Body().Get("construct"))
		if tag.contains("-") {
			continue
		}

		r = append(r, constructField{
			Name: f.Name(),
			Type: f.Type().GoCode(),
		})
	}

	return r
}
