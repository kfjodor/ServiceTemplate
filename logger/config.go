package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Async        asyncLoggerConfig
	Level        Level
	NoLock       bool
	HidePackages []string
}

type Level logrus.Level

type asyncLoggerConfig struct {
	QueueLength int
	BufferSize  int
	Enabled     bool
}

func (l *Level) Decode(data interface{}) error {
	s, ok := data.(string)
	if !ok {
		return errors.Errorf("invalid type: %T", data)
	}

	logrusLevel, err := logrus.ParseLevel(s)
	if err != nil {
		return err
	}

	*l = Level(logrusLevel)
	return nil
}
