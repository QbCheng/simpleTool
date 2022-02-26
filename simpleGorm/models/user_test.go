package models

import (
	"github.com/QbCheng/simpleTool/simpleGorm/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestMock(t *testing.T) {
	t.Log(Mock())
}

func TestCreateData(t *testing.T) {
	err := core.GetDB().Transaction(func(tx *gorm.DB) error {
		for i := 0; i < 100; i++ {
			if err := tx.Create(Mock()).Error; err != nil {
				return err
			}
		}
		return nil
	})
	assert.NoError(t, err)
}
