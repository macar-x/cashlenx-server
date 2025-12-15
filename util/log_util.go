package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	consoleLogger        *zap.Logger
	sugaredConsoleLogger *zap.SugaredLogger
	Logger               = getLogger()
)

func initLogger() {
	if consoleLogger == nil {
		consoleLogger = initConsoleLogger()
	}
	if sugaredConsoleLogger == nil {
		sugaredConsoleLogger = consoleLogger.Sugar()
	}

	consoleLogger.Debug("loggers initialize succeed")
}

func getLogger() *zap.SugaredLogger {
	if sugaredConsoleLogger == nil {
		initLogger()
	}
	return sugaredConsoleLogger
}

func initConsoleLogger() *zap.Logger {
	// 設定 console 輸出
	consoleOutput := zapcore.Lock(os.Stdout)
	// 設定 file 輸出
	fileOutput := zapcore.Lock(zapcore.AddSync(createLogFile()))

	// 設定日誌等級
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 合併輸出和編碼器
	consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, zapcore.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileOutput, zapcore.InfoLevel)

	// 同時使用 console 和 file 輸出
	logger := zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())

	// 記得在程序結束時關閉 logger
	// defer consoleLogger.Sync()

	return logger
}

func createLogFile() *os.File {
	var logFilePath string
	if customPath := GetConfigByKey("logger.file"); customPath != "" {
		logFilePath = customPath
	} else {
		// Create logs directory if it doesn't exist
		if err := os.MkdirAll("./logs", 0755); err != nil {
			log.Fatal("failed to create logs directory: ", err)
		}

		// Format log filename with current date: cashlenx_20251215.log
		currentDate := time.Now().Format("20060102")
		logFilePath = fmt.Sprintf("./logs/cashlenx_%s.log", currentDate)
	}

	// Open file with append mode if it exists, create if it doesn't
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to create log file: ", err)
	}
	return file
}
