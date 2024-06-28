package golightrouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Status  int
	Headers map[string]string
	Body    string // At least for JSON
	File    *os.File
}

// TODO: cookies

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// "" as the value deletes the header
func (r *Response) Header(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	if value == "" {
		delete(r.Headers, key)
		return
	}
	r.Headers[key] = value
}

func (r *Response) JSON(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		r.Status = http.StatusInternalServerError
	}

	r.Header("Content-Type", "application/json")
	r.Header("Content-Length", strconv.Itoa(len(jsonData)))
	r.Body = string(jsonData)
	r.Status = http.StatusOK
}
