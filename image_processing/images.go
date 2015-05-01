package image_processing

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"
	"time"
)

// retrieves images from urls and passes it into channel, resizing the image to fit a rectangle
func RetrieveImages(imageUrls []string, c chan image.Image, regionSize image.Rectangle) {
	client := http.Client{
		Timeout: 500,
	}
	for _, value := range imageUrls {
		resp, _ := client.Get(value)
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		resizedImg := resize.Resize(uint(regionSize.Dx()), uint(regionSize.Dy()), img, resize.NearestNeighbor)
		if err == nil {
			c <- resizedImg
		}
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
