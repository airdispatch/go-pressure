package pressure

type LoginController struct {
	Inner Controller
}

func (g LoginController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if r.Method == "POST" {
		return g.Inner.GetResponse(r, l)
	}
	return nil, HTTPErrorMethodNotAllowed
}

type LoginRequiredController struct {
	Inner Controller
}

func (g LoginRequiredController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if true {
		return g.Inner.GetResponse(r, l)
	} else {
		return nil, nil
	}
}

type LoginControllerChoice struct {
	LoggedIn  Controller
	LoggedOut Controller
}

func (g LoginControllerChoice) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	if true {
		return g.LoggedIn.GetResponse(r, l)
	} else {
		return g.LoggedOut.GetResponse(r, l)
	}
}
