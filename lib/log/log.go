package log

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logger is a simple interface that can be used for logging purposes
type Logger interface {
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})

	Infoln(args ...interface{})
	Infof(format string, args ...interface{})

	Warnln(args ...interface{})
	Warnf(format string, args ...interface{})

	Errorln(args ...interface{})
	Errorf(format string, arg ...interface{})

	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})

	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
}

type logrusLogger struct {
	entry *logrus.Entry
}

func NewLogrus(cfg *Config) (Logger, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	l := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrap(err, "invalid log level")
	}

	l.SetLevel(lvl)
	l.SetReportCaller(cfg.IncludeCaller)

	switch cfg.Format {
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		l.SetFormatter(&logrus.TextFormatter{})
	default:
		return nil, errors.New("Invalid formatter")
	}

	switch cfg.Output {
	case "stdout":
		l.SetOutput(os.Stdout)
	case "stderr":
		l.SetOutput(os.Stderr)
	default:
		return nil, errors.New("invalid output")
	}

	return logrusLogger{
		entry: logrus.NewEntry(l),
	}, nil
}

func (l logrusLogger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l logrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l logrusLogger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l logrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l logrusLogger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

func (l logrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l logrusLogger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l logrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l logrusLogger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l logrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l logrusLogger) WithError(err error) Logger {
	return logrusLogger{l.entry.WithError(err)}
}

func (l logrusLogger) WithField(key string, value interface{}) Logger {
	return logrusLogger{l.entry.WithField(key, value)}
}
