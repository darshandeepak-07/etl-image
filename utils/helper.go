package utils

import (
	"archive/zip"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func ResizeImage(imageFile *image.Gray) image.Image {
	return resize.Resize(1920, 1080, imageFile, resize.Lanczos2)
}

func ZipFolder(sourceFolder, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)

	if err != nil {
		return err
	}
	defer removeDirectory(filepath.Base(sourceFolder))
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(sourceFolder, path)

		if err != nil {
			return err
		}

		zipHeader, _ := zip.FileInfoHeader(info)
		zipHeader.Name = relativePath
		zipHeader.Method = zip.Deflate

		if info.IsDir() {
			zipHeader.Name += "/"
			_, err := zipWriter.CreateHeader(zipHeader)
			return err
		}

		zipWriterEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(zipWriterEntry, file)
		return err
	})
	return err
}

func removeDirectory(directory string) error {
	err := os.RemoveAll(directory)

	if err == nil {
		log.Printf("Removed directory : %s\n", directory)
	} else {
		log.Printf("Failed to remove directory : %s\n", directory)
	}
	return err
}
