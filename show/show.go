package show

import (
	"TEST-LOCAL/eventsbeam/kit"

	"path/filepath"
)

type show struct {
	basePath string
	comper   Comper
	slider   Slider
}

type Shower interface {
	Init() error
	Comper() Comper
	Slider() Slider
}

func NewShower(basePath string) Shower {
	return &show{
		basePath: filepath.Join(kit.ExecutablePath(), basePath),
	}
}

func (s *show) Init() error {
	s.comper = NewComper(s.basePath)
	if err := s.comper.Init(); err != nil {
		return err
	}

	s.slider = NewSlider(s.basePath)
	if err := s.slider.Init(); err != nil {
		return err
	}

	return nil
}

func (s *show) Comper() Comper {
	return s.comper
}

func (s *show) Slider() Slider {
	return s.slider
}
