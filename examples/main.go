package main

import (
	golightrouter "github.com/config-json/go-light-router"
)

func main() {
	r := golightrouter.Default()

	r.GET("/foobar/:user", func(req *golightrouter.Request, res *golightrouter.Response) {
		user := req.GetParam("user")
		body := map[string]string{
			"foobar": user,
		}

		res.JSON(body)
	})

	r.Listen()
}
