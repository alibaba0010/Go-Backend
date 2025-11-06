package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}	
}

func Sync(){
	if Log != nil {
		_ = Log.Sync()
	}
}