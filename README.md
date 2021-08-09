go-mocking
==========

This simple package is intended to provide an example of how to mock dependencies when
unit testing in Go. The package provides a simple REST API that serves up data on your
friends. The data is stored in a file called *friends.go* alongside *main.go*. The API
accepts GET and POST methods. To add new data, pass a JSON object with "name", "age",
and "occupation" fields. For example,

```
$ curl -X POST -d '{"name":Jeff,"age":25,"occupation":"Programmer"}' 'localhost:8080'
```
This call will add a new entry with name "Jeff". If an entry already exists with the same
name, an error is returned in the HTTP response. The data is automatically saved to the
data file and will persist across server instances.

To run the unit tests, run `go test` inside the *cmd* directory. All of the tests should
pass. The *mock.go* file contains the mocked datastore implementation that allows the
unit tests to run without interacting with the file system.

Please reach out to jeff.moorhead1@gmail.com with any questions or concerns about this
package.
