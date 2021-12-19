package simpleLogger

type Logger interface {
	Logf(string, ...interface{})
}
