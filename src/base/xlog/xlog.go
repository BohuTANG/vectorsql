// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	defaultlog *Log
)

type LogLevel int

const (
	DEBUG LogLevel = 1 << iota
	INFO
	WARNING
	ERROR
	FATAL
	PANIC
)

var LevelNames = [...]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
	PANIC:   "PANIC",
}

const (
	D_LOG_FLAGS int = log.LstdFlags | log.Lmicroseconds
)

type Log struct {
	opts *Options
	*log.Logger
}

func NewStdLog(opts ...Option) *Log {
	return NewXLog(os.Stdout, opts...)
}

func NewXLog(w io.Writer, opts ...Option) *Log {
	options := newOptions(opts...)

	l := &Log{
		opts: options,
	}
	l.Logger = log.New(w, l.opts.Name, D_LOG_FLAGS)
	defaultlog = l
	return l
}

func NewLog(w io.Writer, prefix string, flag int) *Log {
	l := &Log{}
	l.Logger = log.New(w, prefix, flag)
	return l
}

func GetLog() *Log {
	if defaultlog == nil {
		log := NewStdLog(Level(INFO))
		defaultlog = log
	}
	return defaultlog
}

func (t *Log) SetLevel(level string) {
	for i, v := range LevelNames {
		if strings.EqualFold(level, v) {
			t.opts.Level = LogLevel(i)
			return
		}
	}
}

func (t *Log) Debug(format string, v ...interface{}) {
	if DEBUG < t.opts.Level {
		return
	}
	t.log("\t [DEBUG] \t%s %s", fmt.Sprintf(format, v...), getFnName())
}

func (t *Log) Info(format string, v ...interface{}) {
	if INFO < t.opts.Level {
		return
	}
	t.log("\t [INFO] \t%s %s", fmt.Sprintf(format, v...), getFnName())
}

func (t *Log) Warning(format string, v ...interface{}) {
	if WARNING < t.opts.Level {
		return
	}
	t.log("\t [WARNING] \t%s %s", fmt.Sprintf(format, v...), getFnName())
}

func (t *Log) Error(format string, v ...interface{}) {
	if ERROR < t.opts.Level {
		return
	}
	t.log("\t [ERROR] \t%s %s", fmt.Sprintf(format, v...), getFnName())
}

func (t *Log) Fatal(format string, v ...interface{}) {
	if FATAL < t.opts.Level {
		return
	}
	t.log("\t [FATAL+EXIT] \t%s %s", fmt.Sprintf(format, v...), getFnName())
	os.Exit(1)
}

func (t *Log) Panic(format string, v ...interface{}) {
	if PANIC < t.opts.Level {
		return
	}
	msg := fmt.Sprintf("\t [PANIC] \t%s %s", fmt.Sprintf(format, v...), getFnName())
	t.log(msg)
	panic(msg)
}

func (t *Log) Close() {
	// nothing
}

func (t *Log) log(format string, v ...interface{}) {
	_ = t.Output(3, strings.Repeat(" ", 3)+fmt.Sprintf(format, v...)+"\n")
}

func getFnName() string {
	var fnName string

	pc, fn, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		fnName = "?()"
	} else {
		names := strings.Split(f.Name(), ".")
		fnName = names[len(names)-1]
	}
	return fmt.Sprintf("<%s@%s:%d>", fnName, filepath.Base(fn), line)
}
