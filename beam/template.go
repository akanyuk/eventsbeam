package beam

import (
	"TEST-LOCAL/eventsbeam/beam/config"
	"errors"

	"io/ioutil"
	"sync"
)

type template struct {
	sync.RWMutex
	templates []config.Template
	baseDir   string
}

type Templater interface {
	Init() error
	Templates() []config.Template
	Read(name string) (config.Template, error)
}

func NewTemplater(baseDir string) Templater {
	return &template{
		baseDir: baseDir,
	}
}

func (t *template) Init() error {
	templates, err := getTemplates(t.baseDir)
	if err != nil {
		return err
	}

	t.Lock()
	defer t.Unlock()

	t.templates = make([]config.Template, 0)
	for _, name := range templates {
		t.templates = append(t.templates, config.Template{
			Name: name,
		})
	}

	return nil
}

func (t *template) Templates() []config.Template {
	return t.templates
}

func (t *template) Read(name string) (config.Template, error) {
	t.Lock()
	defer t.Unlock()

	for _, template := range t.templates {
		if template.Name == name {
			return template, nil
		}
	}

	return config.Template{}, errors.New("template not found")
}

func getTemplates(dir string) ([]string, error) {
	var templates []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			templates = append(templates, f.Name())
		}
	}

	return templates, nil
}
