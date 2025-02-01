package etl

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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

func ExtractImageFilesV2(directory string) (<-chan Task, error) {
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
			fmt.Println("Error walking through directory:", err)
		}
	}()

	return out, nil
}

func TransformImagesV2(inputChannel <-chan Task, outputDirectory string) <-chan Task {
	err := os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
	}
	outputChannel := make(chan Task)

	go func() {
		defer close(outputChannel)
		for task := range inputChannel {
			file, err := os.Open(task.Path)
			if err != nil {
				fmt.Println("Error opening file:", task.Path, err)
				continue
			}
			defer file.Close()

			decodedImage, format, decodeErr := image.Decode(file)
			if decodeErr != nil {
				fmt.Println("Error decoding image:", task.Path, decodeErr)
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
				fmt.Println("Error creating output file:", outputFilePath, err)
				continue
			}

			resizedImage := utils.ResizeImage(grayImage)
			switch format {
			case "jpeg":
				err = jpeg.Encode(outputFile, resizedImage, nil)
			case "png":
				err = png.Encode(outputFile, resizedImage)
			default:
				fmt.Println("Unsupported image format:", format)
				outputFile.Close()
				continue
			}

			if err != nil {
				fmt.Println("Error encoding image:", outputFilePath, err)
				outputFile.Close()
				continue
			}

			outputFile.Close()
			fmt.Println("Grayscale image saved as", outputFilePath)
			task.Image = grayImage
			task.Err = err
			outputChannel <- task
		}
	}()
	return outputChannel
}

func LoadFilesV2(in <-chan Task, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range in {
		if task.Err != nil {
			fmt.Println(task.Err)
			continue
		}

		outputFilePath := filepath.Join(outputDir, filepath.Base(task.Path))
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Printf("Error creating output file %s: %v\n", outputFilePath, err)
			continue
		}
		switch filepath.Ext(task.Path) {
		case ".jpg", ".jpeg":
			err = jpeg.Encode(outputFile, task.Image, nil)
		case ".png":
			err = png.Encode(outputFile, task.Image)
		default:
			fmt.Printf("Unsupported image format for file %s\n", task.Path)
		}

		outputFile.Close()
		if err != nil {
			fmt.Printf("Error encoding image %s: %v\n", outputFilePath, err)
			continue
		}
		fmt.Println("Processed and saved:", outputFilePath)
	}
}
