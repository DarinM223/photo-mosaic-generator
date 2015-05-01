package image_processing

import (
	"fmt"
	"image"
	"image/color"
	"testing"
	"time"
)

func setUpImageMap() (map[color.Color]image.Image, color.Color, color.Color, color.Color, color.Color) {
	testMap := make(map[color.Color]image.Image)
	red := color.RGBA{
		255,
		0,
		0,
		1,
	}
	orange := color.RGBA{
		255,
		165,
		0,
		1,
	}
	blue := color.RGBA{
		0,
		0,
		255,
		1,
	}
	white := color.RGBA{
		255,
		255,
		255,
		1,
	}
	return testMap, red, orange, blue, white
}

func TestFindClosestColor(t *testing.T) {
	testMap, red, orange, blue, white := setUpImageMap()
	testMap[orange] = nil
	testMap[blue] = nil
	testMap[white] = nil
	c := FindClosestColor(red, testMap)
	if c != orange {
		t.Error("FindClosestColor should return orange")
	}
}

func TestCalculateAverageColor(t *testing.T) {
	images, err := OpenTestImages()
	if err != nil {
		t.Error(err)
	}
	ret := make(chan colorChanType)
	timeout := make(chan bool)
	time.Sleep(1 * time.Second)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	for _, img := range images {
		go CalculateAverageColor(img, ret)
	}

	go func() {
		for range images {
			select {
			case img := <-ret:
				fmt.Println(img.Color)
			case <-timeout:
				t.Error("Timed out!")
			}
		}
	}()
	time.Sleep(5 * time.Second)
}
