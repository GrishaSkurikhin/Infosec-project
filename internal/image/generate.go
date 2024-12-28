package image

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

const (
	squareSize = 20
)

func GenerateRandomImage(width, height int) image.Image {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += squareSize {
		for x := 0; x < width; x += squareSize {
			r := uint8(random.Intn(256))
			g := uint8(random.Intn(256))
			b := uint8(random.Intn(256))
			a := uint8(255)
			for yy := y; yy < y+squareSize && yy < height; yy++ {
				for xx := x; xx < x+squareSize && xx < width; xx++ {
					img.Set(xx, yy, color.RGBA{R: r, G: g, B: b, A: a})
				}
			}
		}
	}
	return img
}
