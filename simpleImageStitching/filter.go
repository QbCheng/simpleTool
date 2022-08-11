package simpleImageStitching

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var Config [][]Data

var Conditions []Condition

func LoadConditions(dir string) {
	data, err := os.Open(dir + "filter.json")
	if err != nil {
		panic(err)
	}
	defer data.Close()
	byteData, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteData, &Config)
	if err != nil {
		panic(err)
	}

	for _, v1 := range Config {
		cond := NewCondition()
		for _, v2 := range v1 {
			cond.data[v2] = struct{}{}
		}
		Conditions = append(Conditions, cond)
	}
	fmt.Println(Conditions)
}

type Data struct {
	Level int    `json:"level"`
	Name  string `json:"name"`
}

type Condition struct {
	data map[Data]interface{}
}

func NewCondition() Condition {
	return Condition{
		data: map[Data]interface{}{},
	}
}

// ./input/12/2.png||./input/11/1.png||./input/10/1.png||./input/9/1.png||./input/8/1.png||./input/7/1.png||./input/6/1.png||./input/5/1.png||./input/4/1.png||./input/3/2.png||./input/2/1.png||./input/1/3.png"
func (c Condition) Reach(data []Data) bool {
	reach := 0
	for _, v := range data {
		for cdk, _ := range c.data {
			if cdk == v {
				reach++
			}
		}
	}
	return reach >= len(c.data)
}

// ./input/12/2.png||./input/11/1.png||./input/10/1.png||./input/9/1.png||./input/8/1.png||./input/7/1.png||./input/6/1.png||./input/5/1.png||./input/4/1.png||./input/3/2.png||./input/2/1.png||./input/1/3.png"
func Filter(prepData string) bool {
	d1 := strings.Split(prepData, "||")
	var data []Data
	for _, path := range d1 {
		d2 := strings.Split(path, "/")
		if len(d2) < 3 {
			panic(fmt.Sprintf("Filter, image err. %v", path))
		}
		level, err := strconv.Atoi(d2[len(d2)-2])
		if err != nil {
			panic(err)
		}
		data = append(data, Data{
			Level: level,
			Name:  d2[len(d2)-1],
		})
	}

	for _, condition := range Conditions {
		if condition.Reach(data) {
			return true
		}
	}
	return false
}
