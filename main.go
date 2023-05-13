package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Printf("Trojan-killer v1.0.0 started\n")
	l, _ := net.Listen("tcp4", "127.0.0.1:12345")
	fmt.Printf("Listening on %v\n\n", l.Addr())
	for {
		c, _ := l.Accept()
		go Handle(c)
	}
}

var CCS = []byte{20, 3, 3, 0, 1, 1}

func Handle(c net.Conn) {
	req, err := http.ReadRequest(bufio.NewReader(c))
	if err != nil {
		return
	}
	state := "accepted"
	if !strings.EqualFold(req.Method, "CONNECT") {
		state = "rejected"
	}
	fmt.Printf("%v from %v %v %v\n", time.Now().Format(time.DateTime), c.RemoteAddr(), state, req.URL.Host)
	if state == "rejected" {
		return
	}

	conn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		return
	}
	c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	var mutex sync.Mutex

	uploading := false
	upCount := 0

	downloading := false
	downCount := 0

	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := c.Read(buf)
			if err != nil {
				return
			}
			mutex.Lock()
			if upCount == 0 && n >= 6 && bytes.Equal(buf[:6], CCS) {
				uploading = true
			}
			if uploading {
				upCount += n
			}
			if downloading {
				downloading = false
				//fmt.Printf("%v\tupCount %v\tdownCount %v\n", req.URL.Host, upCount, downCount)
				if upCount >= 650 && upCount <= 750 &&
					((downCount >= 170 && downCount <= 180) || (downCount >= 3000 && downCount <= 7500)) {
					fmt.Printf("%v is Trojan\n", req.URL.Host)
				}
			}
			mutex.Unlock()
			_, err = conn.Write(buf[:n])
			if err != nil {
				return
			}
			if !downloading && downCount != 0 {
				go io.CopyBuffer(conn, c, buf)
				return
			}
		}
	}()

	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			mutex.Lock()
			if uploading {
				uploading = false
				downloading = true
			}
			if downloading {
				downCount += n
			}
			mutex.Unlock()
			_, err = c.Write(buf[:n])
			if err != nil {
				return
			}
			if !downloading && downCount != 0 {
				go io.CopyBuffer(c, conn, buf)
				return
			}
		}
	}()
}
