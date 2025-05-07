package setup

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/jatis/sample-stack-golang/pkg/logger"
)

// InitTestLogger initializes a test logger that won't output anything
func InitTestLogger() {
	testLogger := logrus.New()
	// Use ioutil.Discard to properly handle log output but discard it
	testLogger.SetOutput(ioutil.Discard)
	
	// Set log level to error to minimize output
	testLogger.SetLevel(logrus.ErrorLevel)
	
	// Initialize the logger
	logger.Log = testLogger
	
	// Also set up any other loggers that might be used directly
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.ErrorLevel)
}
