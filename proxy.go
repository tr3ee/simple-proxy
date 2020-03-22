package main

import (
	"net"
	"time"
)

type proxyFunc func(proxyHandle)

type proxyHandle func(local *proxyConn, remote *proxyConn, readed []byte)

type proxyConn struct {
	net.Conn
	expired time.Time
}

func newProxyConn(conn net.Conn) *proxyConn {
	return &proxyConn{
		Conn:    conn,
		expired: time.Now().Add(time.Duration(timeout) * time.Second),
	}
}

func (c *proxyConn) Expired() bool {
	return time.Now().After(c.expired)
}

var proxyModes = map[string]proxyFunc{
	"active":  active,
	"forward": forward,
	"passive": passive,
}

func forward(handle proxyHandle) {
	listener, err := net.Listen(lNet, lAddr)
	if err != nil {
		fatalf("failed to listen on given addres: %s", err)
	}

	infof("listening on %s:%s", lNet, lAddr)

	for {
		lc, err := listener.Accept()
		if err != nil {
			warnf("got error on accept(): %s", err)
			continue
		}

		debugf("received connection from local listener %s", lc.RemoteAddr())

		rc, err := net.Dial(rNet, rAddr)
		if err != nil {
			lc.Close()
			warnf("failed to connect to remote address: %s", err)
			continue
		}

		go handle(newProxyConn(lc), newProxyConn(rc), nil)
	}
}

func active(handle proxyHandle) {
	for {
		lc, err := net.Dial(lNet, lAddr)
		if err != nil {
			fatalf("failed to connect to local address: %s", err)
		}
		lc.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

		buf := make([]byte, 1024)
		n, err := lc.Read(buf)
		if err != nil {
			lc.Close()
			debugf("failed to read from local address: %s", err)
			continue
		}

		rc, err := net.Dial(rNet, rAddr)
		if err != nil {
			lc.Close()
			warnf("failed to connect to remote address: %s", err)
			continue
		}

		go handle(newProxyConn(lc), newProxyConn(rc), buf[:n])
	}
}

func passive(handle proxyHandle) {
	ll, err := net.Listen(lNet, lAddr)
	if err != nil {
		fatalf("failed to listen on given local addres: %s", err)
	}
	infof("listening on address %s:%s (local)", lNet, lAddr)

	rl, err := net.Listen(rNet, rAddr)
	if err != nil {
		fatalf("failed to listen on given remote addres: %s", err)
	}
	infof("listening on address %s:%s (remote)", rNet, rAddr)

	connChannel := make(chan *proxyConn, passiveQSize)

	go func() {
		for {
			lc, err := ll.Accept()
			if err != nil {
				warnf("got error on local address accept(): %s", err)
				continue
			}
			debugf("received connection from local listener %s", lc.RemoteAddr())
			connChannel <- newProxyConn(lc)
		}
	}()

	for {
		rc, err := rl.Accept()
		if err != nil {
			warnf("got error on remote address accept(): %s", err)
			continue
		}
		debugf("received connection from remote listener %s", rc.RemoteAddr())

	fetch:
		remote := newProxyConn(rc)
		expired := time.After(time.Duration(timeout) * time.Second)
		select {
		case local := <-connChannel:
			if local.Expired() {
				local.Close()
				goto fetch
			}
			go handle(local, remote, nil)
		case <-expired:
			remote.Close()
		}
	}
}
