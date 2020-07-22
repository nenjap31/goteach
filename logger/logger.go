package logger

import (
	"context"
	"io"
	// "goteach/config"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	// MiddlewareLog logger.
	MiddlewareLog *rotatelogs.RotateLogs
	goteachLog     *rotatelogs.RotateLogs
	// Logger is logger instance.
	Logger zerolog.Logger
	//Log           *log.Logger
)

func init() {
	
	logdir := viper.GetString("logdir")
	logMaxAge := viper.GetInt("log_max_age")
	debug := viper.GetBool("debug")

	// default goteach log dir setting
	if !strings.HasPrefix(logdir, "/") {
		dir, _ := os.Getwd()
		logdir = dir + "/log"
	}

	if logMaxAge < 1 {
		// default 15 days
		logMaxAge = 15
	}

	// Set Middleware logging.
	MiddlewareLog, _ = rotatelogs.New(
		logdir+"/access_log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(logdir+"/access_log"),
		rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	// Set goteach logging.
	goteachLog, _ = rotatelogs.New(
		logdir+"/goteach_log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(logdir+"/goteach_log"),
		rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	// Set logger, format and level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	Logger = zerolog.New(goteachLog).With().Timestamp().Logger()
	//Logger = log.New(goteachLog, "goteach:", log.Lshortfile)
}

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) zerolog.Logger {
	return Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() zerolog.Context {
	return Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) zerolog.Logger {
	return Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) zerolog.Logger {
	return Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) zerolog.Logger {
	return Logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *zerolog.Event {
	return Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
