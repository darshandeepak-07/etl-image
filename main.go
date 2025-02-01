package main

import (
	"fmt"
	"sync"

	"github.com/darshandeepak-07/etl-image-go/etl"
)

func main() {
	extractChan, err := etl.ExtractImageFilesV2("images/")
	if err != nil {
		fmt.Println("Error extracting files:", err)
		return
	}
	transformChannel := etl.TransformImagesV2(extractChan, "black_white")

	for task := range transformChannel {
		if task.Err != nil {
			fmt.Println("Error processing: ", task.Err)
		} else {
			fmt.Println("Processed : ", task.Path)
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go etl.LoadImagesV2("/home/calibraint/etl-image-go/black_white/", "images_bw.zip", &wg)
	wg.Wait()
}
