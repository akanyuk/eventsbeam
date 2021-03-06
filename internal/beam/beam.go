package beam

import (
	"github.com/akanyuk/eventsbeam/configuration"
	"github.com/akanyuk/eventsbeam/internal/beam/config"
	"github.com/akanyuk/eventsbeam/internal/support"

	"github.com/asticode/go-astilectron"

	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type beam struct {
	asti      *astilectron.Astilectron
	windows   map[string]*window
	templater Templater
}

type Beamer interface {
	Init(appName string) error
	ShowWindow(alias string) error
	WaitInterrupt()

	Templates() []config.Template
	TemplateRead(name string) (config.Template, error)
}

func NewBeamer() Beamer {
	return &beam{
		windows: map[string]*window{},
	}
}

func (b *beam) Init(appName string) error {
	asti, err := astilectron.New(astilectron.Options{
		AppName: appName,
		// BaseDirectoryPath:  b.basePath,
		VersionAstilectron: configuration.Service.VersionAstilectron,
		VersionElectron:    configuration.Service.VersionElectron,
		AppIconDefaultPath: filepath.Join(support.ExecutablePath(), "app", "icon-32x32.png"),
		// AppIconDarwinPath:  "<your .icns icon>", // Same here
	})
	if err != nil {
		return err
	}

	b.asti = asti

	if err := b.asti.Start(); err != nil {
		return err
	}

	b.templater = NewTemplater(filepath.Join(b.asti.Paths().BaseDirectory(), "templates"))
	if err := b.templater.Init(); err != nil {
		return err
	}

	// Перестало работать
	// for _, template := range b.Templates() {
	// 	if err := b.addWindow(template.Name); err != nil {
	// 		return err
	// 	}
	// }

	if configuration.Service.Display > 0 {
		b.setDisplay(configuration.Service.Display)
	}

	return nil
}

func (b *beam) Templates() []config.Template {
	return b.templater.Templates()
}

func (b *beam) TemplateRead(name string) (config.Template, error) {
	return b.templater.Read(name)
}

func (b *beam) ShowWindow(alias string) error {
	w, exist := b.windows[alias]
	if !exist {
		return fmt.Errorf("window not found: %s", alias)
	}

	for a := range b.windows {
		b.windows[a].isActive = false
	}

	if err := w.show(); err != nil {
		return err
	}

	return nil
}

func (b *beam) WaitInterrupt() {
	// Waiting exit signal
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s

		_ = b.asti.Quit()
		os.Exit(0)
	}()

	defer b.asti.Close()
	b.asti.Wait()
}

func (b *beam) setDisplay(d int) {
	var displays = b.asti.Displays()
	if len(displays) < 2 {
		return
	}

	for _, window := range b.windows {
		if err := window.window.MoveInDisplay(displays[1], 0, 0); err != nil {
			log.Printf("unable to change display: %v\n", err)
		}

		if window.isActive {
			_ = window.window.Focus()
		}
	}
}

func (b *beam) addWindow(alias string) error {
	_, exist := b.windows[alias]
	if exist {
		return nil
	}

	w, err := newWindow(b, alias)
	if err != nil {
		return err
	}

	b.windows[alias] = w

	return nil
}
