package simpleLogger

type logger interface {
	Logf(string, ...interface{})
}
