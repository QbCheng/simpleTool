package simpleImageStitching

import (
	"errors"
	"fmt"
	"github.com/QbCheng/simpleTool/simpleShared"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"image"
	"strings"
)

const defaultOutputDir = "./"
const defaultFilterDir = "./"

type Option func(log *BlendPicture)

// WithOutputDir 输出路径
func WithOutputDir(outputDir string) Option {
	return func(log *BlendPicture) {
		if !strings.HasSuffix(outputDir, "/") {
			outputDir += "/"
		}
		log.option.outputDir = outputDir
	}
}

// WithFilterDir 过滤文件地址
func WithFilterDir(filterDir string) Option {
	return func(log *BlendPicture) {
		if !strings.HasSuffix(filterDir, "/") {
			filterDir += "/"
		}
		log.option.filterDir = filterDir
	}
}

/*
BlendPicture 混合图片
NOTE : 图片尺寸保持一致. 保证位置在图片贴合后符合要求
*/
type BlendPicture struct {
	optionName string
	data       map[uint]map[string]*Decorate

	option *blendPictureOption
}

type blendPictureOption struct {
	outputDir string
	filterDir string
}

type Parameter struct {
	Dir     string
	Needful bool
	Level   uint
}

func NewBlendPicture(parameter []Parameter, options ...Option) (*BlendPicture, error) {
	ret := &BlendPicture{
		option: &blendPictureOption{
			outputDir: defaultOutputDir,
			filterDir: defaultFilterDir,
		},
		data: map[uint]map[string]*Decorate{},
	}

	for i := range options {
		options[i](ret)
	}

	for _, v := range parameter {
		var d *Decorate
		var err error
		if v.Needful {
			d, err = NewDecorate(v.Dir, true)
		} else {
			d, err = NewDecorate(v.Dir, false)
		}
		if err != nil {
			return nil, err
		}

		if _, ok := ret.data[v.Level]; !ok {
			ret.data[v.Level] = make(map[string]*Decorate)
		}
		data := ret.data[v.Level]
		if _, ok := data[d.Name()]; ok {
			return nil, errors.New(" Reloading directories.")
		}
		data[d.Name()] = d
	}

	// 加载过滤条件
	LoadConditions(ret.option.filterDir)
	return ret, nil
}

func (s BlendPicture) SlowRun(maxLevel uint) {
	var count = 0
	var limit = 10
	var t []string

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>> 预处理开始 >>>>>>>>>>>>>>>>>>>>>>>")
	for i := maxLevel; i > 0; i-- {
		for j := range s.data[i] {
			if count == 0 {
				t = s.data[i][j].AbsolutePath()
			} else {
				t = pathRecursion(t, s.data[i][j].AbsolutePath())
			}
			if count%limit == 0 {
				fmt.Println(" >>>>>>>>>>>>>>>>>>>>>>> 预处理中 >>>>>>>>>>>>>>>>>>>>>>>")
			}
			count++
		}
	}

	filter := 0
	for i := range t {
		if Filter(t[i]) {
			filter++
		}
	}

	_, err := SaveFile(s.option.outputDir+"paths.pre", t)
	if err != nil {
		panic(err)
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>> 开始输出图片 >>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> 图片总数: %v >>>>>>>>>>>>>>>>>>>>>>>\n", len(t))
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> 过滤数量: %v >>>>>>>>>>>>>>>>>>>>>>>\n", filter)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> 实际数量: %v >>>>>>>>>>>>>>>>>>>>>>>\n", len(t)-filter)
	pictureCount := 0
	for _, paths := range t {
		if Filter(paths) {
			continue
		}
		var outputData image.Image
		pathSlice := strings.Split(paths, "||")
		for _, path := range pathSlice {
			tData, err := imgio.Open(path)
			if err != nil {
				return
			}
			if outputData == nil {
				outputData = tData
			} else {
				outputData = blend.Normal(outputData, tData)
			}
		}
		s.outputPng(outputData)
		pictureCount++
		fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> 进度 %v : %v >>>>>>>>>>>>>>>>>>>>>>> \n", pictureCount, len(t)-filter)
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>> 完成 >>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> 总共输出图片 %v 张 >>>>>>>>>>>>>>>>>>>>>>> \n", len(t)-filter)
}

func pathRecursion(one []string, two []string) (ret []string) {
	for i := range one {
		for j := range two {
			if one[i] == "" && two[j] == "" {
				ret = append(ret, "")
			} else if one[i] == "" {
				ret = append(ret, two[j])
			} else if two[j] == "" {
				ret = append(ret, one[i])
			} else {
				ret = append(ret, one[i]+`||`+two[j])
			}
		}
	}
	return
}

func (s BlendPicture) outputPng(image image.Image) {
	output := s.option.outputDir + simpleShared.GetRandomString(16) + ".png"
	if err := imgio.Save(output, image, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}
