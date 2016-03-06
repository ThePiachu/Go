package Log

// Copyright 2014 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"appengine"
	"runtime"
	"strconv"
	"encoding/json"
)

func Debugf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Debugf(file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Infof(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Infof(file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Warningf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Warningf(file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Errorf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Errorf(file+":"+strconv.Itoa(line)+" - "+format, args...)
}

func Criticalf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Criticalf(file+":"+strconv.Itoa(line)+" - "+format, args...)
}

//JSON debugs

func argsToJSON(args ...interface{}) []interface{} {
	answer:=[]interface{}{}
	for _, v:=range(args) {
		s:=encodeJSONString(v)
		answer = append(answer, s)
	}
	return answer
}

func encodeJSONString(data interface{}) string {
	encoded, _ := json.Marshal(data)
	return string(encoded)
}

func JDebugf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Debugf(file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JInfof(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Infof(file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JWarningf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Warningf(file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JErrorf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Errorf(file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}

func JCriticalf(c appengine.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	c.Criticalf(file+":"+strconv.Itoa(line)+" - "+format, argsToJSON(args)...)
}
