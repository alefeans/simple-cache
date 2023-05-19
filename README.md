# simple-cache

![workflow status](https://github.com/alefeans/simple-cache/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/alefeans/simple-cache/.svg)](https://pkg.go.dev/github.com/alefeans/simple-cache/)

`simple-cache` is a lightweight and thread-safe in-memory key-value cache library for Go. It provides a simple and efficient way to store and retrieve any kind of data with expiration times.

### Installing

```sh
go get github.com/alefeans/simple-cache
```

### Reference

Access [here](https://pkg.go.dev/github.com/alefeans/simple-cache) or execute:

```sh
godoc -http=:6060

# and access http://localhost:6060/pkg/github.com/alefeans/simple-cache/
```

### Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/alefeans/simple-cache"
)

func main() {
	// Create a new cache with a default expiration of 5 minutes and a cleanup interval of 1 minute
	c := cache.New(5*time.Minute, 1*time.Minute)

	// Add an entry to the cache with an expiration of two hours
	c.Set("key", "value", 2*time.Hour)

	// Retrieve the value from the cache
	if value, found := c.Get("key"); found {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Value not found")
	}

	// Delete an entry from the cache
	c.Delete("key")

	// Clear all entries from the cache
	c.Clear()

	// Stop the cleanup goroutine and release resources
	c.Close()
}
```

### Tests

```sh
go test

# to run the benchmarks

go test -bench=. -benchmem
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
