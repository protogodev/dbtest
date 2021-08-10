# dbtest

Database testing made easy in Go.


## Installation

```bash
$ go get -u github.com/RussellLuo/dbtest
```

<details open>
  <summary> Usage </summary>

```bash
$ dbtest -h
dbtest [flags] source-file interface-name
  -fmt
        whether to make the test code formatted (default true)
  -out string
        output filename (default "./<srcPkgName>_test.go")
  -spec string
        the test specification in YAML (default "./dbtest.spec.yaml")
  -tmpl string
        the template to render (default to builtin template)
```

</details>


## Examples

See [examples](examples).


## Documentation

Check out the [Godoc][1].


## License

[MIT](LICENSE)


[1]: https://pkg.go.dev/github.com/RussellLuo/dbtest