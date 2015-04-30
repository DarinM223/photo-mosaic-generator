package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"
	"time"
)

func RetrieveImages(tagSearchResponse TagSearchResponse, c chan image.Image) {
	client := http.Client{
		Timeout: 500,
	}
	for _, value := range tagSearchResponse.Data {
		resp, _ := client.Get(value.Images.Thumbnail.Url)
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err == nil {
			c <- img
		}
	}
}

type colorChanType struct {
	Color color.Color
	Image image.Image
}

func CalculateAverageColor(img image.Image, ret chan colorChanType) {
	r_count, b_count, g_count := uint32(0), uint32(0), uint32(0)
	total_count := uint32(0)

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, b, g, _ := img.At(x, y).RGBA()

			r_count += r
			b_count += b
			g_count += g
			total_count++
		}
	}

	ret <- colorChanType{
		color.RGBA{
			uint8(r_count / total_count),
			uint8(b_count / total_count),
			uint8(g_count / total_count),
			1,
		},
		img,
	}
}

/// receive images as they are sent and calculate average color for each image
func ProcessImages(imagesLength int, imageChan chan image.Image, ret chan map[color.Color]image.Image) {
	index := 0
	imageMap := make(map[color.Color]image.Image)
	averageColorChan := make(chan colorChanType)
	timeout := make(chan bool)

	// timeout after a second
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	for {
		select {
		case img := <-imageChan: // finished receiving a new image
			go CalculateAverageColor(img.(image.Image), averageColorChan)
		case color := <-averageColorChan: // finished computing an average color
			if imageMap[color.Color] == nil {
				imageMap[color.Color] = color.Image // add the new image color pairing
			}
			index++
			if index >= imagesLength {
				ret <- imageMap
				return
			}
		case <-timeout:
			fmt.Println("Image processing timed out!")
			ret <- imageMap
			return
		}
	}
}
