package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelSuccess
	LevelWarn
	LevelError
)

type Logger struct {
	out      io.Writer
	isTTY    bool
	useJSON  bool
	hasColor bool
}

func NewLogger(out io.Writer, useJSON bool) *Logger {
	l := &Logger{
		out:     out,
		useJSON: useJSON,
	}

	// Проверяем, является ли вывод терминалом (валидация TTY)
	if file, ok := out.(*os.File); ok {
		stat, err := file.Stat()
		if err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
			l.isTTY = true
		}
	}

	// Проверяем глобальный стандарт NO_COLOR и наличие TTY
	noColor := os.Getenv("NO_COLOR") != "" || strings.ToLower(os.Getenv("COLORS")) == "no"
	if l.isTTY && !noColor && !useJSON {
		l.hasColor = true
	}

	return l
}

func (l *Logger) log(level Level, msg string, attrs map[string]interface{}) {
	if l.useJSON {
		l.logJSON(level, msg, attrs)
		return
	}
	l.logText(level, msg, attrs)
}

func (l *Logger) logText(level Level, msg string, attrs map[string]interface{}) {
	var levelStr, colorStart, colorReset string

	if l.hasColor {
		colorReset = "\033[0m"
	}

	switch level {
	case LevelInfo:
		levelStr = "[INFO]"
		if l.hasColor {
			colorStart = "\033[36m"
		} // Cyan
	case LevelSuccess:
		levelStr = "[OK]"
		if l.hasColor {
			colorStart = "\033[32m"
		} // Green
	case LevelWarn:
		levelStr = "[WARN]"
		if l.hasColor {
			colorStart = "\033[33m"
		} // Yellow
	case LevelError:
		levelStr = "[ERR]"
		if l.hasColor {
			colorStart = "\033[31m"
		} // Red
	}

	// Сборка строки лога для человека
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.out, "%s %s%s%-7s%s %s", timeStr, colorStart, "", levelStr, colorReset, msg)

	if len(attrs) > 0 {
		fmt.Fprintf(l.out, " | %v", attrs)
	}
	fmt.Fprintln(l.out)
}

func (l *Logger) logJSON(level Level, msg string, attrs map[string]interface{}) {
	var lvl string
	switch level {
	case LevelInfo:
		lvl = "INFO"
	case LevelSuccess:
		lvl = "SUCCESS"
	case LevelWarn:
		lvl = "WARN"
	case LevelError:
		lvl = "ERROR"
	}

	payload := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": lvl,
		"msg":   msg,
	}
	for k, v := range attrs {
		payload[k] = v
	}
	bytes, _ := json.Marshal(payload)
	l.out.Write(append(bytes, '\n'))
}

// Публичные методы логгера
func (l *Logger) Info(msg string, a ...interface{}) { l.log(LevelInfo, fmt.Sprintf(msg, a...), nil) }
func (l *Logger) Success(msg string, a ...interface{}) {
	l.log(LevelSuccess, fmt.Sprintf(msg, a...), nil)
}
func (l *Logger) Warn(msg string, a ...interface{})  { l.log(LevelWarn, fmt.Sprintf(msg, a...), nil) }
func (l *Logger) Error(msg string, a ...interface{}) { l.log(LevelError, fmt.Sprintf(msg, a...), nil) }

func (l *Logger) Prompt(msg string) {
	if l.hasColor {
		fmt.Print("\033[35m" + msg + "\033[0m") // Magenta для ввода
	} else {
		fmt.Print(msg)
	}
}
