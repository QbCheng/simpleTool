package simpleImageStitching

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	LoadConditions("./filter/")
	data := "./input/12/2.png||./input/11/1.png||./input/10/1.png||./input/9/1.png||./input/8/1.png||./input/7/1.png||./input/6/1.png||./input/5/1.png||./input/4/1.png||./input/3/2.png||./input/2/1.png||./input/1/3.png"
	if Filter(data) {
		fmt.Println("过滤")
	} else {
		fmt.Println("不过滤")
	}
}

func TestLoadConditions(t *testing.T) {
	LoadConditions("./filter/")
}
