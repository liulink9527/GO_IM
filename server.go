package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()
	fmt.Println("连接建立成功 来自:", conn.RemoteAddr())
}

func (s *Server) Start() {
	fmt.Println("服务器开始监听.............")
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.accept err:", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}

}
