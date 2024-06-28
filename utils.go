package golightrouter

import (
	"net/http"
	"strings"
)

func (r *Router) NotFound(req *Request, res *Response) {
	if r.MatchingRoute(req) == "" {
		res.Status = http.StatusNotFound
	}
}

func (r *Router) MethodNotAllowed(req *Request, res *Response) {
	matchingRoute := r.MatchingRoute(req)
	if r.HandleMethodNotAllowed && req.method != r.Routes[matchingRoute].method {
		res.Status = http.StatusMethodNotAllowed
	}
}

func (r *Router) TrimSlash(req *Request, res *Response) {
	if r.RedirectTrailingSlash && strings.HasSuffix(req.route, "/") && req.route != "/" {
		req.route = req.route[:len(req.route)-1]
		// If GET, return 301, for all other methods, return 307
		if req.method == GET {
			res.Status = http.StatusMovedPermanently
		}
		res.Status = http.StatusTemporaryRedirect
	}
}

func FormatHeaders(headers map[string]string) string {
	result := ""
	for key, value := range headers {
		result += key + ": " + value + "\r\n"
	}
	return result
}
