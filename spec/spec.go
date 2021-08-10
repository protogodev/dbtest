package spec

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Import struct {
	Alias string `yaml:"alias"`
	Path  string `yaml:"path"`
}

func (i Import) ImportString() string {
	s := fmt.Sprintf("%q", i.Path)
	if i.Alias != "" {
		s = i.Alias + " " + s
	}
	return s
}

type Row map[string]interface{}

func (r Row) LiteralString() string {
	s := fmt.Sprintf("%#v", r)
	return s[len("spec.Row"):]
}

type Rows []Row

func (r Rows) Equal(other Rows) bool {
	if len(r) == 0 && len(other) == 0 {
		return true
	}
	return fmt.Sprintf("%#v", r) == fmt.Sprintf("%#v", other)
}

type DataAssertion struct {
	Query  string `yaml:"query"`
	Result Rows   `yaml:"result"`
}

type Subtest struct {
	Name     string                 `yaml:"name"`
	In       map[string]interface{} `yaml:"in"`
	WantOut  map[string]interface{} `yaml:"wantOut"`
	WantData []DataAssertion        `yaml:"wantData"`
}

type Test struct {
	Name     string          `yaml:"name"`
	Fixture  map[string]Rows `yaml:"fixture"`
	Subtests []Subtest       `yaml:"subtests"`
}

type Spec struct {
	RawImports []string `yaml:"imports"`
	Imports    []Import `yaml:"-"`
	Testee     string   `yaml:"testee"`
	Tests      []Test   `yaml:"tests"`
}

func New(testFilename string) (*Spec, error) {
	b, err := ioutil.ReadFile(testFilename)
	if err != nil {
		return nil, err
	}

	spec := &Spec{}
	err = yaml.Unmarshal(b, spec)
	if err != nil {
		return nil, err
	}

	imports, err := getImports(spec.RawImports)
	if err != nil {
		return nil, err
	}
	spec.Imports = append(spec.Imports, imports...)

	return spec, nil
}

func getImports(rawImports []string) (imports []Import, err error) {
	var path, alias string

	for i, str := range rawImports {
		fields := strings.Fields(str)
		switch len(fields) {
		case 1:
			alias, path = "", fields[0]
		case 2:
			alias, path = fields[0], fields[1]
		default:
			return nil, fmt.Errorf("invalid path in imports[%d]: %s", i, str)
		}

		imports = append(imports, Import{
			Path:  path,
			Alias: alias,
		})
	}

	return
}
