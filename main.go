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
	transformChan := etl.TransformImagesV2(extractChan, "black_white")
	var wg sync.WaitGroup
	wg.Add(1)
	go etl.LoadFilesV2(transformChan, "black_white", &wg)
	go etl.LoadImageFiles("/home/calibraint/etl-image-go/black_white/", "images_bw.zip")
	wg.Wait()
}
