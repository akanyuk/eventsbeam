package show

import (
	"github.com/akanyuk/eventsbeam/internal/beam"
	"github.com/akanyuk/eventsbeam/internal/support"

	"path/filepath"
)

type show struct {
	basePath string
	comper   Comper
	slider   Slider
}

type Shower interface {
	Init(beamer beam.Beamer) error
	Comper() Comper
	Slider() Slider
}

func NewShower(basePath string) Shower {
	return &show{
		basePath: filepath.Join(support.ExecutablePath(), basePath),
	}
}

func (s *show) Init(beamer beam.Beamer) error {
	s.comper = NewComper(s.basePath)
	if err := s.comper.Init(); err != nil {
		return err
	}

	s.slider = NewSlider(s.basePath, beamer, s.comper)
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
