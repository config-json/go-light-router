# Go-light-router

A lightweight HTTP router that serves JSON and static files.

## Getting started

Import the router in your files. The dependency will be automatically added when you try to build/run/test.

```go
import lightrouter "github.com/config-json/go-light-router"
```

### Running go-light-router

A basic example:

```go
// main.go
package main

import (
	lightrouter "github.com/config-json/go-light-router"
)

func main() {
	r := lightrouter.Default()

	r.GET("/ping", func(req *lightrouter.Request, res *lightrouter.Response) {
		res.JSON("pong")
	})

	r.Listen() // Default set to port 8000
}
```

Then run it with:

```
$ go run main.go
```
