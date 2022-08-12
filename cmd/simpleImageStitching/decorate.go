package simpleImageStitching

import (
	"github.com/anthonynsimon/bild/imgio"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Decorate struct {
	loadDir      string        // 加载路径
	data         []image.Image // 图片数据
	absolutePath []string      // 图片路径
	needful      bool          // 必需的
}

func NewDecorate(dir string, needful bool) (*Decorate, error) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	ret := &Decorate{
		loadDir: dir,
		needful: needful,
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for i := range files {
		data, err := imgio.Open(ret.loadDir + files[i].Name())
		if err != nil {
			return nil, err
		}
		ret.data = append(ret.data, data)
		ret.absolutePath = append(ret.absolutePath, ret.loadDir+files[i].Name())
	}

	return ret, nil
}

// Name 当前 Decorate 的名字
func (d Decorate) Name() string {
	t := strings.TrimSuffix(d.loadDir, `/`)
	t = strings.TrimSuffix(t, `\`)
	_, f := filepath.Split(t)
	return f
}

func name(s string) string {
	t := strings.TrimSuffix(s, `/`)
	t = strings.TrimSuffix(t, `\`)
	_, f := filepath.Split(t)
	return f
}

func (d Decorate) Data() []image.Image {
	var ret []image.Image
	if d.needful {
		ret = d.data
	} else {
		ret = append(ret, nil)
		ret = append(ret, d.data...)
	}
	return ret
}

func (d Decorate) AbsolutePath() []string {
	var ret []string
	if d.needful {
		ret = d.absolutePath
	} else {
		ret = append(ret, "")
		ret = append(ret, d.absolutePath...)
	}
	return ret
}

func (d Decorate) NeedFul() bool {
	return d.needful
}
