package pressure

import (
	"errors"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/hoisie/web"
	"strings"
	"os"
	"io"
	"fmt"
	"html/template"
	"path/filepath"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/sessions"
)

type Server struct {
	workingDirectory string
	templateDirectory string
	parsedTemplates map[string]template.Template
	Port string
	DbMap *gorp.DbMap
	WebServer *web.Server
	CookieAuthKey []byte
	CookieEncryptKey []byte
	MainSessionName string
	sessionStore *sessions.CookieStore
	Mailserver string
}

func (s *Server) ConfigServer() {
	s.loadConstants()
	s.WebServer = web.NewServer()
	s.loadTemplates(s.getTemplatePath(), "", true)
	s.WebServer.Config.StaticDir = s.workingDirectory + "/static"
	s.sessionStore = sessions.NewCookieStore(s.CookieAuthKey)
}

func (s *Server) RunServer() {
	s.WebServer.Run("0.0.0.0:" + s.Port)
}

// EVERYTHING BELOW THIS LINE IS BOILERPLATE

func (s *Server) GetMainSession(ctx *web.Context) (*sessions.Session, error) {
	return s.GetSessionForCTXAndName(ctx, s.MainSessionName)
}

func (s *Server) GetSessionForCTXAndName(ctx *web.Context, name string) (*sessions.Session, error) {
	return s.sessionStore.Get(ctx.Request, name)
}

func SaveSessionWithContext(sesh *sessions.Session, ctx *web.Context) error {
	return sesh.Save(ctx.Request, ctx.ResponseWriter)
}

func OpenDatabaseFromURL(url string) (*sql.DB, error) {
	// Take out the Beginning
	url_parts := strings.Split(url,"://")

	if len(url_parts) != 2 || url_parts[0] != "postgres" {
		return nil, errors.New("Database URL is not Postgres")
	}
	url = url_parts[1]

	url_parts = strings.Split(url, "@")
	username := "postgres"

	if len(url_parts) == 2 {
		username = url_parts[0]
		url = url_parts[1]
	} else {
		url = url_parts[0]
	}

	url_parts = strings.Split(url, "/")
	if len(url_parts) != 2 {
		return nil, errors.New("Database URL does not include Database Name")
	}

	db_name := url_parts[1]
	url_parts = strings.Split(url_parts[0], ":")
	port := "5432"

	if len(url_parts) == 2 {
		port = url_parts[1]
	}
	url = url_parts[0]

	return sql.Open("postgres", "user=" + username + " dbname=" + db_name + " host=" + url + " port=" + port + " sslmode=disable")
}

type TemplateView func(ctx *web.Context)
type WildcardTemplateView func(ctx *web.Context, val string)

func (s *Server) DisplayTemplate(templateName string) TemplateView {
	return func(ctx *web.Context) {
		s.WriteTemplateToContext(templateName, ctx, nil)
	}
}

func (s *Server) loadConstants() {
	temp_dir := os.Getenv("WORK_DIR")
	if temp_dir == "" {
		temp_dir, _ = os.Getwd()
	}
	s.workingDirectory = temp_dir

	s.templateDirectory = "templates"
	s.parsedTemplates = make(map[string]template.Template)

	s.CookieAuthKey = []byte(os.Getenv("COOKIE_AUTH"))
	if os.Getenv("COOKIE_AUTH") == "" {
		s.CookieAuthKey = []byte("secret-auth")
	}

	s.CookieEncryptKey = []byte(os.Getenv("COOKIE_ENCRYPTION"))
	if os.Getenv("COOKIE_ENCRYPTION") == "" {
		s.CookieEncryptKey = []byte("secret-encryption-key")
	}
}

func (s *Server) loadTemplates(folder string, append string, linkWithBase bool) {
	// Start looking through the original directory
	dirname := folder + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println("Unable to Read Templates Folder: " + dirname)
		os.Exit(1)
	}
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Loop over all files
	for _, fi := range files {
		if fi.IsDir() {
			// Call yourself if you find more templates
			s.loadTemplates(dirname + fi.Name(), append + fi.Name() + string(filepath.Separator), linkWithBase)
		} else {
			// Parse templates here
			if fi.Name() == "base.html" {
				continue
			}

			linkingOption := linkWithBase
			effectiveName := fi.Name()

			if strings.HasPrefix(fi.Name(), "__nolink ") {
				linkingOption = false
				effectiveName = strings.TrimPrefix(fi.Name(), "__nolink ")
			}

			templateName := append + effectiveName
			s.parseTemplate(templateName, s.getSpecificTemplatePath(append + fi.Name()), linkingOption)
		}
	}
}

func (s *Server) parseTemplate(templateName string, filename string, linkWithBase bool) {
	var tmp *template.Template
	var err error

	if linkWithBase {
		tmp, err = template.New(templateName).ParseFiles(filename, s.getSpecificTemplatePath("base.html"))
	} else {
		tmp, err = template.New(templateName).ParseFiles(filename)	
	}
	if err != nil {
		fmt.Println("Unable to parse template", templateName, "at", filename)
		fmt.Println(err)
		os.Exit(1)
	}

	s.parsedTemplates[templateName] = *tmp
}

func blankResponse() string {
	return ""
}

func (s *Server) writeHeaders(ctx *web.Context) {
}

func (s *Server) getSpecificTemplatePath(templateName string) string {
	return appendPathComponents(s.getTemplatePath(), templateName)
}

func (s *Server) getTemplatePath() string {
	return s.getPath(s.templateDirectory)
}

func (s *Server) getPath(path string) string {
	return appendPathComponents(s.workingDirectory, path)
}

func appendPathComponents(pathComponents ...string) string {
	output := ""
	for i, v := range(pathComponents) {
		if i > 0 {
			output += string(filepath.Separator)
		}
		output += v
	}
	return output
}

func (s *Server) WriteTemplateToBuffer(templatename string, individual string, buffer io.Writer, data interface{}) (error) {
	template, ok := s.parsedTemplates[templatename]
	if !ok {
		return errors.New("Could Not Find Template")
	}
	err := template.ExecuteTemplate(buffer, individual, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) WriteTemplateToContext(templatename string, ctx *web.Context, data interface{}) {
	err := s.WriteTemplateToBuffer(templatename, "base", ctx, data)
	if err != nil {
		displayErrorPage(ctx, "Unable to load template. Template: " + templatename)
	}
}

func (s *Server) getFileContents(filename string) (*os.File, error) {
	file, err := os.Open(s.getPath(filename))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Server) writeFileToContext(filename string, ctx *web.Context) {
	file, err := s.getFileContents(filename)
	if err != nil {
		displayErrorPage(ctx, "Unable to open file. File: " + s.getPath(filename))
		return
	}
	_, err = io.Copy(ctx, file)
	if err != io.EOF && err != nil {
		displayErrorPage(ctx, "Unable to Copy into Buffer. File: " + s.getPath(filename))
		return
	}
}

func displayErrorPage(ctx *web.Context, error string) {
	ctx.WriteString("<!DOCTYPE html><html><head><title>Project Error</title></head>")
	ctx.WriteString("<body><h1>Application Error</h1>")
	ctx.WriteString("<p>" + error + "</p>")
	ctx.WriteString("</body></html>")
}