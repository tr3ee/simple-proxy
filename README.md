# simple-proxy

a simple proxy POC based on my [go-link](https://github.com/tr3ee/go-link)

## Install

```bash
$ go get github.com/tr3ee/simple-proxy
```
this command should install `simple-proxy` in your $GOPATH/bin directory

## Usage
```
Usage of simple-proxy.exe:
  -lNet string
        Specifies the protocol for listening. (default "tcp")
  -lhost string
        Specifies the address for listening.
  -rNet string
        Specifies the remote protocol for dialing (default "tcp")
  -rhost string
        Specifies the remote address for dialing.
```

## TEST

__[SERVER]__
```bash
$ simple-proxy -rhost google.com:80 -lhost 127.0.0.1:8800
2018/04/15 15:53:45 [info] listening on tcp:127.0.0.1:8800
2018/04/15 15:53:51 [info] Hijacking link 127.0.0.1:51650 <==> [2404:6800:4005:80e::200e]:80
2018/04/15 15:53:57 [info] SEND:GET / HTTP/1.1\x0a
2018/04/15 15:54:01 [info] SEND:HOST:google.com\x0a
2018/04/15 15:54:02 [info] SEND:\x0a
2018/04/15 15:54:02 [info] RECV:HTTP/1.1 301 Moved Permanently\x0d\x0aLocation: http://www.google.com/\x0d\x0aContent-Type: text/html; charset=UTF-8\x0d\x0aDate: Sun, 15 Apr 2018 07:54:03 GMT\x0d\x0aExpires: Tue, 15 May 2018 07:54:03 GMT\x0d\x0aCache-Control: public, max-age=2592000\x0d\x0aServer: gws\x0d\x0aContent-Length: 219\x0d\x0aX-XSS-Protection: 1; mode=block\x0d\x0aX-Frame-Options: SAMEORIGIN\x0d\x0a\x0d\x0a<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">\x0a<TITLE>301 Moved</TITLE></HEAD><BODY>\x0a<H1>301 Moved</H1>\x0aThe document has moved\x0a<A HREF="http://www.google.com/">here</A>.\x0d\x0a</BODY></HTML>\x0d\x0a
2018/04/15 15:54:03 [info] 127.0.0.1:51650 <==> [2404:6800:4005:80e::200e]:80 (32 transmitted,540 received)
```
__[CLIENT]__
```bash
$ nc localhost 8800
GET / HTTP/1.1
HOST:google.com

HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Sun, 15 Apr 2018 07:54:03 GMT
Expires: Tue, 15 May 2018 07:54:03 GMT
Cache-Control: public, max-age=2592000
Server: gws
Content-Length: 219
X-XSS-Protection: 1; mode=block
X-Frame-Options: SAMEORIGIN

<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>301 Moved</TITLE></HEAD><BODY>
<H1>301 Moved</H1>
The document has moved
<A HREF="http://www.google.com/">here</A>.
</BODY></HTML>
^C
```

## LICENSE
This project is licensed under the terms of the MIT license.
