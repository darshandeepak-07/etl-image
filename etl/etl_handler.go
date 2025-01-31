package etl

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"github.com/darshandeepak-07/etl-image-go/utils"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExtractImageFiles(directory string) []string {
	files, error := os.ReadDir(directory)

	if error != nil {
		log.Fatal(error)
	}
	check(error)
	var images []string

	for _, value := range files {
		if !value.IsDir() && (filepath.Ext(value.Name()) == ".jpeg" || filepath.Ext(value.Name()) == ".png") {
			images = append(images, filepath.Join(directory, value.Name()))
		}
	}

	return images
}

func TransformImageFiles(imagePaths []string, outputDirectory string) {

	err := os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	for _, path := range imagePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", path, err)
			continue
		}

		decodedImage, format, error := image.Decode(file)
		defer file.Close()
		if error != nil {
			fmt.Println("Error decoding image:", path, err)
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

		outputFilePath := filepath.Join(outputDirectory, filepath.Base(path))

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

		defer outputFile.Close()

		if err != nil {
			fmt.Println("Error encoding image:", outputFilePath, err)
			continue
		}

		fmt.Println("Grayscale image saved as", outputFilePath)
	}
}

func LoadImageFiles(outputDirectory, zipFileName string) error {
	return utils.ZipFolder(outputDirectory, zipFileName)
}
