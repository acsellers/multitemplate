package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	wd, _ := os.Getwd()
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("Could not open file:", path)
				return nil
			}

			fileName := filepath.Base(path)
			nameOnly := fileName[:len(fileName)-len(filepath.Ext(fileName))]
			s := Spec{Title: strings.ToTitle(nameOnly)}
			err = json.Unmarshal(buf, &s)
			if err != nil {
				fmt.Println("Could not unmarshal file:", err)
				return nil
			}
			for i, t := range s.Tests {
				b, _ := json.Marshal(t.Data)
				t.Marshalled = "`" + string(b) + "`"
				t.Template = "`" + t.Template + "`"
				t.Expected = "`" + t.Expected + "`"
				for k, v := range t.Partials {
					t.Partials[k] = "`" + v + "`"
				}
				s.Tests[i] = t
			}

			f, e := os.Create(filepath.Join(wd, "gen", nameOnly+"_test.go"))
			if e == nil {
				tmpl.Execute(f, s)
			}
			f.Close()
		} else {
			fmt.Printf("Not a json file:%s\n", path)
		}
		return nil
	})
}

type Spec struct {
	Title    string
	Overview string `json:"overview"`
	Tests    []Test `json:"tests"`
}

type Test struct {
	Name        string                 `json:"name"`
	Data        map[string]interface{} `json:"data"`
	Expected    string                 `json:"expected"`
	Template    string                 `json:"template"`
	Description string                 `json:"desc"`
	Partials    map[string]string
	Marshalled  string
}

var testTemplate = `
package mussed

import (
  "bytes"
  "encoding/json"
  "testing"
  "html/template"
)

/*
{{.Overview}}
*/

{{ $title := .Title }}
{{ range $index, $test := .Tests }}

func Test{{$title}}{{$index}}(t *testing.T) {
  // {{$test.Name}}

  within(t, func(test *aTest) {
		t := template.New("test").Funcs(RequiredFuncs)
		trees, err := Parse("test.mustache",{{ $test.Template }})
		test.IsNil(err)
    for name, tree := range trees {
      t, err = t.AddParseTree(name, tree)
      test.IsNil(err)
    }
    {{ range $name, $content := $test.Partials }}
    trees, err = Parse("{{ $name }}.mustache",{{ $content }})
    test.IsNil(err)
    for name, tree := range trees {
      t, err = t.AddParseTree(name, tree)
      test.IsNil(err)
    }
    {{end}}

    data := make(map[string]interface{})
    test.IsNil(json.Unmarshal([]byte({{$test.Marshalled}}), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual({{$test.Expected}}, b.String())
  })
}
{{ end }}
`

var tmpl = template.Must(template.New("thing").Parse(testTemplate))
