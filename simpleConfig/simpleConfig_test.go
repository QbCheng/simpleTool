package simpleConfig

import (
	"github.com/stretchr/testify/assert"
	"simpleTool/simpleConfig/model"
	"testing"
)

func TestCreateAndInit(t *testing.T) {
	sc, err := CreateAndInit("app", "yaml", "./config/yaml/")
	assert.NoError(t, err)
	t.Log(sc.AllSettings())
	t.Log(sc.AllKeys())
}

/*
TestUnmarshal 将配置 绑定 到 指定对象.
*/
func TestUnmarshal(t *testing.T) {
	sc, err := CreateAndInit("app", "yaml", "./config/yaml/")
	assert.NoError(t, err)
	config := &model.App{}
	t.Log(sc.AllSettings())
	err = sc.Unmarshal(config)
	assert.NoError(t, err)
	t.Log(config)
}

func TestEnv(t *testing.T) {
	sc, err := CreateAndInit("app", "yaml", "./config/yaml/")
	assert.NoError(t, err)
	sc.AutomaticEnv()
	t.Log(sc.Get("chenTest"))
	t.Log(sc.AllSettings())
	t.Log(sc.AllKeys())
}
