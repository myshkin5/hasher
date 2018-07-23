package logs

import (
	"errors"
	"fmt"
	"log"
)

var (
	Logger         ExternalLogger
	internalLogger *logger
)

type Level uint32

const (
	PanicLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
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
	internalLogger = &logger{level: PanicLevel}
	Logger = internalLogger
}

type logger struct {
	level Level
}

func (l *logger) Error(args ...interface{}) {
	if l.level >= ErrorLevel {
		args[0] = fmt.Sprintf("ERROR %s", args[0])
		log.Println(args...)
	}
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		log.Printf("ERROR "+format, args...)
	}
}

func (l *logger) Info(args ...interface{}) {
	if l.level >= InfoLevel {
		args[0] = fmt.Sprintf("INFO %s", args[0])
		log.Println(args...)
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.level >= InfoLevel {
		log.Printf("INFO "+format, args...)
	}
}

func (l *logger) Panic(args ...interface{}) {
	args[0] = fmt.Sprintf("PANIC %s", args[0])
	log.Panicln(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	log.Panicf("PANIC "+format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	if l.level >= WarnLevel {
		args[0] = fmt.Sprintf("WARN %s", args[0])
		log.Println(args...)
	}
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.level >= WarnLevel {
		log.Printf("WARN "+format, args...)
	}
}

func Init(logLevel string) error {
	switch logLevel {
	case "panic":
		internalLogger.level = PanicLevel
	case "error":
		internalLogger.level = ErrorLevel
	case "warn":
		internalLogger.level = WarnLevel
	case "info":
		internalLogger.level = InfoLevel
	default:
		return errors.New("Invalid log level: " + logLevel)
	}
	return nil
}
