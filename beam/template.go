package beam

import (
	"TEST-LOCAL/eventsbeam/beam/config"

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
