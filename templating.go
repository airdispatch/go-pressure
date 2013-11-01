package pressure

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type TemplateView struct {
	BasicView
}

func (t *TemplateEngine) NewTemplateView(templateName string, context interface{}) TemplateView {
	b := BasicView{
		Status: 200,
		IsHTML: true,
	}

	temp, err := t.GetTemplateNamed(templateName)
	if err != nil {
		t.LogError("Unable to locate template", templateName)
		panic("Cannot recover from errror.")
	}

	var buf *bytes.Buffer = &bytes.Buffer{}
	temp.ExecuteTemplate(buf, "base", context)

	b.Text = buf.String()

	tv := TemplateView{b}
	return tv
}

type TemplateEngine struct {
	Directory        string
	CachedTemplates  map[string]template.Template
	BaseTemplateName string
	Debug            bool
	*Logger
}

func (s *Server) CreateTemplateEngine(directory string, base string) *TemplateEngine {
	t := &TemplateEngine{}
	t.Directory = directory
	t.CachedTemplates = make(map[string]template.Template)
	t.BaseTemplateName = base
	t.Debug = s.Debug
	t.Logger = s.Logger

	t.CompileTemplates()

	return t
}

func (s *TemplateEngine) GetTemplateNamed(templateName string) (*template.Template, error) {
	if s.Debug {
		s.parseTemplate(templateName, filepath.Join(s.Directory, templateName), true)
	}

	t, ok := s.CachedTemplates[templateName]
	if !ok {
		return nil, errors.New("Couldn't find template.")
	}

	return &t, nil
}

func (s *TemplateEngine) CompileTemplates() {
	s.loadTemplatesFromFolder(s.Directory, true)
}

func (s *TemplateEngine) loadTemplatesFromFolder(folder string, linkWithBase bool) {
	// Walk through original Directory
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if folder == path {
			return nil
		}

		if !info.IsDir() {
			// Parse templates here
			if info.Name() == s.BaseTemplateName {
				return nil
			}

			linkingOption := linkWithBase
			effectivePath := path

			if strings.HasPrefix(info.Name(), "__nolink ") {
				linkingOption = false
				effectivePath = strings.Replace(path, "__nolink ", "", 1)
			}

			templateName, _ := filepath.Rel(folder, effectivePath)
			s.parseTemplate(templateName, path, linkingOption)
		}
		return nil
	})
}

func (s *TemplateEngine) parseTemplate(templateName string, filename string, linkWithBase bool) {
	s.LogDebug("Parsing", filename, "as", templateName)
	var tmp *template.Template
	var err error

	if linkWithBase {
		tmp, err = template.New(templateName).ParseFiles(filename, filepath.Join(s.Directory, s.BaseTemplateName))
	} else {
		tmp, err = template.New(templateName).ParseFiles(filename)
	}

	if err != nil {
		s.LogError("Unable to parse template", filename)
		s.LogError(err)
		return
	}

	s.CachedTemplates[templateName] = *tmp
}
