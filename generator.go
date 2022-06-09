package dbtest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/protogodev/dbtest/builtin"
	"github.com/protogodev/dbtest/spec"
	protogocmd "github.com/protogodev/protogo/cmd"
	"github.com/protogodev/protogo/generator"
	"github.com/protogodev/protogo/parser"
	"github.com/protogodev/protogo/parser/ifacetool"
)

func init() {
	protogocmd.MustRegister(&protogocmd.Plugin{
		Name: "dbtest",
		Cmd:  protogocmd.NewGen(&Generator{}),
	})
}

type Generator struct {
	OutFileName      string `name:"out" help:"output filename (default \"./<srcPkgName>_test.go\")"`
	Formatted        bool   `name:"fmt" default:"true" help:"whether to make the test code formatted"`
	TestSpecFileName string `name:"spec" default:"./dbtest.spec.yaml" help:"the test specification in YAML"`
	TemplateFileName string `name:"tmpl" help:"the template to render (default to builtin template)"`
}

func (g *Generator) Generate(data *ifacetool.Data) (*generator.File, error) {
	if g.OutFileName == "" {
		g.OutFileName = fmt.Sprintf("./%s_test.go", data.SrcPkgName)
	}

	testSpec, err := spec.New(g.TestSpecFileName)
	if err != nil {
		return nil, err
	}

	imports := testSpec.Imports
	for _, i := range data.Imports {
		imports = append(imports, spec.Import{Path: i.Path, Alias: i.Alias})
	}

	tmplData := struct {
		DstPkgName    string
		SrcPkgName    string
		InterfaceName string
		Imports       []spec.Import
		Testee        string
		Tests         []spec.Test
	}{
		DstPkgName:    parser.PkgNameFromDir(filepath.Dir(g.OutFileName)),
		SrcPkgName:    data.SrcPkgName,
		InterfaceName: data.InterfaceName,
		Imports:       imports,
		Testee:        testSpec.Testee,
		Tests:         testSpec.Tests,
	}

	methodMap := make(map[string]*ifacetool.Method)
	for _, method := range data.Methods {
		methodMap[method.Name] = method
	}

	template, err := getTemplate(g.TemplateFileName)
	if err != nil {
		return nil, err
	}

	return generator.Generate(template, tmplData, generator.Options{
		Funcs: map[string]interface{}{
			"title": strings.Title,
			"fmtArgCSV": func(csv string, format string) string {
				if csv == "" {
					return ""
				}

				sep := ", "
				args := strings.Split(csv, sep)

				var results []string
				for _, a := range args {
					r := strings.NewReplacer("$Name", a, ">Name", strings.Title(a))
					results = append(results, r.Replace(format))
				}

				return strings.Join(results, sep)
			},
			"interfaceMethod": func(name string) *ifacetool.Method {
				method, ok := methodMap[name]
				if !ok {
					return nil
				}
				return method
			},
			"goString": func(m map[string]interface{}) string {
				return fmt.Sprintf("%#v", m)
			},
		},
		Formatted:      g.Formatted,
		TargetFileName: g.OutFileName,
	})
}

func getTemplate(fileName string) (string, error) {
	if fileName == "" {
		return builtin.Template, nil
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
