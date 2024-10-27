package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// Basic colors
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorReset  = "\033[0m"

	// Text styles
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
	colorItalic  = "\033[3m"
	colorReverse = "\033[7m"

	// Background colors
	bgRed     = "\033[41m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func getColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return colorBold + colorCyan // Bold cyan for better debug visibility
	case INFO:
		return colorGreen
	case WARN:
		return colorBold + colorYellow + bgYellow + colorGray // Yellow background for warnings
	case ERROR:
		return colorBold + colorRed + bgRed + colorGray // Red background for errors
	case FATAL:
		return colorBold + colorGray + bgMagenta // Magenta background for fatal errors
	default:
		return colorReset
	}
}

func getLevelEmoji(level LogLevel) string {
	switch level {
	case DEBUG:
		return "🔍"
	case INFO:
		return "ℹ️ "
	case WARN:
		return "⚠️ "
	case ERROR:
		return "❌"
	case FATAL:
		return "💀"
	default:
		return "•"
	}
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

func formatTimestamp(t time.Time) string {
	return fmt.Sprintf("%s%s%s",
		colorBlue+colorBold,
		t.Format("2006-01-02 15:04:05.000"),
		colorReset,
	)
}

func formatLevel(level LogLevel) string {
	color := getColor(level)
	emoji := getLevelEmoji(level)
	levelStr := strings.ToUpper(level.String())
	padding := strings.Repeat(" ", 5-len(levelStr)) // Align level text

	return fmt.Sprintf("%s%s %s%s%s",
		color,
		emoji,
		padding+levelStr,
		colorReset,
		colorDim+" |"+colorReset, // Add separator
	)
}

func formatCaller(caller string) string {
	return fmt.Sprintf("%s%s%s%s",
		colorGray+colorItalic,
		caller,
		colorReset,
		colorDim+" |"+colorReset, // Add separator
	)
}

func formatMessage(msg string, level LogLevel) string {
	color := getColor(level)
	return fmt.Sprintf("%s%s%s",
		color,
		msg,
		colorReset,
	)
}

func formatField(key string, value interface{}) string {
	return fmt.Sprintf("%s%s%s=%s%v%s",
		colorCyan+colorBold,
		key,
		colorReset,
		colorItalic,
		value,
		colorReset,
	)
}

func logMessage(level LogLevel, message string, fields ...map[string]interface{}) {
	timestamp := formatTimestamp(time.Now())
	levelStr := formatLevel(level)
	callerInfo := formatCaller(getCallerInfo())
	formattedMsg := formatMessage(message, level)

	// Base log message
	logMsg := fmt.Sprintf("%s %s %s %s",
		timestamp,
		levelStr,
		callerInfo,
		formattedMsg,
	)

	// Add fields if present
	if len(fields) > 0 && fields[0] != nil {
		logMsg += colorDim + " |" + colorReset
		for k, v := range fields[0] {
			logMsg += " " + formatField(k, v)
		}
	}

	logger.Println(logMsg)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func Debug(message string, fields ...map[string]interface{}) {
	logMessage(DEBUG, message, fields...)
}

// Info logs an info message
func Info(message string, fields ...map[string]interface{}) {
	logMessage(INFO, message, fields...)
}

// Warn logs a warning message
func Warn(message string, fields ...map[string]interface{}) {
	logMessage(WARN, message, fields...)
}

// Error logs an error message
func Error(message string, fields ...map[string]interface{}) {
	logMessage(ERROR, message, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(message string, fields ...map[string]interface{}) {
	logMessage(FATAL, message, fields...)
}

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// Helper function to create fields
func Fields(fields map[string]interface{}) map[string]interface{} {
	return fields
}
