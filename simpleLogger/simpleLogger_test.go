package simpleLogger

import (
	"log"
	"testing"
)

type SimpleObj struct {
	logger Logger
}

func NewTestObj() *SimpleObj {
	ret := SimpleObj{
		logger: DefaultLogger(WithCallPath(2)),
	}
	return &ret
}

func TestLogger(t *testing.T) {
	o := NewTestObj()
	o.logger.Logf("hello world")
}

type advancedObj struct {
	logger     Logger
	detailOpen bool
}

func NewTestCallPathObj(detailOpen bool) *advancedObj {
	ret := advancedObj{
		logger: NewLogger(
			WithCallPath(3),
			WithFlag(log.Lshortfile),
		),
		detailOpen: detailOpen,
	}
	return &ret
}

func (t advancedObj) err(layout string, v ...interface{}) {
	t.logger.Logf(layout, v...)
}

func (t advancedObj) info(layout string, v ...interface{}) {
	if !t.detailOpen {
		return
	}
	t.logger.Logf(layout, v...)
}

func ExampleDetailOpenTrue() {
	o := NewTestCallPathObj(true)
	o.info("hello world")
	o.err("hello world")
	// output:
	//simpleLogger_test.go:53: hello world
	//simpleLogger_test.go:54: hello world
}

func ExampleDetailOpenFalse() {
	o := NewTestCallPathObj(false)
	o.info("hello world")
	o.err("hello world")
	// output:
	//simpleLogger_test.go:63: hello world
}
