[![Go Report Card](https://goreportcard.com/badge/myshkin5/hasher)](https://goreportcard.com/report/myshkin5/hasher)

# hasher

Hasher hashes passwords posted to one endpoint and then returns the encoded hashes on another endpoint after a 5 second delay.

## Running the tests

To run the tests in a console, change to the directory of this readme file then execute the following:

```bash
go test -race -bench=. ./...
```

## Running the service

### Running the service locally

To run the service locally in a console, change to the directory of this readme file then execute the following:

```bash
go run api/main.go
```

### Running from an executable

Productions servers typically only get an executable to run. To build an executable, first set the `GOOS` and `GOARCH` environment variables from the [list of supported values](https://github.com/golang/go/blob/master/src/go/build/syslist.go). Then execute the following:

```bash
go build -o hasher-api github.com/myshkin5/hasher/api
```

### Environment variables

The following environment variables can be used to control how hasher runs:

Environment variable | Default value | Description
--- | --- | ---
`HASH_STORE_COUNT` | `10000` | The maximum number of most recently created hashes to keep in memory at any given time. 
`LOG_LEVEL` | `info` | Determines what logs are emitted at runtime. Acceptable values are: `info`, `warn`, `error`, and `panic`
`PORT` | `8080` | The port the service will accept request on.
`SERVER_ADDR` | `localhost` | The server address the service will accept requests on. Set to `0.0.0.0` to accept requests from the network.

## API Documentation

Endpoint |
--- |
[`GET /`](#get-) |
[`GET /health`](#get-health) |
[`POST /hash`](#post-hash) |
[`GET /hash/:requestId`](#get-hashrequestid) |
[`GET /stats`](#get-stats) |
[`[ANY] /shutdown`](#any-shutdown) |

### `GET /`

#### Response Statuses

`200 - OK`: Only `200 - OK` is returned.

#### OK Response Body

Field | Description
--- | ---
`health` | Link to the health endpoint.

##### Example

```json
{
    "health": "http://localhost:8080/health"
}
```

### `GET /health`

#### Response Statuses

`200 - OK`: Currently only `200 - OK` is returned.

#### OK Response Body

Field | Description
--- | ---
`health` | The health of the service. Currently only `GOOD` is supported.

##### Example

```json
{
    "health": "GOOD"
}
```

### `POST /hash`

#### Request Body

Field | Description
--- | ---
`password` | The password to be hashed.

##### Example

```http request
password=angryMonkey
```

#### Response Statuses

`201 - Created`: The hash request was successfully created.

`400 - Bad Request`: The request was malformed and could not be processed.

#### Created Response Body

Returns the `requestId`. See [`GET /hash/:requestId`](#get-hashrequestid) for the definition of a request id.

##### Example

```
42
```

### `GET /hash/:requestId`

#### Request Parameters

Field | Description
--- | ---
`requestId` | The identifier of a previously requested hash. See [`POST /hash`](#post-hash) to request a hash.

#### Response Statuses

`200 - OK`: Returned on success.

`400 - Bad Request`: The request was malformed and could not be processed.

`404 - Not Found`: The requested hash could not be found possibly because it isn't available yet.

`500 - Internal Server Error`: Returned when there is an internal server error.

#### OK Response Body

The hash is returned in the body.

##### Example

```
ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
```

### `GET /stats`

#### Response Statuses

`200 - OK`: Currently only `200 - OK` is returned.

#### OK Response Body

Field | Description
--- | ---
`total` | The total number of hashes performed since the service started.
`average` | The average time in microseconds for hashes performed.

##### Example

```json
{
    "total": 1,
    "average": 123
}
```

### `[ANY] /shutdown`

#### Response Statuses

`200 - OK`: Currently only `200 - OK` is returned.
