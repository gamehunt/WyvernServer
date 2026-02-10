package logger

type LogLevel int

const (
	Debug LogLevel = iota
	Info 
	Warn 
	Error
)

func LevelToString(level LogLevel) string {
	switch level {
	case Debug:
		return "debug";
	case Info:
		return "info";
	case Warn:
		return "warn";
	case Error:
		return "error";
	default:
		return "info";
	}
}

func LevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

type LoggerConfig struct {
	Level LogLevel
}


