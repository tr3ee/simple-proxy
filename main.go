/*
 * @Author: tr3e
 * @Date: 2019-11-26 20:50:34
 * @Last Modified by: tr3e
 * @Last Modified time: 2019-11-26 21:20:22
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/tr3ee/go-link"
)

var (
	verbose        bool
	verboseverbose bool
	decrypt        bool
	lNet           string
	lAddr          string
	rNet           string
	rAddr          string
	method         string
	secret         string
	timeout        int
)

func init() {
	methods := make([]string, 0, len(cipherMethod))
	for name := range cipherMethod {
		methods = append(methods, name)
	}

	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&verboseverbose, "vv", false, "more verbose mode")
	flag.BoolVar(&decrypt, "d", false, "decryption mode")

	flag.StringVar(&lNet, "ln", "tcp", "local network protocal will be used when listening")
	flag.StringVar(&lAddr, "l", "", "local address to listen on")

	flag.StringVar(&rNet, "rn", "tcp", "remote network protocal will be used when connecting")
	flag.StringVar(&rAddr, "r", "", "remote address to connect")

	flag.StringVar(&method, "m", "", fmt.Sprintf("cipher method (currently support: %s)", strings.Join(methods, "|")))
	flag.StringVar(&secret, "k", "", "secret key for cipher")

	flag.IntVar(&timeout, "t", 10, "idle timeout for each connection")
}

func main() {
	flag.Parse()
	if len(lNet) == 0 || len(lAddr) == 0 || len(rNet) == 0 || len(rAddr) == 0 || timeout <= 0 {
		flag.Usage()
		return
	}

	if verboseverbose {
		verbose = true
	}

	cipher, err := NewCipher(method, secret)
	if err != nil {
		log.Fatalf("[FATAL] failed to initialize cipher: %s", err)
	}

	listener, err := net.Listen(lNet, lAddr)
	if err != nil {
		log.Fatalf("[FATAL] failed to listen on given addres: %s", err)
	}

	if verbose {
		log.Printf("[INFO] listening on %s:%s", lNet, lAddr)
	}

	for {
		lc, err := listener.Accept()
		if err != nil {
			log.Printf("[WARN] error occured on Accept: %s", err)
			continue
		}

		rc, err := net.Dial(rNet, rAddr)
		if err != nil {
			log.Printf("[WARN] failed to connect remote address: %s", err)
			continue
		}

		go func() {
			defer lc.Close()
			defer rc.Close()
			c := cipher.Copy()
			ehook, dhook := NewHook(), NewHook()
			ehook.Add(c.Encrypt)
			dhook.Add(c.Decrypt)
			if verbose {
				ehook.Add(func(p []byte) []byte {
					fmt.Printf("[SEND] %s", p)
					return p
				})
				dhook.Add(func(p []byte) []byte {
					fmt.Printf("[RECV] %s", p)
					return p
				})
			}
			ehook.Add(func(p []byte) []byte {
				lc.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				return p
			})
			dhook.Add(func(p []byte) []byte {
				rc.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				return p
			})
			if decrypt {
				ehook, dhook = dhook, ehook
			}
			lsend, rsend, lerr, rerr := link.TwoWayLinkSpec(nil, lc, rc, nil, nil, ehook.Run, dhook.Run)
			if verbose {
				log.Printf("[INFO] %v <==> %v (%d transmitted, %d received)", lc.RemoteAddr(), rc.RemoteAddr(), lsend, rsend)
				if verboseverbose {
					if lerr != nil {
						log.Printf("[DEBUG] %v --> %v: %s", lc.RemoteAddr(), rc.RemoteAddr(), lerr)
					}
					if rerr != nil {
						log.Printf("[DEBUG] %v <-- %v: %s", lc.RemoteAddr(), rc.RemoteAddr(), rerr)
					}
				}
			}
		}()
	}
}
