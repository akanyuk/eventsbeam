package main

import (
	"TEST-LOCAL/events_beam/beam"
	"TEST-LOCAL/events_beam/show"
	"TEST-LOCAL/events_beam/web"

	"fmt"
	"log"
)

func main() {
	beamer := beam.NewBeamer("resources")
	if err := beamer.Init("Events Beam"); err != nil {
		log.Fatalf("unable to initialize beamer: %v\n", err)
	}

	shower := show.NewShower("")
	if err := shower.Init(); err != nil {
		log.Fatalf("unable to initialize shower: %v\n", err)
	}

	go web.Start(beamer, shower)

	fmt.Println("events beam started")

	beamer.WaitInterrupt()
}
