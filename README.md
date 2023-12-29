[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE.txt)
[![Webex Events](https://github.com/SocioEvents/webex-events-go-sdk/actions/workflows/go-test.yml/badge.svg)](https://github.com/SocioEvents/webex-events-go-sdk/actions)

⚠️ This library has not been released yet.
# Webex Events Api Go SDK

Webex Events provides a range of additional SDKs to accelerate your development process.
They allow a standardized way for developers to interact with and leverage the features and functionalities.
Pre-built code modules will help access the APIs with your private keys, simplifying data gathering and update flows.

Requirements
-----------------

Go 1.21+

Installation
-----------------

`go get github.com/SocioEvents/webex-events-go-sdk`

Configuration
-----------------

```go
	var ctx = context.Background()
	config := NewConfig()
	config.SetAccessToken("sk_live_your_access_token")
	config.SetMaxRetries(5) // Default is 5
	config.SetTimeout(time.Duration(30) * time.Second) // default is 30 seconds
	var loggerHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	config.SetLoggerHandler(loggerHandler) // default is Error loglevel to stdout
```

Usage
-----------------
```go
    var variables = make(map[string]any)
	var query = "query CurrenciesList{ currenciesList{ isoCode }}"
	var operationName = "CurrenciesList"

	queryRequest := QueryRequest{
		Query:         query,
		OperationName: operationName,
		Variables:     variables,
	}
	var response, err = client.Query(ctx, &queryRequest)
	fmt.Println(response, err)
```

For Introspection
-----------------
```go
	var response, err = client.DoIntrospectionQuery(ctx)
	fmt.Println(response, err)
```

Idempotency
-----------------
The API supports idempotency for safely retrying requests without accidentally performing the same operation twice.
When doing a mutation request, use an idempotency key. If a connection error occurs, you can repeat
the request without risk of creating a second object or performing the update twice.

To perform mutation request, you must add a header which contains the idempotency key such as
`Idempotency-Key: <your key>`. The SDK does not produce an Idempotency Key on behalf of you if it is missed.
Here is an example like the following:

```go
var variables = map[string]any{
    "input": map[string]any{
    "ids":     []int{1, 2, 3},
    "eventId": 1,
    },
}
var query = `
          mutation TrackDelete($input: TrackDeleteInput!) {
            trackDelete(input: $input) {
              success
            }
          }
`
var operationName = "TrackDelete"

queryRequest := QueryRequest{
	Query:         query,
    OperationName: operationName,
    Variables:     variables,
    IdempotencyKey: "61672155-56d3-4375-a864-52e7bba4f445",
}
var response, err = client.Query(ctx, &queryRequest)
fmt.Println(response, err)

```

Telemetry Data Collection
-----------------
Webex Events collects telemetry data, including hostname, operating system, language and SDK version, via API requests.
This information allows us to improve our services and track any usage-related faults/issues. We handle all data with
the utmost respect for your privacy. For more details, please refer to the Privacy Policy at https://www.cisco.com/c/en/us/about/legal/privacy-full.html

Development
-----------------

After checking out the repo, install dependencies. Then, run `go test` to run the tests.

Contributing
-----------------
Please see the [contributing guidelines](CONTRIBUTING.md).

License
-----------------

The library is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

Code of Conduct
-----------------

Everyone interacting in the Webex Events API project's codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](https://github.com/SocioEvents/webex-events-go-sdk/blob/main/CODE_OF_CONDUCT.md).
