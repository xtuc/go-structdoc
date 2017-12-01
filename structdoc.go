package structdoc

import (
	"html/template"
	"os"
	"reflect"
	"regexp"
)

var (
	jsonTagRegexp = regexp.MustCompile(`json:"(.*)"`)

	docTemplate = `
<a name="{{.Title}}"></a>
### ` + "`" + `{{.Title}}` + "`" + `
{{ if .Fields }}
| Property name | Type | Representation in the JSON |
|---|---|
{{ range $value := .Fields }}| {{ $value.Name }} | ` + "`" + `{{ $value.Type }}` + "`" + ` | ` + "`" + `{{ $value.TypeInJson }}` + "`" + ` |
{{ end }}{{ end }}
`
)

type DocField struct {
	Name       string
	Type       string
	TypeInJson string
}

type DocEntry struct {
	Title  string
	Fields []DocField
}

type Generator struct {
	normalizeTypeName func(string) string
	runtimeTypes      map[string]string
}

func MakeGenerator(normalizeTypeName func(string) string, runtimeTypes map[string]string) Generator {
	return Generator{
		normalizeTypeName: normalizeTypeName,
		runtimeTypes:      runtimeTypes,
	}
}

func (g Generator) GeneratorFor(_struct interface{}) {
	structType := reflect.TypeOf(_struct)

	entry := DocEntry{
		Title:  structType.Name(),
		Fields: make([]DocField, 0),
	}

	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		typeName := g.normalizeTypeName(structField.Type.String())

		docField := DocField{
			Name:       structField.Name,
			Type:       typeName,
			TypeInJson: g.runtimeTypes[typeName],
		}

		if docField.Type == "" {
			docField.Type = "unknown"
		}

		if structField.Tag != "" {
			matches := jsonTagRegexp.FindSubmatch([]byte(structField.Tag))
			docField.Name = string(matches[1])
		}

		entry.Fields = append(entry.Fields, docField)
	}

	tmpl, err := template.New("").Parse(docTemplate)

	err = tmpl.Execute(os.Stdout, entry)

	if err != nil {
		panic(err)
	}
}
