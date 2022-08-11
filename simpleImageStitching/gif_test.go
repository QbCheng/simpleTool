package simpleImageStitching

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateGif(t *testing.T) {
	err := GenerateGif()
	assert.NoError(t, err)
}
