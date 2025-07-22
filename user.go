package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser 创建一个用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听
	go user.ListenMessage()

	return user
}

// ListenMessage 监听消息
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn.Write err:", err)
		}
	}
}
