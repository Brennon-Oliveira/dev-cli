package logger

import (
	"fmt"
	"io"

	loggerutils "github.com/Brennon-Oliveira/dev-cli/internal/logger/logger_utils"
)

func (l *realLogger) SetVerbose(verbose bool) {
	l.verbose = verbose
	l.w.SetAllowWriting(verbose)
}

func SetVerbose(verbose bool) {
	logger.SetVerbose(verbose)
}

func (l *realLogger) SetOutput(writer io.Writer) {
	l.w = NewLoggerWriter(writer, l.verbose)
}

func SetOutput(writer io.Writer) {
	logger.SetOutput(writer)
}

func (l *realLogger) Info(format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s%s➜%s %s\n", loggerutils.RegularCyanColor, loggerutils.BoldStyle, loggerutils.ResetColor, msg)
}

func Info(format string, args ...any) {
	logger.Info(format, args...)
}

func (l *realLogger) Verbose(format string, args ...any) {
	if !l.verbose {
		return
	}
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s  │ %s%s\n", loggerutils.HighIntensityBlackColor, msg, loggerutils.ResetColor)
}

func Verbose(format string, args ...any) {
	logger.Verbose(format, args...)
}

func (l *realLogger) VerboseFromOutput(output io.Reader) {
	if !l.verbose {
		return
	}

	fmt.Fprintf(l.w.dest, "%s  │ ", loggerutils.HighIntensityBackgroundBlackColor)

	io.Copy(l.w.dest, output)

	fmt.Fprintf(l.w.dest, "%s\n", loggerutils.ResetColor)
}

func (l *realLogger) Debug(format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s%s➜ %s%s\n", loggerutils.RegularBlueColor, loggerutils.BoldStyle, msg, loggerutils.ResetColor)
}

func Debug(format string, args ...any) {
	logger.Debug(format, args...)
}

func VerboseFromOutput(output io.Reader) {
	logger.VerboseFromOutput(output)
}

func (l *realLogger) Success(format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s✓ %s%s\n", loggerutils.RegularGreenColor, msg, loggerutils.ResetColor)
}

func Success(format string, args ...any) {
	logger.Success(format, args...)
}

func (l *realLogger) Warn(format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s⚠ %s%s\n", loggerutils.RegularYellowColor, msg, loggerutils.ResetColor)
}

func Warn(format string, args ...any) {
	logger.Warn(format, args...)
}

func (l *realLogger) Error(format string, args ...any) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	fmt.Fprintf(l.w.dest, "%s✗ %s%s\n", loggerutils.RegularRedColor, msg, loggerutils.ResetColor)
}

func Error(format string, args ...any) {
	logger.Error(format, args...)
}

func (l *realLogger) GetWriter() io.Writer {
	return l.w
}

func GetWriter() io.Writer {
	return logger.GetWriter()
}
