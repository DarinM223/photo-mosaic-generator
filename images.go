package main

import (
	"image"
	"image/color"
	"net/http"
)

func RetrieveImages(tagSearchResponse TagSearchResponse, c chan<- interface{}) {
	go func() {
		for _, value := range tagSearchResponse.Data {
			resp, err := http.Get(value.Images.Thumbnail.Url)
			defer resp.Body.Close()
			if err != nil {
				c <- err
				return
			}

			img, _, err := image.Decode(resp.Body)
			if err != nil {
				c <- err
			} else {
				c <- img
			}
		}
	}()
}

func CalculateAverageColor(img image.Image) color.Color {
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

	return color.RGBA{
		uint8(r_count / total_count),
		uint8(b_count / total_count),
		uint8(g_count / total_count),
		1,
	}
}

/// receive images as they are sent and calculate average color for each image
func ProcessImages(imagesLength int, c <-chan interface{}, retChan chan<- map[color.Color]image.Image) {
	index := 0
	imageMap := make(map[color.Color]image.Image)
	for {
		img := <-c
		switch img.(type) {
		case image.Image:
			// calculate the average color
			c := CalculateAverageColor(img.(image.Image))
			if imageMap[c] != nil {
				imageMap[c] = img.(image.Image)
			}
		case error:
			break
		}

		index++
		// send return map after processing all of the images
		if index >= imagesLength {
			retChan <- imageMap
			break
		}
	}
}
