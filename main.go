/*
 * @Author: tr3e
 * @Date: 2019-11-26 20:50:34
 * @Last Modified by: tr3e
 * @Last Modified time: 2019-11-26 21:20:22
 */

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tr3ee/go-link"
)

var (
	nocolor        bool
	verbose        bool
	verboseverbose bool
	decrypt        bool
	hexadecimal    bool
	mode           string
	lNet           string
	lAddr          string
	rNet           string
	rAddr          string
	method         string
	secret         string
	timeout        int
)

func init() {
	flag.BoolVar(&nocolor, "no-color", false, "disable color output")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&verboseverbose, "vv", false, "more verbose mode")
	flag.BoolVar(&decrypt, "d", false, "decryption mode")
	flag.BoolVar(&hexadecimal, "hex", false, "print inbound/outbound data in hexadecimal format")

	flag.StringVar(&lNet, "ln", "tcp", "local network protocal will be used when listening")
	flag.StringVar(&lAddr, "l", "", "local address to listen on")

	flag.StringVar(&rNet, "rn", "tcp", "remote network protocal will be used when connecting")
	flag.StringVar(&rAddr, "r", "", "remote address to connect")

	flag.StringVar(&method, "m", "plain", fmt.Sprintf("cipher method (currently support: %s)", strings.Join(supportedCiphers, "|")))
	flag.StringVar(&secret, "k", "", "secret key for cipher")
	flag.StringVar(&mode, "mode", "forward", fmt.Sprintf("proxy mode (currently support: %s)", strings.Join(supportedModes, "|")))

	flag.IntVar(&timeout, "t", 10, "idle timeout for each connection")
}

func main() {
	flag.Parse()
	if len(lNet) == 0 || len(lAddr) == 0 || len(rNet) == 0 || len(rAddr) == 0 || timeout <= 0 {
		flag.Usage()
		return
	}

	if nocolor {
		color.NoColor = nocolor
	}

	if verboseverbose {
		verbose = true
	}

	cipher, err := NewCipher(method, secret)
	if err != nil {
		fatalf("failed to initialize cipher: %s", err)
	}

	proxy := proxyModes[mode]

	if proxy == nil {
		fatalf("unrecognized proxy mode %q (currently support %s)", mode, fmt.Sprintf(strings.Join(supportedModes, "|")))
	}

	proxy(func(local, remote *proxyConn, lread []byte) {
		defer local.Close()
		defer remote.Close()
		c := cipher.Copy()
		lhook, rhook := NewHook(), NewHook()
		if !decrypt {
			rhook.Add(c.Decrypt)
		} else {
			lhook.Add(c.Decrypt)
		}
		if verbose {
			lhook.Add(func(p []byte) []byte {
				output := interface{}(p)
				if hexadecimal {
					output = hex.Dump(p)
				}
				if color.NoColor {
					fmt.Printf("[LOCAL]\n%s", output)
				} else {
					fmt.Print(colorLocal("%s", output))
				}
				return p
			})
			rhook.Add(func(p []byte) []byte {
				output := interface{}(p)
				if hexadecimal {
					output = hex.Dump(p)
				}
				if color.NoColor {
					fmt.Printf("[REMOTE]\n%s", output)
				} else {
					fmt.Print(colorRemote("%s", output))
				}
				return p
			})
		}
		if !decrypt {
			lhook.Add(c.Encrypt)
		} else {
			rhook.Add(c.Encrypt)
		}

		lhook.Add(func(p []byte) []byte {
			local.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			return p
		})
		rhook.Add(func(p []byte) []byte {
			remote.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			return p
		})
		var lsendfix int64
		if len(lread) > 0 {
			lread = lhook.Run(lread)
			n, err := remote.Write(lread)
			if err != nil {
				warnf("failed to send to remote address: %s", err)
				return
			}
			lsendfix += int64(n)
		}
		lsend, rsend, lerr, rerr := link.TwoWayLinkSpec(nil, local, remote, lread, nil, lhook.Run, rhook.Run)
		if lerr != nil {
			debugf("%v --> %v: %s", local.RemoteAddr(), remote.RemoteAddr(), lerr)
		}
		if rerr != nil {
			debugf("%v <-- %v: %s", local.RemoteAddr(), remote.RemoteAddr(), rerr)
		}
		infof("%v <==> %v (%d transmitted, %d received)", local.RemoteAddr(), remote.RemoteAddr(), lsend+lsendfix, rsend)
	})

	infof("bye!")
}
