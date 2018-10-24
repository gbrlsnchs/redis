# redis (Redis client for Go)
[![Build Status](https://travis-ci.org/gbrlsnchs/redis.svg?branch=master)](https://travis-ci.org/gbrlsnchs/redis)
[![Sourcegraph](https://sourcegraph.com/github.com/gbrlsnchs/redis/-/badge.svg)](https://sourcegraph.com/github.com/gbrlsnchs/redis?badge)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/redis?status.svg)](https://godoc.org/github.com/gbrlsnchs/redis)
[![Minimal Version](https://img.shields.io/badge/minimal%20version-go1.10%2B-5272b4.svg)](https://golang.org/doc/go1.10)

## About
This package is a simple [Redis](https://redis.io) client for [Go](https://golang.org). It is context-aware and uses a resizable connection pool internally.

## Usage
Full documentation [here](https://godoc.org/github.com/gbrlsnchs/redis).

### Installing
#### Go 1.10
`vgo get -u github.com/gbrlsnchs/redis`
#### Go 1.11 or after
`go get -u github.com/gbrlsnchs/redis`

### Importing
```go
import (
	// ...

	"github.com/gbrlsnchs/redis"
)
```

### Pinging the database
```go
db, err := redis.Open(":6379")
if err != nil {
	// handle error
}
if _, err = db.Ping(); err != nil {
	// handle error
}
```

### Configuring the connection pool
```go
db, err := redis.Open(":6379")
if err != nil {
	// handle error
}
db.SetMaxIdleConns(20) // reuses up to 20 connections without closing them
db.SetMaxOpenConns(45) // opens up to 45 connections (20 remain open), otherwise waits
```

### Sending commands
```go
r, err := db.Send("SET", "foo", 1)
if err != nil {
	// handle error
}
fmt.Println(r.String()) // prints "OK"
fmt.Println(r.IsOK())   // prints "true"

if r, err = db.Send("GET", "foo"); err != nil {
	// handle error
}
fmt.Println(r.Int64()) // prints "1"
```

## Contributing
### How to help
- For bugs and opinions, please [open an issue](https://github.com/gbrlsnchs/connpool/issues/new)
- For pushing changes, please [open a pull request](https://github.com/gbrlsnchs/connpool/compare)
