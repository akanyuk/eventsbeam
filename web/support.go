package web

import (
	"bitbucket.org/nyuk/eventsbeam/show/config"
	"encoding/json"
	"log"
	"net/http"
)

func ExtractCompo(r *http.Request) (config.Compo, error) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error with closing request body: %v", err)
		}
	}()

	var compo config.Compo
	err := json.NewDecoder(r.Body).Decode(&compo)

	return compo, err
}

func ExtractSlide(r *http.Request) (config.Slide, error) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error with closing request body: %v", err)
		}
	}()

	var slide config.Slide
	err := json.NewDecoder(r.Body).Decode(&slide)

	return slide, err
}
