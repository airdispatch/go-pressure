package pressure

import (
	"net/http"
	"strconv"
)

type Request struct {
	*http.Request
	Form   map[string][]string
	URL    URLCapture
	Method string
	Path   string
	r      *http.Request
}

type Controller interface {
	GetResponse(*Request, *Logger) (View, *HTTPError)
}

// ---- CONTROLLER METHODS ---- //

func NewHandlerFuncFromController(c Controller, capture URLCapture, l *Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Get Request
		pressureRequest := NewRequestFromHTTPRequest(req)
		pressureRequest.URL = capture

		// Log Request
		l.LogDebug(pressureRequest.Method, pressureRequest.Path)

		// Get Response View
		pressureResponse, err := c.GetResponse(pressureRequest, l)
		if err != nil {
			l.LogError("Error at"+pressureRequest.Path, "::", err.Print())
			err.Write(res)
			return
		}

		// Set Headers
		res.Header().Add("Content-Type", pressureResponse.ContentType())
		contentLength := strconv.Itoa(pressureResponse.ContentLength())
		res.Header().Add("Content-Length", contentLength)

		// Set Custom Headers
		for i, v := range pressureResponse.Headers() {
			res.Header().Add(i, v)
		}

		if c, ok := pressureResponse.(CookieView); ok {
			c.AddCookies(res)
		}

		res.WriteHeader(pressureResponse.StatusCode())

		// Write Body
		pressureResponse.WriteBody(res)
	}
}

func NewRequestFromHTTPRequest(req *http.Request) *Request {
	req.ParseForm()
	return &Request{
		Request: req,
		Form:    req.Form,
		Method:  req.Method,
		Path:    req.URL.Path,
		r:       req,
	}
}
