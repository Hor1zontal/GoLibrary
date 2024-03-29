/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/10/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 * Desc: compatible log framework
 *******************************************************************************/
package log

import (
	//"os"
	"os"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"runtime"
	"path"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/pkg/errors"
)

var format = &log.TextFormatter{}

var logger = NewLogger("gok", format, false)

var logRoot string

var DEBUG = false

//调试版本日志带颜色
func Init(debug bool, tag string, logDir string) {
	DEBUG = debug
	format.ForceColors = DEBUG
	format.DisableTimestamp = DEBUG
	if tag == "" {
		tag = "gok"
	}
	logRoot = logDir
	configLocalFilesystemLogger(tag, logger)
}

func NewLogger(name string, formatter log.Formatter, local bool) *log.Logger {
	logger := log.New()
	logger.Formatter = formatter
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.Out = os.Stdout
	// Only log the warning severity or above.
	logger.Level = log.DebugLevel
	if local {
		configLocalFilesystemLogger(name, logger)
	}
	return logger
}

// config logrus log to amqp  rabbitMQ
//func ConfigAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
//	hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
//	log.AddHook(hook)
//}

// config logrus log to elasticsearch
//func ConfigESLogger(esUrl string, esHOst string, index string) {
//	client, err := elastic.NewClient(elastic.SetURL(esUrl))
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	log.AddHook(esHook)
//}

//config logrus log to local file
func configLocalFilesystemLogger(name string, logger *log.Logger) {
	maxAge := 30 * 24 * time.Hour
	rotationTime := 24 * time.Hour
	logFileName := name + ".log"
	baseLogPath := path.Join(logRoot, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	errWriter, err1 := rotatelogs.New(
		baseLogPath+".err.%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err1 != nil {
		log.Errorf("config local file system err logger error. %+v", errors.WithStack(err1))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: errWriter,
		log.FatalLevel: errWriter,
		log.PanicLevel: errWriter,
	}, logger.Formatter)
	logger.AddHook(lfHook)
}

//Debugf Printf Infof Warnf Warningf Errorf Panicf Fatalf

//做一层适配，方便后续切换到其他日志框架或者自己写

//-----------format
func WithField(key string, value interface{}) *log.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields log.Fields) *log.Entry {
	return logger.WithFields(fields)
}

func getLocation() string {
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("[%s:%d] ", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}

func Debug(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Debugf(format, arg...)
}

func Print(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Printf(format, arg...)
}

func Info(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Infof(format, arg...)
}

func Warn(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Warnf(format, arg...)
}

func Error(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Errorf(format, arg...)
}

func Panic(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Panicf(format, arg...)
}

func Fatal(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Fatalf(format, arg...)
}

