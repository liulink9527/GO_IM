package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessager 监听消息 当Server广播管道中有消息时 发送给所有用户的管道
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMessage := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMessage
}

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn)
	// 记录上线用户
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	// 用户上线广播
	s.BroadCast(user, "已上线")
	// 当前handler阻塞
	select {}
}

func (s *Server) Start() {
	fmt.Println("服务器开始监听.............")
	// 监听socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	//启动监听广播消息的协程
	go s.ListenMessager()

	for {
		// 监听客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.accept err:", err)
			continue
		}

		// 处理连接
		go s.Handler(conn)
	}
}
