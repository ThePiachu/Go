package Log

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"golang.org/x/net/context"
	//"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"runtime"
	"strconv"
)

func Debugf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Debugf(c, file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Infof(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Infof(c, file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Warningf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Warningf(c, file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Errorf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Errorf(c, file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Criticalf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Criticalf(c, file+":"+strconv.Itoa(line)+" - "+format, args...)
}

//JSON debugs

func argsToJSON(args ...interface{}) []interface{} {
	answer := []interface{}{}
	for _, v := range args {
		s := encodeJSONString(v)
		answer = append(answer, s)
	}
	return answer
}

func encodeJSONString(data interface{}) string {
	encoded, _ := json.Marshal(data)
	return string(encoded)
}

func JDebugf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Debugf(c, file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JInfof(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Infof(c, file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JWarningf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Warningf(c, file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JErrorf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Errorf(c, file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JCriticalf(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Criticalf(c, file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}
