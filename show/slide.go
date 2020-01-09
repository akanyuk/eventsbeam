package show

import (
	"TEST-LOCAL/eventsbeam/beam"
	"TEST-LOCAL/eventsbeam/show/config"
	"TEST-LOCAL/eventsbeam/show/storage"
	"TEST-LOCAL/eventsbeam/web/response"
	"path/filepath"
	"sync"
)

const slideConfigFileName = "compo.yaml"

type slide struct {
	sync.RWMutex
	slides     []config.Slide
	configPath string

	beamer beam.Beamer
}

type Slider interface {
	Init() error
	Validate(config.Slide, config.Slide) []response.ErrorItem
	//Compos() []config.Compo
	//Create(config.Slide) error
	//Read(string) (config.Compo, error)
	//Update(alias string, compo config.Compo) error
	//Delete(alias string) error
}

func NewSlider(basePath string, beamer beam.Beamer) Slider {
	return &slide{
		configPath: filepath.Join(basePath, slideConfigFileName),
		beamer:     beamer,
	}
}

func (s *slide) Init() error {
	slides, err := storage.LoadSlides(s.configPath)
	if err != nil {
		slides = []config.Slide{}
	}

	s.Lock()
	defer s.Unlock()

	s.slides = slides

	return nil
}

func (s *slide) Validate(slide config.Slide, oldSlide config.Slide) []response.ErrorItem {
	errorItems := make([]response.ErrorItem, 0)

	_, exist := s.beamer.Template(slide.Template)
	if !exist {
		errorItems = append(errorItems, response.ErrorItem{Code: "template", Message: "template not found"})
	}

	return errorItems
}

func (s *slide) getSlide(id int) (config.Slide, bool) {
	s.Lock()
	defer s.Unlock()

	for _, slide := range s.slides {
		if slide.ID == id {
			return slide, true
		}
	}

	return config.Slide{}, false
}
