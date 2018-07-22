[![Go Report Card](https://goreportcard.com/badge/myshkin5/hasher)](https://goreportcard.com/report/myshkin5/hasher)

# hasher

## API Documentation

Endpoint |
--- |
[`GET /`](#get-) |
[`GET /health`](#get-health) |
[`POST /hash`](#post-hash) |
[`GET /hash/:requestId`](#get-hashrequestid) |
[`GET /stats`](#get-stats) |

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
