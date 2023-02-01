package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
)

func Master(l net.Conn) {
	defer func() {
		l.Close()
	}()
	//读取服务器返回的数据
	go func() {
		buf := make([]byte, 1024)
		for {
			l, err := l.Read(buf)
			if err != nil {
				fmt.Print(err)
				return
			}
			if l > 0 {
				fmt.Print("s:")
				fmt.Println(string(buf[:l]))
				fmt.Print("c:")
			}
		}
	}()
	//读取标准输入写入到server
	fmt.Print("c:")
	for {
		buf := make([]byte, 1024)
		len, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		if len > 0 {
			l.Write(buf[:len])
		}

	}
}

func main() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	l, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	//启动master处理 监听器
	go Master(l)
	<-s
	fmt.Printf("client end===>\r\n")
}
