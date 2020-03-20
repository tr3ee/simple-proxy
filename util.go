package main

import (
	"log"

	"github.com/fatih/color"
)

var (
	colorDebug  = color.BlueString
	colorInfo   = color.MagentaString
	colorWarn   = color.YellowString
	colorFatal  = color.RedString
	colorLocal  = color.CyanString
	colorRemote = color.GreenString
)

func debugf(fmtstr string, args ...interface{}) {
	if verboseverbose {
		log.Println(colorDebug("[-] "+fmtstr, args...))
	}
}

func infof(fmtstr string, args ...interface{}) {
	if verbose {
		log.Println(colorInfo("[+] "+fmtstr, args...))
	}
}

func warnf(fmtstr string, args ...interface{}) {
	log.Println(colorWarn("[*] "+fmtstr, args...))
}

func fatalf(fmtstr string, args ...interface{}) {
	log.Fatalln(colorFatal("[!] "+fmtstr, args...))
}
