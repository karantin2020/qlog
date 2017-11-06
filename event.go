package qlog

type Event struct {
	Logger  *Logger
	Data    []Field
	Time    time.Time
	Level   Level
	Message string
}
