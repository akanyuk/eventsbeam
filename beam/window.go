package beam

import (
	"TEST-LOCAL/eventsbeam/configuration"
	"fmt"
	"github.com/asticode/go-astilectron"
	"path/filepath"
)

type window struct {
	window   *astilectron.Window
	isActive bool
}

func newWindow(b *beam, alias string) (*window, error) {
	alwaysOnTop := true
	frameWindow := true

	if configuration.Service.Debug > 0 {
		alwaysOnTop = false
		frameWindow = false

		fmt.Printf("\033[1;35m debug mode\033[0m\n")
	}

	var w, err = b.asti.NewWindow(filepath.Join(b.asti.Paths().BaseDirectory(), "templates", alias, "index.html"), &astilectron.WindowOptions{
		Fullscreen:      astilectron.PtrBool(true),
		Frame:           astilectron.PtrBool(frameWindow),
		AlwaysOnTop:     astilectron.PtrBool(alwaysOnTop),
		Show:            astilectron.PtrBool(false),
		BackgroundColor: astilectron.PtrStr("#000000"),
	})
	if err != nil {
		return nil, err
	}

	if err := w.Create(); err != nil {
		return nil, err
	}

	if configuration.Service.Debug > 0 {
		_ = w.OpenDevTools()
	}

	return &window{
		window: w,
	}, nil
}

func (w *window) show() error {
	if err := w.window.Show(); err != nil {
		return err
	}

	_ = w.window.Focus()

	w.isActive = true

	return nil
}
