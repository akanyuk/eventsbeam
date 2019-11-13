package beam

import (
	"TEST-LOCAL/events_beam/configuration"
	"TEST-LOCAL/events_beam/kit"
	"fmt"
	"github.com/asticode/go-astilectron"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type beam struct {
	asti    *astilectron.Astilectron
	windows map[string]*window
}

type Beamer interface {
	Init(appName string) error
	ShowWindow(alias string) error
	WaitInterrupt()
}

func NewBeamer() Beamer {
	return &beam{
		windows: map[string]*window{},
	}
}

func (b *beam) Init(appName string) error {
	asti, err := astilectron.New(astilectron.Options{
		AppName: appName,
		//BaseDirectoryPath:  b.basePath,
		VersionAstilectron: configuration.Service.VersionAstilectron,
		VersionElectron:    configuration.Service.VersionElectron,
		AppIconDefaultPath: filepath.Join(kit.ExecutablePath(), "app", "icon-32x32.png"),
		//AppIconDarwinPath:  "<your .icns icon>", // Same here
	})
	if err != nil {
		return err
	}

	b.asti = asti

	if err := b.asti.Start(); err != nil {
		return err
	}

	templates, err := getTemplates(filepath.Join(b.asti.Paths().BaseDirectory(), "templates"))
	if err != nil {
		return err
	}

	for _, template := range templates {
		if err := b.addWindow(template); err != nil {
			return err
		}
	}

	if configuration.Service.Display > 0 {
		b.setDisplay(configuration.Service.Display)
	}

	return nil
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
