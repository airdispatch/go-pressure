package pressure

import "net/http"

type ViewHeaders map[string]string

type View interface {
	WriteBody(http.ResponseWriter)
	StatusCode() int
	ContentType() string
	ContentLength() int
	Headers() ViewHeaders
}

type CookieView interface {
	AddCookies(http.ResponseWriter)
}

type BasicView struct {
	Status int
	Text   string
	IsHTML bool
}

func (b BasicView) StatusCode() int {
	return b.Status
}

func (b BasicView) WriteBody(w http.ResponseWriter) {
	w.Write([]byte(b.Text))
	return
}

func (b BasicView) ContentLength() int {
	return len(b.Text)
}

func (b BasicView) ContentType() string {
	if b.IsHTML {
		return "text/html"
	}
	return "text/plain"
}

func (b BasicView) Headers() ViewHeaders {
	return nil
}

type HTMLView struct {
	BasicView
}

func NewHTMLView(text string) HTMLView {
	b := BasicView{}
	b.Status = 200
	b.Text = text
	b.IsHTML = true
	return HTMLView{b}
}

type FormView struct{}

type ModelDetailView struct{}
type ModelListView struct{}

type JSONView struct{}
