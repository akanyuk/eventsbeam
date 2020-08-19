package beam

import (
	"encoding/json"
	"errors"
	"github.com/akanyuk/eventsbeam/beam/config"
	"os"
	"path/filepath"

	"io/ioutil"
	"sync"
)

const templateParamsFilename = "params.json"

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

	t.templates = templates
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

func getTemplates(dir string) ([]config.Template, error) {
	var templates []config.Template

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			params, err := getTemplateParams(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}

			templates = append(templates, config.Template{
				Name:   f.Name(),
				Params: params,
			})
		}
	}

	return templates, nil
}

func getTemplateParams(dir string) (map[string]config.TemplateParam, error) {
	filename := filepath.Join(dir, templateParamsFilename)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return map[string]config.TemplateParam{}, nil
	} else if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	template := config.Template{}
	if err := json.Unmarshal(content, &template); err != nil {
		return nil, err
	}

	return template.Params, nil
}
