package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	//"github.com/nfnt/resize"
)

var grayscale [8]string = [8]string{"@", "%", "#", "*", "+", "=", ":", "."}

func mapRange(x, in_min, in_max, out_min, out_max int) int {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}

func ImageToAscii(r io.Reader) (text string, err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return
	}

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixel := (r + g + b) / 3
			text += grayscale[mapRange(int(pixel), 0, 65535, 0, 7)]
		}
		text += "\n"
	}
	return
}
