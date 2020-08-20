package web

import (
	"github.com/akanyuk/eventsbeam/beam"
	"github.com/akanyuk/eventsbeam/configuration"
	"github.com/akanyuk/eventsbeam/show"

	"github.com/gorilla/mux"

	"fmt"
	"log"
	"net/http"
)

type staticResource struct {
	path string
	f    func(http.ResponseWriter, *http.Request)
}

var staticResources []staticResource

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
	router := mux.NewRouter().StrictSlash(true)

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

	// Adding generated static resources
	for _, resource := range staticResources {
		router.HandleFunc(resource.path, resource.f)
	}

	router.Use(handleMiddleware)

	http.Handle("/", router)

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
