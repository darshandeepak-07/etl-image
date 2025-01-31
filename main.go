package main

import (
	"fmt"

	"github.com/darshandeepak-07/etl-image-go/etl"
)

func main() {
	fmt.Println("---------------------------------")
	fmt.Println("ETL Pipeline for Image Processing")
	fmt.Println("---------------------------------")
	fmt.Println("---------------------------------")

	fmt.Println("Files Found")
	fmt.Println("---------------------------------")
	files := etl.ExtractImageFiles("images/")
	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println("---------------------------------")

	etl.TransformImageFiles(files, "black_white")

	err := etl.LoadImageFiles("/home/calibraint/etl-image-go/black_white/", "images_bw.zip")
	if err != nil {
		fmt.Println("Image files zipped successfully")
	}
}
