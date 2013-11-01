package pressure

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type StaticFileRoute struct {
	URLRoute
	removePattern string
	directory     string
}

func NewStaticFileRoute(pattern string, directory string) StaticFileRoute {
	u := NewURLRoute(pattern, nil)
	return StaticFileRoute{u, pattern, directory}
}

func (s StaticFileRoute) GetMatch(path string) (*RouteMatch, bool) {
	_, b := s.URLRoute.GetMatch(path)
	if !b {
		return nil, false
	}

	// We have a match... Check for File
	relative_location := s.URLRoute.compiledRegex.ReplaceAllString(path, "")
	file_location := filepath.Join(s.directory, relative_location)

	// Make sure that we aren't just getting the same directory
	if filepath.Clean(file_location) == filepath.Clean(s.directory) {
		return nil, false
	}

	fi, err := os.Open(file_location)
	if err != nil {
		// Couldn't find file, return false
		return nil, false
	}

	return &RouteMatch{
		Path:       path,
		Controller: StaticFileViewController{fi, path},
	}, true
}

type StaticFileViewController struct {
	File *os.File
	Name string
}

func (s StaticFileViewController) GetResponse(r *Request, l *Logger) (View, *HTTPError) {
	l.LogDebug("Serving File", s.Name)
	return s, nil
}

func (b StaticFileViewController) StatusCode() int {
	return 200
}

func (b StaticFileViewController) WriteBody(w io.Writer) {
	io.Copy(w, b.File)
	b.File.Close()
}

func (b StaticFileViewController) ContentLength() int {
	info, _ := b.File.Stat()
	return int(info.Size())
}

func (b StaticFileViewController) ContentType() string {
	return mime.TypeByExtension("." + strings.Split(b.Name, ".")[1])
}

func (b StaticFileViewController) Headers() ViewHeaders {
	return nil
}
