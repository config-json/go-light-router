package golightrouter

import (
	"bufio"
	"net"
	"strings"
)

type Request struct {
	method  HTTPMethod
	route   string
	params  map[string]string
	headers map[string]string
	// body    string
}

func (req *Request) GetParam(key string) string {
	return req.params[key]
}

func (req *Request) GetHeader(key string) string {
	return req.headers[key]
}

// "" as the value deletes the param
func (r *Request) setParam(key, value string) {
	if r.params == nil {
		r.params = make(map[string]string)
	}
	if value == "" {
		delete(r.params, key)
		return
	}
	r.params[key] = value
}

// Looks for the params (prefixed by :) of the request's matching route and stores their values into a map
func (req *Request) routeToParams(route string) {
	splitReqRoute := strings.Split(req.route, "/")
	splitRoute := strings.Split(route, "/")

	for i, part := range splitRoute {
		if strings.HasPrefix(part, ":") {
			req.setParam(part[1:], splitReqRoute[i])
		}
	}
}

// Reads the request from the connection
func (req *Request) readRequest(conn *net.Conn) error {
	reader := bufio.NewReader(*conn)

	// Read the request line (Method, Route)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	split := strings.Split(requestLine, " ")
	req.method = HTTPMethod(split[0])
	req.route = split[1]

	// Headers
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Empty line means we're done
		if header == "\r\n" {
			break
		}

		splitHeader := strings.Split(header, ": ")

		if req.headers == nil {
			req.headers = make(map[string]string)
		}
		req.headers[splitHeader[0]] = splitHeader[1]
	}

	// Read the body
	// TODO: fix
	/*
		body, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		req.body = body
	*/

	return nil
}
