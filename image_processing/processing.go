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

type colorChanType struct {
	Color color.Color
	Image image.Image
}

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
func GetBestImage(c color.Color, imageMap map[color.Color]image.Image,
	savedBestColorMap map[color.Color]color.Color) image.Image {

	if imageMap[c] != nil { // if there is an exact color match
		return imageMap[c]
	} else if savedBestColorMap[c] != nil { // if the closest color has already been calculated
		return imageMap[savedBestColorMap[c]]
	} else {
		closestColor := FindClosestColor(c, imageMap)
		savedBestColorMap[c] = closestColor // save best color mapping
		return imageMap[closestColor]
	}
}

// Calculates the average color of an image and sends the result through a channel
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
