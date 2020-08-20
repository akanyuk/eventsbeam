//go:generate go run swagger.go -swagger-work-dir ../../ -swagger-output-file ../../web/static/sources/swagger.json
//go:generate swagger serve ../../web/static/sources/swagger.json /flavor:redoc

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/akanyuk/eventsbeam/generators/swagger/options"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	projectName        = "EventsBeam API"
	projectDescription = "# EventsBeam\nСистема управления показом слайдов. Интегрируется с [events.retroscene.org](http://events.retroscene.org)"
)

func main() {
	var swaggerWorkDir = flag.String("swagger-work-dir", ".", "working directory for swagger generation")
	var swaggerOutputFile = flag.String("swagger-output-file", "./web/static/sources/swagger.json", "output swagger file")

	flag.Parse()

	Generate(*swaggerWorkDir, *swaggerOutputFile)
}

func Generate(workDir, outputFile string, opt ...options.Option) {
	callOptions := options.Do(opt...)

	metaFile, err := ioutil.TempFile("", "meta.yaml")
	if err != nil {
		log.Fatalf("create tmp meta file error: %v", err)
	}

	if err := swaggerMeta(metaFile.Name()); err != nil {
		log.Fatalf("swagger meta generation error: %v", err)
	}

	if err := generateSwagger(workDir, outputFile, metaFile.Name(), callOptions.SkipRemoveUnusedModels, callOptions.ExcludePackages); err != nil {
		log.Fatalf("generate swagger error: %v", err)
	}

	_ = os.Remove(metaFile.Name())
}

func generateSwagger(workDir string, outputFile string, metaFile string, skipRemoveUnusedModels bool, excludePackages []string) error {
	var buf bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &buf)

	args := []string{"generate", "spec", "--scan-models", "--work-dir=" + workDir, "--output", outputFile}

	for _, excludePackage := range excludePackages {
		args = append(args, "--exclude="+excludePackage)
	}

	cmd := exec.Command("swagger", args...)
	cmd.Stdout = mw
	cmd.Stderr = mw
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swagger generate spec error: %v", err)
	}

	cmd = exec.Command("swagger", "mixin", outputFile, metaFile, "--output", outputFile)
	cmd.Stdout = mw
	cmd.Stderr = mw
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swagger mixin error: %v", err)
	}

	if !skipRemoveUnusedModels {
		cmd = exec.Command("swagger", "flatten", "--with-flatten", "remove-unused", "--output", outputFile, outputFile)
		cmd.Stdout = mw
		cmd.Stderr = mw
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("swagger flattern error: %v", err)
		}
	}

	return nil
}

func swaggerMeta(filepath string) error {
	swaggerMeta := MetaYaml{
		Swagger:  "2.0",
		BasePath: "/",
		Host:     "localhost",
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Schemes:  []string{"http"},
		MetaInfo: MetaInfo{
			Version:     gitVersion(),
			Title:       projectName,
			Description: projectDescription,
		},
	}

	data, err := yaml.Marshal(swaggerMeta)
	if err != nil {
		return fmt.Errorf("yaml marshal error: %v", err)
	}

	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("write file error: %v", err)
	}

	return nil
}

func gitVersion() string {
	args := []string{"describe", "--tags"}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "snapshot"
	}

	return strings.TrimRight(string(out), "\n")
}

type MetaYaml struct {
	Swagger  string   `yaml:"swagger"`
	BasePath string   `yaml:"basePath"`
	Host     string   `yaml:"host"`
	Consumes []string `yaml:"consumes"`
	Produces []string `yaml:"produces"`
	Schemes  []string `yaml:"schemes"`
	MetaInfo `yaml:"info"`
}

type MetaInfo struct {
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}
