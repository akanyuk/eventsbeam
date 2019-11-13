package beam

import (
	"io/ioutil"
)

func getTemplates(dir string) ([]string, error) {
	var templates []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			templates = append(templates, f.Name())
		}
	}

	return templates, nil
}
