package pressure

import (
	"net/http"
	"strconv"
)

type HTTPError struct {
	Code int
	Text string
}

func (h *HTTPError) Write(res http.ResponseWriter) {
	http.Error(res, h.Text, h.Code)
}

func (h *HTTPError) Print() string {
	return strconv.Itoa(h.Code) + " (" + h.Text + ")"
}

var HTTPErrorMethodNotAllowed = &HTTPError{405, "405: Method Not Allowed"}
var HTTPErrorNotFound = &HTTPError{404, "404: Unable to resolve route."}
