package redis_test

import (
	"fmt"

	"github.com/gbrlsnchs/redis"
)

func Example() {
	db, err := redis.Open("localhost:6379")
	if err != nil {
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
	// Output:
	// bar
}
