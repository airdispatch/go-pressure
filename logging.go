package pressure

import (
	"fmt"
	"time"
)

const ERROR int = 2
const WARNING int = 1
const DEBUG int = 0

const (
	reset      = "\x1b[0m"
	bright     = "\x1b[1m"
	dim        = "\x1b[2m"
	underscore = "\x1b[4m"
	blink      = "\x1b[5m"
	reverse    = "\x1b[7m"
	hidden     = "\x1b[8m"

	fgBlack   = "\x1b[30m"
	fgRed     = "\x1b[31m"
	fgGreen   = "\x1b[32m"
	fgYellow  = "\x1b[33m"
	fgBlue    = "\x1b[34m"
	fgMagenta = "\x1b[35m"
	fgCyan    = "\x1b[36m"
	fgWhite   = "\x1b[37m"

	bgBlack   = "\x1b[40m"
	bgRed     = "\x1b[41m"
	bgGreen   = "\x1b[42m"
	bgYellow  = "\x1b[43m"
	bgBlue    = "\x1b[44m"
	bgMagenta = "\x1b[45m"
	bgCyan    = "\x1b[46m"
	bgWhite   = "\x1b[47m"
)

type Logger struct {
	LogLevel int
}

func NewLogger(log_level int) *Logger {
	return &Logger{log_level}
}

func (l *Logger) LogError(objects ...interface{}) {
	fmt.Print(fgRed, bright)
	l.logMessageAtLevel(ERROR, objects...)
}

func (l *Logger) LogWarning(objects ...interface{}) {
	fmt.Print(fgYellow)
	l.logMessageAtLevel(WARNING, objects...)
}

func (l *Logger) LogDebug(objects ...interface{}) {
	fmt.Print(fgCyan)
	l.logMessageAtLevel(DEBUG, objects...)
}

func (l *Logger) logMessageAtLevel(level int, objects ...interface{}) {
	if l.LogLevel <= level {
		// Mon Jan 2 15:04:05 -0700 MST 2006
		fmt.Print("[", time.Now().Format("02/01/2006 - 3:04:05.99 PM (MST)"), "]: ")
		fmt.Println(objects...)
		fmt.Print(reset)
	}
}
