package simpleImageStitching

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDecorateNeedful(t *testing.T) {
	decorate, err := NewDecorate("./input", true)
	assert.NoError(t, err)
	data := decorate.Data()
	count := 0
	for i := range data {
		if data[i] == nil {
			continue
		}
		count++
	}
	assert.Equal(t, count, len(data))
}

func TestNewDecorateOptional(t *testing.T) {
	decorate, err := NewDecorate("./input", false)
	assert.NoError(t, err)
	data := decorate.Data()
	count := 0
	for i := range data {
		if data[i] == nil {
			continue
		}
		count++
	}
	assert.Equal(t, count+1, len(data))
}

func TestDecorate_Name(t *testing.T) {
	assert.Equal(t, name("/aaa/aaa/bbb"), "bbb")
	assert.Equal(t, name("/aaa/aaa/bbb/"), "bbb")
	assert.Equal(t, name("F:\\selfSvn\\grammar\\go\\tool\\simpleTool\\simpleBlendPicture\\input"), "input")
	assert.Equal(t, name("F:\\selfSvn\\grammar\\go\\tool\\simpleTool\\simpleBlendPicture\\input\\"), "input")
}
