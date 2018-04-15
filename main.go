package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"unicode"

	link "github.com/tr3ee/go-link"
)

var (
	lNetwork, lAddr string
	rNetwork, rAddr string
)

func init() {
	flag.StringVar(&lNetwork, "lNet", "tcp", "Specifies the protocol for listening.")
	flag.StringVar(&rNetwork, "rNet", "tcp", "Specifies the remote protocol for dialing")
	flag.StringVar(&lAddr, "lhost", "", "Specifies the address for listening.")
	flag.StringVar(&rAddr, "rhost", "", "Specifies the remote address for dialing.")
}

func cLog(buf []byte) []byte {
	var logStr string
	for _, v := range buf {
		if unicode.IsPrint(rune(v)) && v != '\n' {
			logStr += string(v)
		} else {
			logStr += fmt.Sprintf("\\x%02x", v)
		}
	}
	infoLogf("SEND:%s", logStr)
	return buf
}

func sLog(buf []byte) []byte {
	var logStr string
	for _, v := range buf {
		if unicode.IsPrint(rune(v)) && v != '\n' {
			logStr += string(v)
		} else {
			logStr += fmt.Sprintf(`\x%02x`, v)
		}
	}
	infoLogf("RECV:%s", logStr)
	return buf
}

func main() {
	flag.Parse()
	if len(lAddr) == 0 || len(rAddr) == 0 {
		flag.Usage()
		return
	}

	l, err := net.Listen(lNetwork, lAddr)
	if err != nil {
		panic(err)
	}
	infoLogf("listening on %s:%s", lNetwork, lAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			target, err := net.Dial(rNetwork, rAddr)
			if err != nil {
				errorLogf("failed to connect to %s:%s", rNetwork, rAddr)
				return
			}
			defer target.Close()
			connAddr, tgtAddr := conn.RemoteAddr(), target.RemoteAddr()
			infoLogf("Hijacking link %v <==> %v", connAddr, tgtAddr)
			n1, n2, e1, e2 := link.TwoWayLink(nil, conn, target, cLog, sLog)
			infoLogf("%v <==> %v (%d transmitted,%d received)", connAddr, tgtAddr, n1, n2)
			if e1 != nil {
				debugLogf(" %v --> %v: %s", connAddr, tgtAddr, e1)
			}
			if e2 != nil {
				debugLogf(" %v <-- %v: %s", connAddr, tgtAddr, e1)
			}
		}(conn)
	}
}
