package main

import lightrouter "github.com/config-json/go-light-router"

func main() {
	r := lightrouter.Default()

	r.GET("/ping", func(req *lightrouter.Request, res *lightrouter.Response) {
		res.JSON("pong")
	})

	r.GET("/foobar/:user", func(req *lightrouter.Request, res *lightrouter.Response) {
		user := req.GetParam("user")
		body := map[string]string{
			"foobar": user,
		}

		res.JSON(body)
	})

	r.Listen()
}
