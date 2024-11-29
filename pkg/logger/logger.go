package logger

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"time"
)

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case ErrorLevel:
		return "Error"
	default:
		return ""
	}
}

type Logger struct {
	out   io.Writer
	level Level
}

var l *Logger

func New(out io.Writer, level Level) *Logger {
	return &Logger{
		out:   out,
		level: level,
	}
}

func PrintDebug(message string, properties map[string]any) {
	l.print(DebugLevel, message, properties)
}
func PrintInfo(message string, properties map[string]any) {
	l.print(InfoLevel, message, properties)
}
func PrintError(err error, properties map[string]any) {
	l.print(ErrorLevel, err.Error(), properties)
}

func (l *Logger) print(level Level, message string, properties map[string]any) (int, error) {
	if level < l.level {
		return 0, nil
	}

	aux := struct {
		Level      string         `json:"level"`
		Time       string         `json:"time"`
		Message    string         `json:"message,omitempty"`
		Properties map[string]any `json:"properties,omitempty"`
		Trace      string         `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= ErrorLevel {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(ErrorLevel.String() + ": unable to marshal log message: " + err.Error())
	}

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(ErrorLevel, string(message), nil)
}
