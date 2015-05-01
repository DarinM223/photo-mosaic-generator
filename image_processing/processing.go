package image_processing

import (
	pkgColor "code.google.com/p/sadbox/color"
	"image"
	"image/color"
	"math"
)

const (
	WEIGHT_HUE        = 0.8
	WEIGHT_SATURATION = 0.1
	WEIGHT_VALUE      = 0.1
)

func convertRGBToHSV(c color.Color) pkgColor.HSV {
	r, g, b, _ := c.RGBA()
	h, s, v := pkgColor.RGBToHSV(uint8(r), uint8(g), uint8(b))
	return pkgColor.HSV{
		h,
		s,
		v,
	}
}

// Finds the color closest to the specific color and returns it
func FindClosestColor(c color.Color, imageMap map[color.Color]image.Image) color.Color {
	// otherwise, find the closest image that matches the color
	minDistance := math.MaxFloat64
	var minColor color.Color = color.White
	cHSV := convertRGBToHSV(c)
	for color, _ := range imageMap {
		imgHSV := convertRGBToHSV(color)
		dH := float64(imgHSV.H) - float64(cHSV.H)
		dS := float64(imgHSV.S) - float64(cHSV.S)
		dV := float64(imgHSV.V) - float64(cHSV.V)

		distance := math.Sqrt(WEIGHT_HUE*math.Pow(dH, 2) +
			WEIGHT_SATURATION*math.Pow(dS, 2) +
			WEIGHT_VALUE*math.Pow(dV, 2))

		if distance < minDistance {
			minDistance = distance
			minColor = color
		}
	}
	return minColor
}

// Gets the best image for the color
// if image is not in the color map, then create a new goroutine to calculate the closest image
// and sent it through a return channel
func GetBestImage(c color.Color, imageMap map[color.Color]image.Image) image.Image {
	if imageMap[c] != nil { // if there is an exact color match
		return imageMap[c]
	} else {
		closestColor := FindClosestColor(c, imageMap)
		return imageMap[closestColor]
	}
}

// Calculates the average color of an image and sends the result through a channel
func CalculateAverageColor(img image.Image, ret chan colorChanType) {
	var averageCalc averageColorCalcType = &averageColorCalcStruct{
		0,
		0,
		0,
		0,
	}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			averageCalc.Increment(img.At(x, y))
		}
	}

	ret <- colorChanType{
		averageCalc.CalcAverage(),
		img,
	}
}

// Breaks the image into multiple regions and returns through a channel a map of region index to the calculated
// average color of each region. The region index is specified by a point starting from (0, 0)
// to (number of x regions - 1, number of y regions - 1)
func CalculateAverageColorRegions(img image.Image, region image.Rectangle, ret chan map[image.Point]color.Color) {
	regionMap := make(map[image.Point]averageColorCalcType)
	averageColorMap := make(map[image.Point]color.Color)

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			currRegionX := int(float64(x) / float64(region.Dx()))
			currRegionY := int(float64(y) / float64(region.Dy()))

			currRegion := image.Point{
				currRegionX,
				currRegionY,
			}
			if regionMap[currRegion] == nil {
				regionMap[currRegion] = &averageColorCalcStruct{
					0,
					0,
					0,
					0,
				}
			} else {
				regionMap[currRegion].Increment(img.At(x, y))
			}
		}
	}

	for r, calc := range regionMap {
		averageColorMap[r] = calc.CalcAverage()
	}
	ret <- averageColorMap
}
