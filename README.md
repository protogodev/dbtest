# dbtest

Database testing made easy in Go.


## Features

1. Human-friendly
  
    Define the minimum test specification in a YAML-based DSL, then all tests can be generated automatically.

2. Customizable

    You can specify your own Go template when generating tests.

3. Opinionated
  
    - [The database is a detail][1], which means the database-related code always implements a Go interface.
    - [Run database tests against a real database][2], instead of relying on mocks.
        + The limitation is that running tests in parallel is not supported.
    - Use the same [test fixture][3] for all subtests of the same [method under test][4].
    - Only [test the public methods][5], thus the generated tests are in a separate `_test` package.
    - Prefer [table driven tests][6].


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

Check out the [Godoc][7].


## License

[MIT](LICENSE)


[1]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[2]: https://github.com/go-testfixtures/testfixtures
[3]: https://en.wikipedia.org/wiki/Test_fixture#Software
[4]: http://xunitpatterns.com/SUT.html
[5]: https://martinfowler.com/articles/practical-test-pyramid.html#WhatToTest
[6]: https://github.com/golang/go/wiki/TableDrivenTests
[7]: https://pkg.go.dev/github.com/RussellLuo/dbtest