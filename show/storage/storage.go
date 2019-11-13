package storage

import (
	"TEST-LOCAL/events_beam/show/config"

	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
	"os"
)

func LoadCompos(configPath string) ([]config.Compo, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("show loading error: %v", err)
		}
	}()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	compos := make([]config.Compo, 0)
	if err = yaml.Unmarshal(bytes, &compos); err != nil {
		return nil, err
	}

	return compos, nil
}

func SaveCompos(compos []config.Compo, configPath string) error {
	data, err := yaml.Marshal(compos)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("show saving error: %v", err)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}
