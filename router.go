package pressure

import (
	"net/http"
	"regexp"
)

type Router struct {
	Routes []Route
	Logger *Logger
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Declare Variables
	var match *RouteMatch
	var follow bool

	// Loop through routes
	for _, v := range r.Routes {
		match, follow = v.GetMatch(req.URL.Path)
		if follow {
			break
		}
	}

	// 404 if cannot find route
	if !follow || match == nil {
		r.Logger.LogWarning(req.Method, req.URL.Path, "Unable to Resolve Route")
		HTTPErrorNotFound.Write(res)
		return
	}

	// Run Route
	match.runRoute(r.Logger, res, req)
}

func (r *Router) AddRoute(route Route) {
	r.Routes = append(r.Routes, route)
}

type Route interface {
	GetMatch(path string) (*RouteMatch, bool)
}

type RouteMatch struct {
	Path       string
	Controller Controller
	URL        URLCapture
}

func (r *RouteMatch) runRoute(l *Logger, res http.ResponseWriter, req *http.Request) {
	NewHandlerFuncFromController(r.Controller, r.URL, l)(res, req)
}

type URLCapture map[string]string

type URLRoute struct {
	Path          string
	Controller    Controller
	compiledRegex *regexp.Regexp
}

func NewURLRoute(path string, c Controller) URLRoute {
	// Compile the Path
	reg := regexp.MustCompile(path)

	return URLRoute{
		Path:          path,
		Controller:    c,
		compiledRegex: reg,
	}
}

func (r URLRoute) GetMatch(path string) (*RouteMatch, bool) {
	match := r.compiledRegex.FindStringSubmatch(path)
	if match == nil {
		return nil, false
	}

	captures := make(map[string]string)
	for i, name := range r.compiledRegex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}

	return &RouteMatch{path, r.Controller, captures}, true
}

type IncludeRoute struct {
	Path   string
	Routes []Route
}

func (r IncludeRoute) GetMatch(path string) (*RouteMatch, bool) {
	return nil, true
}
