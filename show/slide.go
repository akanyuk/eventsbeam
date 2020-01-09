package show

import (
	"TEST-LOCAL/eventsbeam/beam"
	"TEST-LOCAL/eventsbeam/show/config"
	"TEST-LOCAL/eventsbeam/show/storage"
	"TEST-LOCAL/eventsbeam/web/response"
	"errors"
	"path/filepath"
	"sync"
)

const slideConfigFileName = "slide.yaml"

type slide struct {
	sync.RWMutex
	slides     []config.Slide
	configPath string

	beamer beam.Beamer
	comper Comper
}

type Slider interface {
	Init() error
	Validate(config.Slide) []response.ErrorItem
	//Compos() []config.Compo
	Create(config.Slide) error
	//Read(string) (config.Compo, error)
	//Update(alias string, compo config.Compo) error
	//Delete(alias string) error
}

func NewSlider(basePath string, beamer beam.Beamer, comper Comper) Slider {
	return &slide{
		configPath: filepath.Join(basePath, slideConfigFileName),
		beamer:     beamer,
		comper:     comper,
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

func (s *slide) Validate(slide config.Slide) []response.ErrorItem {
	errorItems := make([]response.ErrorItem, 0)

	if _, err := s.beamer.TemplateRead(slide.Template); err != nil {
		errorItems = append(errorItems, response.ErrorItem{Code: "template", Message: err.Error()})
	}

	if _, err := s.comper.Read(slide.Compo); err != nil {
		errorItems = append(errorItems, response.ErrorItem{Code: "compo", Message: err.Error()})
	}

	return errorItems
}

func (s *slide) Create(slide config.Slide) error {
	if _, exist := s.getSlide(slide.ID); exist {
		return errors.New("slide already exist")
	}

	slide.ID = s.nextID()
	slide.Position = s.nextPosition(slide.Compo)
	s.slides = append(s.slides, slide)

	s.Lock()
	defer s.Unlock()

	if err := storage.SaveSlides(s.slides, s.configPath); err != nil {
		return err
	}

	return nil
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

func (s *slide) nextID() int {
	s.Lock()
	defer s.Unlock()

	id := 1

	for _, slide := range s.slides {
		if slide.ID >= id {
			id = slide.ID + 1
		}
	}

	return id
}

func (s *slide) nextPosition(compo string) int {
	s.Lock()
	defer s.Unlock()

	position := 1

	for _, slide := range s.slides {
		if slide.Compo == compo && slide.Position >= position {
			position = slide.Position + 1
		}
	}

	return position
}
