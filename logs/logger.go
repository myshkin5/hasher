package logs

import "log"

var (
	Logger ExternalLogger
)

type ExternalLogger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

func init() {
	Logger = logger{}
}

type logger struct{}

func (l logger) Error(args ...interface{}) {
	log.Print(args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l logger) Info(args ...interface{}) {
	log.Print(args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l logger) Panic(args ...interface{}) {
	log.Panic(args...)
}

func (l logger) Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func (l logger) Warn(args ...interface{}) {
	log.Print(args)
}

func (l logger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Init(logLevel string) error {
	return nil
}
