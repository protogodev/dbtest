package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/pkgtool"

	"github.com/protogodev/dbtest/builtin"
	"github.com/protogodev/dbtest/spec"
)

type Options struct {
	OutFileName      string
	Formatted        bool
	TestSpecFileName string
	TemplateFileName string
}

type userFlags struct {
	Options

	args []string
}

func main() {
	var flags userFlags
	flag.StringVar(&flags.TestSpecFileName, "spec", "./dbtest.spec.yaml", "the test specification in YAML")
	flag.StringVar(&flags.TemplateFileName, "tmpl", "", "the template to render (default to builtin template)")
	flag.StringVar(&flags.OutFileName, "out", "", `output filename (default "./<srcPkgName>_test.go")`)
	flag.BoolVar(&flags.Formatted, "fmt", true, "whether to make the test code formatted")

	flag.Usage = func() {
		fmt.Println(`dbtest [flags] source-file interface-name`)
		flag.PrintDefaults()
	}

	flag.Parse()
	flags.args = flag.Args()

	if err := run(flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}
}

func run(flags userFlags) error {
	if len(flags.args) != 2 {
		return errors.New("need 2 arguments")
	}

	srcFilename, interfaceName := flags.args[0], flags.args[1]

	srcFilename, err := filepath.Abs(srcFilename)
	if err != nil {
		return err
	}

	// Non-empty pkgName makes all type names used in the interface full-qualified.
	pkgName := "x"
	data, err := pkgtool.ParseInterface(pkgName, srcFilename, interfaceName)
	if err != nil {
		return err
	}

	if flags.Options.OutFileName == "" {
		flags.Options.OutFileName = fmt.Sprintf("./%s_test.go", data.SrcPkgName)
	}

	file, err := generate(&flags.Options, data)
	if err != nil {
		return err
	}

	if err := file.Write(); err != nil {
		return err
	}

	return nil
}

func generate(opts *Options, ifaceData *ifacetool.Data) (*generator.File, error) {
	testSpec, err := spec.New(opts.TestSpecFileName)
	if err != nil {
		return nil, err
	}

	imports := testSpec.Imports
	for _, i := range ifaceData.Imports {
		imports = append(imports, spec.Import{Path: i.Path, Alias: i.Alias})
	}

	data := struct {
		DstPkgName    string
		SrcPkgName    string
		InterfaceName string
		Imports       []spec.Import
		Testee        string
		Tests         []spec.Test
	}{
		DstPkgName:    pkgtool.PkgNameFromDir(filepath.Dir(opts.OutFileName)),
		SrcPkgName:    ifaceData.SrcPkgName,
		InterfaceName: ifaceData.InterfaceName,
		Imports:       imports,
		Testee:        testSpec.Testee,
		Tests:         testSpec.Tests,
	}

	methodMap := make(map[string]*ifacetool.Method)
	for _, method := range ifaceData.Methods {
		methodMap[method.Name] = method
	}

	template, err := getTemplate(opts.TemplateFileName)
	if err != nil {
		return nil, err
	}

	return generator.Generate(template, data, generator.Options{
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
		Formatted:      opts.Formatted,
		TargetFileName: opts.OutFileName,
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
