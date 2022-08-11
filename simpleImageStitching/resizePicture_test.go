package simpleImageStitching

import (
	"fmt"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	stdWith = 274.0
	stdHigh = 316.0
)

func TestResize(t *testing.T) {
	tData, err := imgio.Open("./input/test.jpg")
	assert.NoError(t, err)
	std := 0.0
	if float64(tData.Bounds().Max.X)/stdWith > float64(tData.Bounds().Max.Y)/stdHigh {
		std = float64(tData.Bounds().Max.X) / stdWith
	} else {
		std = float64(tData.Bounds().Max.Y) / stdHigh
	}
	nData := transform.Resize(tData, int(float64(tData.Bounds().Max.X)/std), int(float64(tData.Bounds().Max.Y)/std), transform.Linear)
	if err := imgio.Save("./output/test.png", nData, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}
