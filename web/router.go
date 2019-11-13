package web

import (
	"TEST-LOCAL/events_beam/beam"
	"TEST-LOCAL/events_beam/configuration"
	"TEST-LOCAL/events_beam/kit"
	"TEST-LOCAL/events_beam/show"

	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
)

type handler struct {
	beamer beam.Beamer
	shower show.Shower
}

func newHandler(beamer beam.Beamer, shower show.Shower) *handler {
	return &handler{
		beamer: beamer,
		shower: shower,
	}
}

func Start(beamer beam.Beamer, shower show.Shower) {
	handler := newHandler(beamer, shower)
	router := mux.NewRouter()

	router.HandleFunc("/setup/compos", handler.handleCompos)
	router.HandleFunc("/setup/compo/create", handler.handleCompoCreate)
	router.HandleFunc("/setup/compo/read/{alias}", handler.handleCompoRead)
	router.HandleFunc("/setup/compo/update/{alias}", handler.handleCompoUpdate)
	router.HandleFunc("/setup/compo/delete/{alias}", handler.handleCompoDelete)

	router.HandleFunc("/setup/slide/create", handler.handleSlideCreate)

	// Static resources
	router.HandleFunc("/openapi/swagger.json", handleOpenapi)
	router.PathPrefix("/openapi").Handler(http.StripPrefix("/openapi", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "static", "openapi")))))
	router.PathPrefix("/setup").Handler(http.StripPrefix("/setup", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "static", "setup")))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "static", "control")))))
	http.Handle("/", router)

	fmt.Printf("beam control starting at: %s\n", configuration.Service.BindAddress)

	if err := http.ListenAndServe(configuration.Service.BindAddress, nil); err != nil {
		log.Fatalf("unable to start control server: %v\n", err)
	}
}

func handleOpenapi(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(swaggerJson)
}
