package redis_test

import (
	"fmt"

	"github.com/gbrlsnchs/redis"
)

func Example() {
	var (
		db  *redis.DB
		err error
	)

	if db, err = redis.Open("localhost", 6379, 10); err != nil {
		// Handler error...
	}

	if err = db.Exec("SET", "foo", "bar"); err != nil {
		// Handle error...
	}

	var foo string

	if err = db.QueryRow("GET", "foo").Scan(&foo); err != nil {
		// Handle error...
	}

	fmt.Println(foo)
}
