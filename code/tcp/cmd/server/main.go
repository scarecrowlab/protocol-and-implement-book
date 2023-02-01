package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
)

func Master(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		//启动worker处理 客户端的每个连接
		go Worker(conn)
	}
}

func Worker(conn net.Conn) {
	defer func() {
		fmt.Printf("defer close")
		conn.Close()
	}()
	buf := make([]byte, 1024)
	for {
		l, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("read error")
			return
		}
		if l > 0 {
			l, err = conn.Write(buf[:l])
			if err != nil {
				fmt.Println("echo error")
			} else {
				fmt.Printf("c %s %s", conn.RemoteAddr(), buf[:l])
			}
			//TODO 处理粘包 和断包
		} else {
			fmt.Println("[warning] read zero")
		}
	}
}

func main() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	//启动master处理 监听器
	go Master(l)
	<-s
	fmt.Printf("end===>\r\n")
}
