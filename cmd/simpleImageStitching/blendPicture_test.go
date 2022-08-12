package simpleImageStitching

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlendPicture_SlowRun2(t *testing.T) {
	b, err := NewBlendPicture([]Parameter{
		{
			Dir:     "./input/12/",
			Needful: true,
			Level:   12,
		},
		{
			Dir:     "./input/11/",
			Needful: true,
			Level:   11,
		},
		{
			Dir:     "./input/10/",
			Needful: true,
			Level:   10,
		},
		{
			Dir:     "./input/9/",
			Needful: true,
			Level:   9,
		},
		{
			Dir:     "./input/8/",
			Needful: true,
			Level:   8,
		},
		{
			Dir:     "./input/7/",
			Needful: true,
			Level:   7,
		},

		{
			Dir:     "./input/6/",
			Needful: true,
			Level:   6,
		},
		{
			Dir:     "./input/5/",
			Needful: true,
			Level:   5,
		},
		{
			Dir:     "./input/4/",
			Needful: true,
			Level:   4,
		},
		{
			Dir:     "./input/3/",
			Needful: true,
			Level:   3,
		},
		{
			Dir:     "./input/2/",
			Needful: true,
			Level:   2,
		},
		{
			Dir:     "./input/1/",
			Needful: true,
			Level:   1,
		},
	},
		WithOutputDir("./output"),
		WithFilterDir("./filter"),
	)
	assert.NoError(t, err)
	b.SlowRun(12)
}
