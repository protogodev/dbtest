# User Store

This example illustrates how to apply declarative testing for a typical database application.


## Prerequisites

1. Prepare the database

   Execute the SQL [user.sql](user.sql):
   
    ```bash
    $ mysql -uroot -p -e 'source ./user.sql'
    ```

2. Implement the testee factory

   See [NewTestee()](store.go#L117-L133).

3. Write the test specification

    See [dbtest.spec.yaml](dbtest.spec.yaml).


## Generate tests

```bash
$ go generate
```

See [store_test.go](store_test.go) for generated tests.


## Run tests

```bash
$ go test -v -race
```

<details>
  <summary> Result </summary>

```bash
=== RUN   TestCreateUser
=== RUN   TestCreateUser/new_user
=== RUN   TestCreateUser/duplicate_user
--- PASS: TestCreateUser (0.05s)
    --- PASS: TestCreateUser/new_user (0.01s)
    --- PASS: TestCreateUser/duplicate_user (0.02s)
=== RUN   TestGetUser
=== RUN   TestGetUser/ok
=== RUN   TestGetUser/not_found
--- PASS: TestGetUser (0.02s)
    --- PASS: TestGetUser/ok (0.00s)
    --- PASS: TestGetUser/not_found (0.01s)
=== RUN   TestUpdateUser
=== RUN   TestUpdateUser/ok
=== RUN   TestUpdateUser/not_found
--- PASS: TestUpdateUser (0.04s)
    --- PASS: TestUpdateUser/ok (0.02s)
    --- PASS: TestUpdateUser/not_found (0.02s)
=== RUN   TestDeleteUser
=== RUN   TestDeleteUser/ok
=== RUN   TestDeleteUser/not_found
--- PASS: TestDeleteUser (0.04s)
    --- PASS: TestDeleteUser/ok (0.01s)
    --- PASS: TestDeleteUser/not_found (0.02s)
PASS
ok      github.com/RussellLuo/dbtest/examples   0.175s
```

</details>
