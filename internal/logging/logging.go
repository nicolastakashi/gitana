package logging

import (
	"bytes"
	"os"

	"github.com/sirupsen/logrus"
)

type OutputSplitter struct{}

func (splitter *OutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=error")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

func Configure(logLevel string) error {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Error("error to parse log level %v", err)
		return err
	}

	logrus.SetLevel(lvl)
	logrus.SetOutput(&OutputSplitter{})
	return nil
}
