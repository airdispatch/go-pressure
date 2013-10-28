package pressure

type GETController struct {
	Inner Controller
}

func (g GETController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "GET" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, &HTTPError{405, "405: Method Not Allowed"}
}

type POSTController struct {
	Inner Controller
}

func (g POSTController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "POST" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, &HTTPError{405, "405: Method Not Allowed"}
}

type LoginController struct {
	Inner Controller
}

func (g LoginController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "POST" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, &HTTPError{405, "405: Method Not Allowed"}
}
