package web

import (
	"bitbucket.org/nyuk/eventsbeam/beam"
	"bitbucket.org/nyuk/eventsbeam/configuration"
	"bitbucket.org/nyuk/eventsbeam/kit"
	"bitbucket.org/nyuk/eventsbeam/show"

	"github.com/gorilla/mux"

	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type handler struct {
	beamer beam.Beamer
	comper show.Comper
	slider show.Slider
}

func newHandler(beamer beam.Beamer, shower show.Shower) *handler {
	return &handler{
		beamer: beamer,
		comper: shower.Comper(),
		slider: shower.Slider(),
	}
}

func Start(beamer beam.Beamer, shower show.Shower) {
	handler := newHandler(beamer, shower)
	router := mux.NewRouter()

	router.HandleFunc("/setup/templates", handler.handleTemplates)

	router.HandleFunc("/setup/compos", handler.handleCompos)
	router.HandleFunc("/setup/compo/create", handler.handleCompoCreate)
	router.HandleFunc("/setup/compo/read/{alias}", handler.handleCompoRead)
	router.HandleFunc("/setup/compo/update/{alias}", handler.handleCompoUpdate)
	router.HandleFunc("/setup/compo/delete/{alias}", handler.handleCompoDelete)

	router.HandleFunc("/setup/slides/{compo}", handler.handleSlides)
	router.HandleFunc("/setup/slide/create", handler.handleSlideCreate)
	router.HandleFunc("/setup/slide/read/{id}", handler.handleSlideRead)
	router.HandleFunc("/setup/slide/update/{id}", handler.handleSlideUpdate)
	router.HandleFunc("/setup/slide/delete/{id}", handler.handleSlideDelete)

	// Static resources
	router.HandleFunc("/openapi.yaml", handleOpenapi)
	router.PathPrefix("/openapi").Handler(http.StripPrefix("/openapi", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "app/static", "openapi")))))
	router.PathPrefix("/setup").Handler(http.StripPrefix("/setup", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "app/static", "setup")))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(filepath.Join(kit.ExecutablePath(), "app/static", "control")))))
	http.Handle("/", router)
	router.Use(handleMiddleware)

	fmt.Printf("beam control starting at: %s\n", configuration.Service.BindAddress)

	if err := http.ListenAndServe(configuration.Service.BindAddress, nil); err != nil {
		log.Fatalf("unable to start control server: %v\n", err)
	}
}

func handleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Add("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleOpenapi(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(swaggerJson)
}
