package etl

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/darshandeepak-07/etl-image-go/utils"
)

type Task struct {
	Path  string
	Image image.Image
	Err   error
}

func ExtractImage(directory string) (<-chan Task, error) {
	out := make(chan Task)
	go func() {
		defer close(out)
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".png") {
				out <- Task{Path: path}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking through directory: %v\n", err)
		}
	}()
	return out, nil
}

func TransforImage(inputChannel <-chan Task, outputDirectory string) <-chan Task {
	err := os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		log.Printf("Error creating output directory: %v\n", err)
	}
	outputChannel := make(chan Task)

	go func() {
		defer close(outputChannel)
		for task := range inputChannel {
			file, err := os.Open(task.Path)
			if err != nil {
				log.Printf("Error opening file: %s %v\n", task.Path, err)
				continue
			}
			defer file.Close()

			decodedImage, format, decodeErr := image.Decode(file)
			if decodeErr != nil {
				log.Printf("Error decoding image: %s %v\n", task.Path, decodeErr)
				continue
			}

			bounds := decodedImage.Bounds()
			grayImage := image.NewGray(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					originalColor := decodedImage.At(x, y)
					grayColor := color.GrayModel.Convert(originalColor)
					grayImage.Set(x, y, grayColor)
				}
			}

			outputFilePath := filepath.Join(outputDirectory, filepath.Base(task.Path))

			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				log.Printf("Error creating output file: %s %v\n", outputFilePath, err)
				continue
			}

			resizedImage := utils.ResizeImage(grayImage)
			switch format {
			case "jpeg":
				err = jpeg.Encode(outputFile, resizedImage, nil)
			case "png":
				err = png.Encode(outputFile, resizedImage)
			default:
				log.Printf("Unsupported image format: %s", format)
				outputFile.Close()
				continue
			}

			if err != nil {
				log.Printf("Error encoding image: %s %v\n", outputFilePath, err)
				outputFile.Close()
				continue
			}

			outputFile.Close()
			log.Printf("Grayscale image saved as %s\n", outputFilePath)
			task.Image = resizedImage
			task.Err = err
			outputChannel <- task
		}
	}()
	return outputChannel
}

func LoadImages(source, zipFileName string, wg *sync.WaitGroup) error {
	defer wg.Done()
	return utils.ZipFolder(source, zipFileName)
}
