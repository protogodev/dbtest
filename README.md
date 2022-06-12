# dbtest

Database testing made easy in Go.


## Features

1. Declarative
  
    Define the minimum test specification in a YAML-based DSL, then all tests can be generated automatically.

2. Customizable

    You can specify your own Go template when generating tests.

3. Opinionated
  
    - [The database is a detail][1], which means the database-related code should always implement a Go interface.
    - [Run database tests against a real database][2], instead of relying on mocks.
        + The limitation is that running tests in parallel is not supported.
    - Use the same [test fixture][3] for all subtests of the same [method under test][4].
    - Only [test the public methods][5], thus the generated tests are in a separate `_test` package.
    - Leverage [TestMain][6] to do global setup and teardown.
    - Prefer [table driven tests][7].


## Installation

Make a custom build of [protogo](https://github.com/protogodev/protogo):

```bash
$ protogo build --plugin=github.com/protogodev/dbtest
```

Or build from a local fork:

```bash
$ protogo build --plugin=github.com/protogodev/dbtest=../my-fork
```

<details open>
  <summary> Usage </summary>

```bash
$ protogo dbtest -h
Usage: protogo dbtest <source-file> <interface-name>

Arguments:
  <source-file>       source-file
  <interface-name>    interface-name

Flags:
  -h, --help                         Show context-sensitive help.

      --out=STRING                   output filename (default "./<srcPkgName>_test.go")
      --fmt                          whether to make the test code formatted
      --spec="./dbtest.spec.yaml"    the test specification in YAML
      --tmpl=STRING                  the template to render (default to builtin template)
```

</details>


## Examples

See [examples](examples).


## Documentation

Check out the [Godoc][8].


## License

[MIT](LICENSE)


[1]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[2]: https://github.com/go-testfixtures/testfixtures
[3]: https://en.wikipedia.org/wiki/Test_fixture#Software
[4]: http://xunitpatterns.com/SUT.html
[5]: https://martinfowler.com/articles/practical-test-pyramid.html#WhatToTest
[6]: https://pkg.go.dev/testing#hdr-Main
[7]: https://github.com/golang/go/wiki/TableDrivenTests
[8]: https://pkg.go.dev/github.com/protogodev/dbtest
