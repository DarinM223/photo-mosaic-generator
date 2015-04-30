package image_processing

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"testing"
	"time"
)

func OpenTestImages() ([]image.Image, error) {
	var images []image.Image
	f_img1, err := os.Open("test/images/gilg.jpg")
	if err != nil {
		return nil, err
	}
	f_img2, err := os.Open("test/images/illya.jpg")
	if err != nil {
		return nil, err
	}
	f_img3, err := os.Open("test/images/jesus_tatsuya.jpg")
	if err != nil {
		return nil, err
	}
	f_img4, err := os.Open("test/images/aldnoah.jpg")
	if err != nil {
		return nil, err
	}

	defer f_img1.Close()
	defer f_img2.Close()
	defer f_img3.Close()
	defer f_img4.Close()

	img1, _, err := image.Decode(f_img1)
	if err != nil {
		return nil, err
	}
	img2, _, err := image.Decode(f_img2)
	if err != nil {
		return nil, err
	}
	img3, _, err := image.Decode(f_img3)
	if err != nil {
		return nil, err
	}
	img4, _, err := image.Decode(f_img4)
	if err != nil {
		return nil, err
	}

	images = append(images, img1)
	images = append(images, img2)
	images = append(images, img3)
	images = append(images, img4)
	return images, nil
}

// should build a map of size 4 of colors to images
func shouldProperlyBuildMap(t *testing.T, images []image.Image) {
	imageChan := make(chan image.Image)
	retChan := make(chan map[color.Color]image.Image)
	go ProcessImages(4, imageChan, retChan)

	for _, image := range images {
		imageChan <- image
	}

	timeout := make(chan bool)

	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()

	select {
	case colorMap := <-retChan:
		fmt.Println(colorMap)
	case <-timeout:
		t.Error("ProcessImages timed out!")
	}
}

// should time out after sending only three images but after the timeout is hit it should
// still send the map, which will have only three key/value pairs
func shouldTimeOutAfterSendingThree(t *testing.T, images []image.Image) {
	imageChan := make(chan image.Image)
	retChan := make(chan map[color.Color]image.Image)
	go ProcessImages(4, imageChan, retChan)

	for i := 0; i < len(images)-1; i++ {
		imageChan <- images[i]
	}

	timeout := make(chan bool)

	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()

	select {
	case colorMap := <-retChan:
		if len(colorMap) != 3 {
			t.Error("colorMap's length should be 3")
		}
		fmt.Println(colorMap)
	case <-timeout:
		t.Error("ProcessImages timed out!")
	}
}

func TestProcessImages(t *testing.T) {
	images, err := OpenTestImages()
	if err != nil {
		t.Error(err.Error())
		return
	}

	shouldProperlyBuildMap(t, images)
	shouldTimeOutAfterSendingThree(t, images)
}
