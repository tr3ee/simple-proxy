# simple-proxy

a simple forward/active/passive proxy POC based on my [go-link](https://github.com/tr3ee/go-link)

## Install

```bash
$ go get github.com/tr3ee/simple-proxy
```
this command should install `simple-proxy` in your $GOPATH/bin directory

## Overview

- **Forward Mode**

```txt


        +--------+                +--------------+                  +--------+
        | Client |                | Forward Mode |                  | Server |
        +---+----+                +------+-------+                  +----+---+
            |                            |                               |
            |        Dial to Proxy       |                               |
            |+-------------------------->|                               |
            |         (connected)        |         Dial to Server        |
            |                            |+----------------------------->|
            |   Send message to Proxy    |          (connected)          |
            |+-------------------------->|                               |
            |                            |   Forward message to Server   |
            |                            |+----------------------------->|
            |                            |                               |
            |                            |     Send reply to Proxy       |
            |                            |<-----------------------------+|
            |   Forward reply to Client  |                               |
            |<--------------------------+|                               |
            |                            |                               |
            |          ......            |             ......            |
            |                            |                               |
        +---+----+                +------+-------+                  +----+---+
        | Client |                | Forward Mode |                  | Server |
        +--------+                +--------------+                  +--------+


```

- **Active Mode**

```txt


        +--------+                +-------------+                   +--------+
        | Client |                | Active Mode |                   | Server |
        +---+----+                +------+------+                   +----+---+
            |                            |                               |
            |        Dial to Client      |                               |
            |<--------------------------+|                               |
            |         (connected)        |         Dial to Server        |
            |                            |+----------------------------->|
            |   Send message to Proxy    |          (connected)          |
            |+-------------------------->|                               |
            |                            |   Forward message to Server   |
            |                            |+----------------------------->|
            |                            |                               |
            |                            |     Send reply to Proxy       |
            |                            |<-----------------------------+|
            |   Forward reply to Client  |                               |
            |<--------------------------+|                               |
            |                            |                               |
            |          ......            |             ......            |
            |                            |                               |
        +---+----+                +------+------+                   +----+---+
        | Client |                | Active Mode |                   | Server |
        +--------+                +-------------+                   +--------+


```

  

- **Passive Mode**

```txt


        +--------+                +--------------+                  +--------+
        | Client |                | Passive Mode |                  | Server |
        +---+----+                +------+-------+                  +----+---+
            |                            |                               |
            |        Dial to Proxy       |                               |
            |+-------------------------->|                               |
            |         (connected)        |         Dial to Proxy         |
            |                            |<-----------------------------+|
            |   Send message to Proxy    |          (connected)          |
            |+-------------------------->|                               |
            |                            |   Forward message to Server   |
            |                            |+----------------------------->|
            |                            |                               |
            |                            |     Send reply to Proxy       |
            |                            |<-----------------------------+|
            |   Forward reply to Client  |                               |
            |<--------------------------+|                               |
            |                            |                               |
            |          ......            |             ......            |
            |                            |                               |
        +---+----+                +------+-------+                  +----+---+
        | Client |                | Passive Mode |                  | Server |
        +--------+                +--------------+                  +--------+


```

## Quick Start

__[SERVER]__
```bash
$ simple-proxy -r google.com:80 -l 127.0.0.1:8800 -v
2020/03/20 13:43:10 [+] listening on tcp:127.0.0.1:8800
GET / HTTP/1.1
Host: google.com

HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Fri, 20 Mar 2020 05:43:14 GMT
Expires: Sun, 19 Apr 2020 05:43:14 GMT
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
2020/03/20 13:43:24 [+] 127.0.0.1:56513 <==> 127.0.0.1:56514 (33 transmitted, 528 received)
```
__[CLIENT]__
```bash
$ nc localhost 8800
GET / HTTP/1.1
Host: google.com

HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Fri, 20 Mar 2020 05:43:14 GMT
Expires: Sun, 19 Apr 2020 05:43:14 GMT
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

## Usage

```
Usage of simple-proxy:
  -d    decryption mode
  -hex
        print inbound/outbound data in hexadecimal format
  -k string
        secret key for cipher
  -l string
        local address to listen on
  -ln string
        local network protocal will be used when listening (default "tcp")
  -m string
        cipher method (currently support: plain|xor) (default "plain")
  -mode string
        proxy mode (currently support: active|forward|passive) (default "forward")
  -no-color
        disable color output
  -passive-qsize int
        max queue size for local listener in passive mode (qsize > 0) (default 32)
  -r string
        remote address to connect
  -rn string
        remote network protocal will be used when connecting (default "tcp")
  -t int
        idle timeout for each connection (t > 0) (default 30)
  -v    verbose mode
  -vv
        more verbose mode
```

## LICENSE
This project is licensed under the terms of the MIT license.
