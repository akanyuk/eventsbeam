package main

import (
	"github.com/akanyuk/eventsbeam/beam"
	"github.com/akanyuk/eventsbeam/show"
	"github.com/akanyuk/eventsbeam/web"

	"fmt"
	"log"
)

func main() {
	beamer := beam.NewBeamer()
	if err := beamer.Init("EventsBeam"); err != nil {
		log.Fatalf("unable to initialize beamer: %v\n", err)
	}

	shower := show.NewShower("")
	if err := shower.Init(beamer); err != nil {
		log.Fatalf("unable to initialize shower: %v\n", err)
	}

	go web.Start(beamer, shower)

	fmt.Println("eventsbeam started")

	beamer.WaitInterrupt()
}
