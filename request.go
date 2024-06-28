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
	body    string
}

// "" as the value deletes the param
func (r *Request) Param(key, value string) {
	if r.params == nil {
		r.params = make(map[string]string)
	}
	if value == "" {
		delete(r.params, key)
		return
	}
	r.params[key] = value
}

func (req *Request) GetParam(key string) string {
	return req.params[key]
}

func (req *Request) RouteToParams(route string) {
	splitReqRoute := strings.Split(req.route, "/")
	splitRoute := strings.Split(route, "/")

	for i, part := range splitRoute {
		if strings.HasPrefix(part, ":") {
			req.Param(part[1:], splitReqRoute[i])
		}
	}
}

func (req *Request) ReadRequest(conn *net.Conn) error {
	reader := bufio.NewReader(*conn)

	// Read the request line (Method, Route)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	split := strings.Split(requestLine, " ")
	req.method = HTTPMethod(split[0])
	req.route = split[1]

	// Read the headers
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// Empty line means we're done, next is body
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
	body, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	req.body = body

	return nil
}
