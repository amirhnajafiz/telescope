package logr

import "go.uber.org/zap"

// NewZapLogger creates a new zap logger instance
func NewZapLogger(debug bool) (*zap.Logger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	// create a new zap logger instance based on the debug flag
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}
