# redis (Go Redis client)
[![Build Status](https://travis-ci.org/gbrlsnchs/redis.svg?branch=master)](https://travis-ci.org/gbrlsnchs/redis)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/redis?status.svg)](https://godoc.org/github.com/gbrlsnchs/redis)

## About
This package is a [Redis] client for [Go]. It tries to work similar to the `database/sql` package.
It is so simple that it can be run without any additional setup.

Out of the box, it uses an internal connection pool for reusing connections, if desired.  
If there are available connections in the pool, it simply reuses an open connection,
otherwise it spawns a new connection and, if it can't store it in the pool, closes it.

## Usage
Full documentation [here] (work in progress).

## Example
```go
db, err := redis.Open("localhost:6379")
if err != nil {
	return err
}

if err = db.Exec("SET", "foo", 1); err != nil {
	return err
}

var foo int8
if err = db.QueryRow("GET", "foo").Scan(&foo); err != nil {
	return err
}
log.Println(foo) // prints 1
```

## Contribution
### How to help:
- Pull Requests
- Issues
- Opinions

[Redis]: https://redis.io
[Go]: https://golang.org
[here]: https://godoc.org/github.com/gbrlsnchs/redis