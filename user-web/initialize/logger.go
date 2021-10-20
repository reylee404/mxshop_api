package initialize

import "go.uber.org/zap"

func MustInitLogger() {
	err := InitLogger()
	if err != nil {
		panic(err)
	}
}

func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}
