/*
 * @Author: tr3e
 * @Date: 2019-11-26 20:50:14
 * @Last Modified by: tr3e
 * @Last Modified time: 2019-11-26 20:50:14
 */

package main

import (
	"fmt"
	"log"
)

const (
	logDebug = iota
	logInfo
	logError
	logFatal
	logNull
)

var (
	logLevel = logInfo
)

func debugLog(v ...interface{}) {
	if logLevel <= logDebug {
		prefixLog("[debug]", v...)
	}
}

func debugLogf(format string, v ...interface{}) {
	debugLog(fmt.Sprintf(format, v...))
}

func infoLog(v ...interface{}) {
	if logLevel <= logInfo {
		prefixLog("[info]", v...)
	}
}

func infoLogf(format string, v ...interface{}) {
	infoLog(fmt.Sprintf(format, v...))
}

func errorLog(v ...interface{}) {
	if logLevel <= logError {
		prefixLog("[error]", v...)
	}
}

func errorLogf(format string, v ...interface{}) {
	errorLog(fmt.Sprintf(format, v...))
}

func fatalLog(v ...interface{}) {
	if logLevel <= logFatal {
		log.Print("[FATAL]")
		log.Panicln(v...)
	}
}

func fatalLogf(format string, v ...interface{}) {
	fatalLog(fmt.Sprintf(format, v...))
}

func prefixLog(pre string, v ...interface{}) {
	if logLevel < logNull {
		log.Println(pre, fmt.Sprint(v...))
		// log.SetPrefix("")
	}
}
