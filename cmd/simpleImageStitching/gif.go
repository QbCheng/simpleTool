package simpleImageStitching

import (
	"github.com/anthonynsimon/bild/imgio"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

func GenerateGif() error {
	files := []string{"./input/1/1.png", "./input/1/2.png", "./input/1/3.png"}

	output := "./test.gif"
	delay := 1
	anim := gif.GIF{}
	for _, file := range files {
		img, err := imgio.Open(file)
		if err != nil {
			return err
		}
		paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.Point{})
		anim.Image = append(anim.Image, paletted)
		anim.Delay = append(anim.Delay, delay)
	}
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	err = gif.EncodeAll(f, &anim)
	if err != nil {
		return err
	}
	return nil
}
