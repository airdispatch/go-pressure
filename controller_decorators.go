package pressure

type GETController struct {
	Inner Controller
}

func (g GETController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "GET" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, HTTPErrorMethodNotAllowed
}

type POSTController struct {
	Inner Controller
}

func (g POSTController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "POST" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, HTTPErrorMethodNotAllowed
}
