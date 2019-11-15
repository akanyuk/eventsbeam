package main

import (
	"TEST-LOCAL/events-beam/beam"
	"TEST-LOCAL/events-beam/show"
	"TEST-LOCAL/events-beam/web"

	"fmt"
	"log"
)

func main() {
	beamer := beam.NewBeamer()
	if err := beamer.Init("EventsBeam"); err != nil {
		log.Fatalf("unable to initialize beamer: %v\n", err)
	}

	shower := show.NewShower("")
	if err := shower.Init(); err != nil {
		log.Fatalf("unable to initialize shower: %v\n", err)
	}

	go web.Start(beamer, shower)

	fmt.Println("eventsbeam started")

	beamer.WaitInterrupt()
}
