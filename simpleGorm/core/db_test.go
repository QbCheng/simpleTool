package core

import (
	"github.com/QbCheng/simpleTool/simpleGorm/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigrator(t *testing.T) {
	err := Migrator(models.User{})
	assert.NoError(t, err)
}
