# redis (Go Redis client)
[![Build Status](https://travis-ci.org/gbrlsnchs/redis.svg?branch=master)](https://travis-ci.org/gbrlsnchs/redis)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/redis?status.svg)](https://godoc.org/github.com/gbrlsnchs/redis)

## About
This package is a simple [Redis] client. It is context-aware and uses a total customizable internal connection pool.

## Usage
Full documentation [here] (work in progress).

## Example
```go
db, err := redis.Open("localhost:6379")
db.SetMaxIdleConns(20)  // reuses up to 20 connections without closing them
db.SetMaxOpensConns(45) // opens up to 45 connections (20 remain open), otherwise waits
if err != nil {
	return err
}

r, err := db.Send("SET", "foo", 1)
if err != nil {
	// handle error
}
log.Print(r.String()) // prints "OK"
log.Print(r.IsOK())   // prints "true"

if r, err = db.Send("GET", "foo"); err != nil {
	// handle error
}
log.Print(r.Int64()) // prints "1"
```

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Redis]: https://redis.io
[Go]: https://golang.org
[here]: https://godoc.org/github.com/gbrlsnchs/redis