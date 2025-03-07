# httpmock

[![Go Report Card](https://goreportcard.com/badge/github.com/oupo1337/httpmock)](https://goreportcard.com/report/github.com/oupo1337/httpmock)

HTTPMock is a library to easily mock your http clients and describe their behavior.

It's designed to be as simple as possible, first you describe the client behavior:
 - How many calls should it make? On which path?
 - What HTTP status codes should it return?
 - Should it receive headers? Query parameters? A body?

You can now use the mock client as a regular HTTP client in your code.

Finally, you check that all expected calls were done with the mock client by the tested code.

## Examples

### Basic

```go
func Test_basic(t *testing.T) {
    // We declare a mock variable and its behavior.
    // It's an HTTP client waiting for a GET request on /path.
    // It will return a 200 status code.
    mock := httpmock.New(t).
        WithRequest(http.MethodGet, "/path",
            httpmock.ReturnStatus(http.StatusOK),
        )

    // Do something with mock variable (It's a regular http.Client object).
    doSomething(mock)

    // We check that all expected requests were done.
    mock.AssertExpectations()
}

// Or you can declare your mock in a different way without using option functions
// similar way to the testify/mock package.
func Test_basic_bis(t *testing.T) {
	mock := httpmock.New(t)
	mock.On(http.MethodGet, "/path").ReturnStatus(http.StatusOK)

	doSomething(mock)
	mock.AssertExpectations()
}
```

### Advanced

```go
func Test_advanced(t *testing.T) {
    // We declare a mock variable and its behavior.
    // It's an HTTP client waiting for a POST request on /form.
    // It will return a 201 status code.
    // It expects a body, an authorization header and will return a body.
    //
    // It's also waiting for a DELETE request on /route
    // returning a 204 status code.
    // It expects a query param.
    mock := httpmock.New(t).
        WithRequest(http.MethodPost, "/form",
            httpmock.ExpectBody(`{"some": "data"}`),
            httpmock.ExpectHeader("Authorization", []string{"Bearer token"}),
            httpmock.ReturnStatus(http.StatusCreated),
            httpmock.ReturnBody(`{"a": "response"}`),
        ).
        WithRequest(http.MethodDelete, "/route",
            httpmock.ExpectQueryParam("param", "value"),
            httpmock.ReturnStatus(http.StatusNoContent),
        )

    // Do something with the mock variable (It's a regular http.Client object).
    doSomething(mock)

    // We check that all expected requests were made.
    mock.AssertExpectations()
}

// Or you can declare your mock in a different way without using option functions
// similar way to the testify/mock package.
func Test_advanced_bis(t *testing.T) {
	mock := httpmock.New(t)
	mock.On(http.MethodPost, "/form").
		ExpectBody(`{"some": "data"}`).
		ExpectHeader("Authorization", []string{"Bearer token"}).
		ReturnStatus(http.StatusCreated).
		ReturnBody(`{"a": "response"}`)

	mock.On(http.MethodDelete, "/route").
		ExpectQueryParam("param", "value").
		ReturnStatus(http.StatusNoContent)

	doSomething(mock)
	mock.AssertExpectations()
}
```

### More examples

See example file [here](examples/example_test.go)

## Documentation

### Options functions

| Name                   | Description                                                                                      | Type             |
|------------------------|--------------------------------------------------------------------------------------------------|------------------|
| ReturnStatus           | Sets the http status code returned by the request.                                               | int              |
| ReturnBodyRaw          | Sets the body returned by the request.                                                           | string           |
| ReturnBodyFromObject   | Sets the body returned by the request from an object. (Using json.Marshal function)              | interface{}      |
| ReturnError            | Sets an error returned by the http client.                                                       | error            |
| ExpectBody             | Will expect a body in the received request and asserts that strings are equal.                   | string           |
| ExpectJSON             | Will expect a body in the received request and asserts that the JSONs are equal.                 | string           |
| ExpectHeader           | Will expect a header in the received request and asserts that the name and value are equal.      | string, string   |
| ExpectQueryParamValues | Will expect a query param in the received request and assert that the name and values are equal. | string, []string |
| ExpectQueryParam       | Will expect a query param in the received request and assert that the name and values are equal. | string, string   |

