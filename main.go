package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/darshandeepak-07/etl-image-go/etl"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Dir struct {
		InputDirectory  string `yaml:"input"`
		OutputDirectory string `yaml:"output"`
		ZipFile         string `yaml:"zipFile"`
	} `yaml:"dir"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	outputDirectory := config.Dir.OutputDirectory
	log.Printf("File : %s %s %s", outputDirectory, config.Dir.InputDirectory, config.Dir.ZipFile)
	defer cancel()
	extractChan, err := etl.ExtractImage(config.Dir.InputDirectory)
	if err != nil {
		log.Printf("Error extracting files: %v\n", err)
		return
	}
	transformChannel := etl.TransforImage(extractChan, outputDirectory)

	for task := range transformChannel {
		if task.Err != nil {
			log.Printf("Error processing: %v\n", task.Err)
		} else {
			log.Printf("Processed : %s\n", task.Path)
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go etl.LoadImages(outputDirectory, config.Dir.ZipFile, &wg)
	wg.Wait()
	ctx.Done()
}
