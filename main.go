package main

import (
	"bitbucket.org/nyuk/eventsbeam/beam"
	"bitbucket.org/nyuk/eventsbeam/show"
	"bitbucket.org/nyuk/eventsbeam/web"

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
