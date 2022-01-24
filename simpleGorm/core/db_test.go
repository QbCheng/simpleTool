package core

import (
	"github.com/stretchr/testify/assert"
	"simpleTool/simpleGorm/models"
	"testing"
)

func TestMigrator(t *testing.T) {
	err := Migrator(models.User{})
	assert.NoError(t, err)
}
