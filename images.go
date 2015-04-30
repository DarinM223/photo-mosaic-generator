package main

import (
	"image"
	"image/color"
	"net/http"
)

func RetrieveImages(tagSearchResponse TagSearchResponse, c chan interface{}) {
	for _, value := range tagSearchResponse.Data {
		go func() {
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
		}()
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
func ProcessImages(imagesLength int, imageChan chan interface{}, ret chan map[color.Color]image.Image) {
	index := 0
	imageMap := make(map[color.Color]image.Image)
	colorChan := make(chan colorChanType)

	for {
		select {
		case img := <-imageChan: // finished receiving a new image
			switch img.(type) {
			case image.Image:
				go CalculateAverageColor(img.(image.Image), colorChan) // calculate the average color
			case error:
				// do nothing
				break
			}
			index++
			break
		case color := <-colorChan: // finished computing an average color
			if imageMap[color.Color] != nil {
				imageMap[color.Color] = color.Image // add the new image color pairing
			}
			if index >= imagesLength {
				ret <- imageMap
			}
			break
		}
	}
}
