/*****************************************************************************
File name: logger.go
Description: 日志
Author: failymao
Version: V1.0
Date: 2018/06/14
History:
*****************************************************************************/
package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"platform/config"
	"path/filepath"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
	size	   int64
}

func New(strLevel string, pathname string)  error {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()

		filename := fmt.Sprintf("%02d%02d_%02d_%02d_%02d%03d.log",
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			now.Nanosecond()/1000000)

		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return err
		}

		baseLogger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	// new
	gLogger = new(Logger)
	gLogger.level = level
	gLogger.baseLogger = baseLogger
	gLogger.baseFile = baseFile
	gLogger.size = 0
	return nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}
	if logger.size > int64(config.Gconfig.Logcfg.LogMax * 1024 * 1024) {
		logger.rotate()
	}
	format = printLevel + format
	outstr := fmt.Sprintf(format, a...)
	logger.size += int64(len(outstr))
	logger.baseLogger.Output(3, outstr)

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) rotate() error {

	if logger.baseFile != nil {

		fileInfo, err := logger.baseFile.Stat()
		if err != nil {
			fmt.Printf("log file get file info err:%s", err.Error())
			return err
		}

		now := time.Now()
		filename := fmt.Sprintf("%02d%02d_%02d_%02d_%02d:%03d.log",
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			now.Nanosecond()/1000000)

		absPath, err := filepath.Abs(fileInfo.Name())
		if err != nil {
			fmt.Printf("log file get file path err:%s", err.Error())
			return err
		}

		absDir := filepath.Dir(absPath)
		file, err := os.Create(path.Join(absDir, filename))
		if err != nil {
			fmt.Printf("log create new err:%s", err.Error())
			return err
		}
		logger.baseLogger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		logger.baseFile = file
		logger.size = 0
	}
	return nil
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger *Logger


func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}


