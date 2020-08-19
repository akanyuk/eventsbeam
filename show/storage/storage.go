package storage

import (
	"github.com/akanyuk/eventsbeam/show/config"

	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
	"os"
)

func LoadCompos(configPath string) ([]config.Compo, error) {
	bytes, err := loadData(configPath)
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

	return saveData(data, configPath)
}

func LoadSlides(configPath string) ([]config.Slide, error) {
	bytes, err := loadData(configPath)
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

	return saveData(data, configPath)
}

func loadData(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("file loading error: %v", err)
		}
	}()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func saveData(data []byte, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("file error: %v", err)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}
