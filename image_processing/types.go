package image_processing

import (
	_ "fmt"
	"image"
	"image/color"
)

// Color channel return type that includes both the computed average color and the image that
// it used to compute it
type colorChanType struct {
	Color color.Color
	Image image.Image
}

// Interface for representing values needed to calculate the average color of an image
// and abstracting the functions to increment by a color and to get the average
type averageColorCalcType interface {
	R() uint32
	G() uint32
	B() uint32
	TotalCount() uint32
	Increment(color.Color)
	CalcAverage() color.Color
}

type averageColorCalcStruct struct {
	r           uint32
	g           uint32
	b           uint32
	total_count uint32
}

func (a averageColorCalcStruct) R() uint32 {
	return a.r
}

func (a averageColorCalcStruct) G() uint32 {
	return a.g
}

func (a averageColorCalcStruct) B() uint32 {
	return a.b
}

func (a averageColorCalcStruct) TotalCount() uint32 {
	return a.b
}

func (a *averageColorCalcStruct) Increment(color color.Color) {
	r, g, b, _ := color.RGBA()
	a.r += r
	a.g += g
	a.b += b
	a.total_count++
}

func (a averageColorCalcStruct) CalcAverage() color.Color {
	return color.RGBA{
		uint8(a.r / a.total_count),
		uint8(a.g / a.total_count),
		uint8(a.b / a.total_count),
		1,
	}
}
