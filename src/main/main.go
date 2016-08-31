package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	dirname := "." + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			fmt.Println(fi.Name(), fi.Size(), "bytes")
			jpgMatched, _ := regexp.Compile(`\w+.jpeg|jpg|JPG|JPEG`)
			pngMatched, _ := regexp.Compile(`\w+.png|PNG`)
			if jpgMatched.MatchString(fi.Name()) {
				resizeJpg(fi.Name())
			} else if pngMatched.MatchString(fi.Name()) {
				resizePng(fi.Name())
			}
		}
	}

}

func resizeJpg(jpgFileName string) {
	file, err := os.Open(jpgFileName)
	if err != nil {
		log.Fatal(err)
	}
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(0, 150, img, resize.Lanczos3)
	var m1 image.Image
	if true {
		m1 = convertToGrayScale(m)
	}
	out, err := os.Create(jpgFileName + "_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// write new image to file
	jpeg.Encode(out, m1, nil)
}

func resizePng(pngFileName string) {
	file, err := os.Open(pngFileName)
	if err != nil {
		log.Fatal(err)
	}
	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(0, 150, img, resize.Lanczos3)
	var m1 image.Image
	if true {
		m1 = convertToGrayScale(m)
	}
	out, err := os.Create(pngFileName + "_resized.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// write new image to file
	png.Encode(out, m1)
}

func convertToGrayScale(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := MyGrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}
	return gray
}

func grayModel(c color.Color) color.Color {
	if _, ok := c.(color.Gray); ok {
		return c
	}
	r, g, b, _ := c.RGBA()

	//y := 0.2126*r + 0.7152*g + 0.0222*b
	//BT709 Greyscale: Red: 0.2125 Green: 0.7154 Blue: 0.0721
	//Y-Greyscale (YIQ/NTSC): Red: 0.299 Green: 0.587 Blue: 0.114
	//RMY Greyscale: Red: 0.5 Green: 0.419 Blue: 0.081
	y := (500*r + 419*g + 81*b + 500) / 1000
	return color.Gray{uint8(y >> 8)}
}

var MyGrayModel color.Model = color.ModelFunc(grayModel)
