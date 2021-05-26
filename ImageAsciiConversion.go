package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/nfnt/resize"
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

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y

	/*if width > height {
		if width > 150 {
			img = resize.Resize(150, 0, img, resize.Lanczos3)
		}
	} else if height > width {
		if height > 100 {
			img = resize.Resize(0, 100, img, resize.Lanczos3)
		}
	} else {
		if width > 150 {
			img = resize.Resize(150, 0, img, resize.Lanczos3)
		} else if height > 100 {
			img = resize.Resize(0, 100, img, resize.Lanczos3)
		}
	}*/
	img = resize.Resize(150, 0, img, resize.Lanczos3)
	//img = resize.Thumbnail(150, 100, img, resize.Lanczos3)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixel := (r + g + b) / 3
			text += grayscale[mapRange(int(pixel), 0, 65535, 0, 7)]
		}
		text += "\n"
	}
	return
}
