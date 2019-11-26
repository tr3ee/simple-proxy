# simple-proxy

a simple proxy POC based on my [go-link](https://github.com/tr3ee/go-link)

## Install

```bash
$ go get github.com/tr3ee/simple-proxy
```
this command should install `simple-proxy` in your $GOPATH/bin directory

## Usage
```
Usage of simple-proxy:
  -d    decryption mode
  -k string
        secret key for cipher
  -l string
        local address to listen on
  -ln string
        local network protocal will be used when listening (default "tcp")
  -m string
        cipher method (currently support: plain|xor)
  -r string
        remote address to connect
  -rn string
        remote network protocal will be used when connecting (default "tcp")
  -v    verbose mode
```

## TEST

__[SERVER]__
```bash
$ simple-proxy -r google.com:80 -l 127.0.0.1:8800 -m plain -v
2019/11/26 21:20:28 [INFO] listening on tcp:127.0.0.1:8800
[SEND] GET / HTTP/1.1
[SEND] Host: google.com
[SEND] 
[RECV] HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Tue, 26 Nov 2019 13:20:39 GMT
Expires: Thu, 26 Dec 2019 13:20:39 GMT
Cache-Control: public, max-age=2592000
Server: gws
Content-Length: 219
X-XSS-Protection: 0
X-Frame-Options: SAMEORIGIN

<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>301 Moved</TITLE></HEAD><BODY>
<H1>301 Moved</H1>
The document has moved
<A HREF="http://www.google.com/">here</A>.
</BODY></HTML>
2019/11/26 21:20:42 [INFO] 127.0.0.1:52723 <==> 127.0.0.1:52724 (33 transmitted, 528 received)
```
__[CLIENT]__
```bash
$ nc localhost 8800
GET / HTTP/1.1
Host: google.com

HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Tue, 26 Nov 2019 13:20:39 GMT
Expires: Thu, 26 Dec 2019 13:20:39 GMT
Cache-Control: public, max-age=2592000
Server: gws
Content-Length: 219
X-XSS-Protection: 0
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
