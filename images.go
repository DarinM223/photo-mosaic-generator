package main

import (
	"image"
	"image/color"
	"net/http"
)

func RetrieveImages(tagSearchResponse TagSearchResponse, c chan<- *image.Image) {
	go func() {
		for _, value := range tagSearchResponse.Data {
			resp, err := http.Get(value.Images.Thumbnail.Url)
			defer resp.Body.Close()
			if err != nil {
				c <- nil
			} else {
				img, _, err := image.Decode(resp.Body)
				if err != nil {
					c <- &img
				} else {
					c <- nil
				}
			}

		}
	}()
}

// TODO: implement this
func CalculateAverageColor(img *image.Image) color.Color {
	return (*img).At(0, 0)
}

/// receive images as they are sent and calculate average color for each image
func ProcessImages(imagesLength int, c <-chan *image.Image, retChan chan<- map[color.Color]*image.Image) {
	index := 0
	imageMap := make(map[color.Color]*image.Image)
	for {
		img := <-c
		if img != nil {
			// calculate the average color
			c := CalculateAverageColor(img)
			if imageMap[c] != nil {
				imageMap[c] = img
			}
		}
		index++
		// send return map after processing all of the images
		if index >= imagesLength {
			retChan <- imageMap
			break
		}
	}
}
