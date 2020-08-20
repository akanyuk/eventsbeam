package configuration

import (
	"log"
)

var Service = struct {
	VersionAstilectron string `default:"0.33.0" usage:"astilectron version"`
	VersionElectron    string `default:"4.0.1" usage:"electron version"`
	BindAddress        string `default:"127.0.0.1:4848" usage:"ip and port for control"`
	Display            int    `default:"0" usage:"start beams on selected monitor"`
	Debug              int    `default:"0" usage:"start beams in debug mode"`
}{}

func Init() {
	if err := Load(&Service, "config.toml"); err != nil {
		log.Printf("configuration load failed: %v\n", err)
	}
}
