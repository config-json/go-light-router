package golightrouter

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RouteHandler func(req *Request, res *Response)

type Route struct {
	method  HTTPMethod
	handler RouteHandler
}

type Router struct {
	// Redirect /foo/ to /foo, return 301 for GET, 307 for other methods
	RedirectTrailingSlash bool
	// Checks if there's other method allowed (returns 405 - Method Not Allowed) else 404 - Not Found
	HandleMethodNotAllowed bool
	Routes                 map[string]Route
}

func Default() *Router {
	r := &Router{
		RedirectTrailingSlash:  true,
		HandleMethodNotAllowed: true,
		Routes:                 make(map[string]Route),
	}
	return r
}

func (r *Router) Listen(ports ...int) {
	// Use 8000 as default port
	var port string
	if len(ports) > 0 {
		port = fmt.Sprintf(":%s", strconv.Itoa(ports[0]))
	} else {
		port = ":8000"
	}

	// Start listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	// TODO: Add callback

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(r, conn)
	}
}

func handleConnection(r *Router, conn net.Conn) {
	defer conn.Close()

	res := Response{}
	req := Request{}

	err := req.readRequest(&conn)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading request:", err)
		return
	}

	// Modifies the request route if the router has RedirectTrailingSlash enabled
	r.trimSlash(&req, &res)

	// Handle route not existent
	r.handleNotFound(&req, &res)
	r.handleMethodNotAllowed(&req, &res)

	// Handler executes if status 404/405 was not added
	if res.Status != http.StatusNotFound && res.Status != http.StatusMethodNotAllowed {
		matchingRoute := r.matchingRoute(&req)
		req.routeToParams(matchingRoute)
		r.Routes[matchingRoute].handler(&req, &res)

		// Handle status automatically if the user/the handler didn't set it already
		if res.Status == 0 {
			res.Status = http.StatusOK
		}
	}

	formattedStatus := fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.Status, http.StatusText(res.Status))
	formattedRes := formattedStatus + formatHeaders(res.Headers) + "\r\n" + res.Body
	_, err = conn.Write([]byte(formattedRes))

	if err != nil {
		fmt.Println("Error writing response:", err)
	}

	if res.File != nil {
		_, err = io.Copy(conn, res.File)
		if err != nil {
			fmt.Println("Error copying file to connection:", err)
		}
		res.File.Close()
	}
}

/************************************/
/********** HTTP METHODS ************/
/************************************/

func (r *Router) GET(route string, handler RouteHandler) {
	r.Routes[route] = Route{
		method:  GET,
		handler: handler,
	}
}

func (r *Router) POST(route string, handler RouteHandler) {
	r.Routes[route] = Route{
		method:  POST,
		handler: handler,
	}
}

func (r *Router) PATCH(route string, handler RouteHandler) {
	r.Routes[route] = Route{
		method:  PATCH,
		handler: handler,
	}
}

func (r *Router) PUT(route string, handler RouteHandler) {
	r.Routes[route] = Route{
		method:  PUT,
		handler: handler,
	}
}

func (r *Router) DELETE(route string, handler RouteHandler) {
	r.Routes[route] = Route{
		method:  DELETE,
		handler: handler,
	}
}

func (r *Router) ServeFile(route, path string) {
	r.Routes[route] = Route{
		method: GET,
		handler: func(req *Request, res *Response) {
			file, err := os.Open(path)
			if err != nil {
				res.Status = http.StatusNotFound
				return
			}

			stat, err := file.Stat()
			if err != nil || stat.IsDir() {
				res.Status = http.StatusNotFound
				return
			}

			res.Header("Content-Length", strconv.Itoa(int(stat.Size())))

			fileExtension := strings.Split(route, ".")[1]
			res.Header("Content-Type", ContentTypeHeader(fileExtension))
			res.File = file
		},
	}
}

/************************************/
/********** ROUTE HANDLING **********/
/************************************/

func removeRouteParams(route string) []string {
	if route == "/" {
		return nil
	}
	splitRoute := strings.Split(route, "/")
	for i, part := range splitRoute {
		if strings.HasPrefix(part, ":") {
			splitRoute[i] = ""
		}
	}
	return splitRoute
}

// Slap some tests onto this bitch
func (r *Router) matchingRoute(req *Request) string {
	for route := range r.Routes {

		if req.route == route {
			return route
		}

		splitRoute := removeRouteParams(route)
		splitReqRoute := strings.Split(req.route, "/")
		if len(splitRoute) != len(splitReqRoute) {
			continue
		}

		matches := true
		for i, part := range splitRoute {
			if part != "" && part != splitReqRoute[i] {
				matches = false
				break
			}
		}
		if matches {
			return route
		}
	}
	return ""
}

// Handle 404
func (r *Router) handleNotFound(req *Request, res *Response) {
	if r.matchingRoute(req) == "" {
		res.Status = http.StatusNotFound
	}
}

// Handle 405
func (r *Router) handleMethodNotAllowed(req *Request, res *Response) {
	matchingRoute := r.matchingRoute(req)
	if !r.HandleMethodNotAllowed && req.method != r.Routes[matchingRoute].method {
		res.Status = http.StatusMethodNotAllowed
	}

	// If it's being handled, change the method to the supported one
	req.method = r.Routes[matchingRoute].method
}

// RedirectTrailingSlash
func (r *Router) trimSlash(req *Request, res *Response) {
	if r.RedirectTrailingSlash && strings.HasSuffix(req.route, "/") && req.route != "/" {
		req.route = req.route[:len(req.route)-1]
		// If GET, return 301, for all other methods, return 307
		if req.method == GET {
			res.Status = http.StatusMovedPermanently
		}
		res.Status = http.StatusTemporaryRedirect
	}
}
