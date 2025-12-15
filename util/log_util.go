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
	// Get the log folder from config, default to ./logs
	logFolder := GetConfigByKey("logger.file")

	// Extract folder path if a full filename was provided
	if logFolder != "" {
		// If the path has a file extension, extract just the directory
		if len(logFolder) > 4 && logFolder[len(logFolder)-4:] == ".log" {
			// Check if it's just a filename without directory
			if logFolder == "./cashlenx.log" || logFolder == "cashlenx.log" {
				logFolder = "./logs"
			} else {
				// Extract directory part
				lastSlash := -1
				for i := len(logFolder) - 1; i >= 0; i-- {
					if logFolder[i] == '/' || logFolder[i] == '\\' {
						lastSlash = i
						break
					}
				}
				if lastSlash >= 0 {
					logFolder = logFolder[:lastSlash]
				} else {
					logFolder = "."
				}
			}
		}
	} else {
		logFolder = "./logs"
	}

	// Create the log folder if it doesn't exist
	if err := os.MkdirAll(logFolder, 0755); err != nil {
		log.Fatal("failed to create logs directory: ", err)
	}

	// Always use date-based filename in the configured folder
	currentDate := time.Now().Format("20060102")
	logFilePath := fmt.Sprintf("%s/cashlenx_%s.log", logFolder, currentDate)

	// Open file with append mode if it exists, create if it doesn't
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to create log file: ", err)
	}
	return file
}
