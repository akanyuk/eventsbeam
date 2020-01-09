package storage

import (
	"TEST-LOCAL/eventsbeam/show/config"

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
			log.Printf("compos loading error: %v", err)
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
			log.Printf("compos saving error: %v", err)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}

func LoadSlides(configPath string) ([]config.Slide, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("slides loading error: %v", err)
		}
	}()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	slides := make([]config.Slide, 0)
	if err = yaml.Unmarshal(bytes, &slides); err != nil {
		return nil, err
	}

	return slides, nil
}

func SaveSlides(slides []config.Slide, configPath string) error {
	data, err := yaml.Marshal(slides)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("slides saving error: %v", err)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}
