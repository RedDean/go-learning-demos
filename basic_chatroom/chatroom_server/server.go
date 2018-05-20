package main

import (
	"net"
	"fmt"
	"strings"
)

type Server struct {
	OnlineClients map[string]net.Conn  // 在线客户端映射
	MessageQueue  chan string		   // 消息队列
	QuitQueue     chan bool		   // 退出消息队列
}

/**
	@param : conn net.Conn
	功能: 接受客户端发来的消息
 */
func (s *Server) MsgAccepter(conn net.Conn) {
	// 读取消息
	input := make([]byte, 1024)
	defer func(conn net.Conn) {
		// 关闭连接，同时从在线映射中删除连接
		addr := conn.RemoteAddr().String()
		delete(s.OnlineClients, addr)
		fmt.Printf("client %s is disconnecting ... \n", addr)
		conn.Close()
	}(conn)

	for {
		bytesnum,err := conn.Read(input)
		if err!=nil {
			fmt.Printf("Error was occured while reading : %s\n", err.Error())
			break
		}

		if bytesnum != 0 {
		   msg := string(input[:bytesnum])
			if strings.ToUpper(msg) == "QUIT" {
				s.QuitQueue <- true
			}
		   s.MessageQueue <- msg
		}
	}
}

/**
	@param : conn net.Conn
	功能: 处理消息消息队列
 */
func (s *Server) HandleMsgQueue() {
	for {
		select {
			case msg := <- s.MessageQueue :
				 s.HandleMsg(msg)
			case <- s.QuitQueue :
			break
		}
	}
}

/**
	@param : msg string
	功能: 处理消息
    消息格式 : 地址+"#"+内容 IP#xxxxxxxx
 */
func (s *Server) HandleMsg(Message string) {
	ctx := strings.Split(Message, "#")
	if len(ctx) > 1 {
		addr := strings.Trim(ctx[0]," ")
		msg  := strings.Join(ctx[1:], "#")  // 防止消息内容中带有"#"

		if conn,ok := s.OnlineClients[addr]; ok {
			_,err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Printf("error:%s , while writing msg to %",err.Error(), addr)
			}
		}
	}
}

func main() {

	o := make(map[string]net.Conn)  // 在线客户端映射
	m := make(chan string , 1024) // 消息队列
	q := make(chan bool)  // 退出消息队列
	server := &Server{o,m,q}

	// 监听本地8080端口
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	HandleError(err)
	defer listener.Close()

	go server.HandleMsgQueue()

	fmt.Printf("App is listening at port 8080...\n")
	for {
		conn,err := listener.Accept()
		HandleError(err)
		addr := conn.RemoteAddr().String()
		fmt.Printf("%s is connecting to server!", addr)
		server.OnlineClients[addr] = conn
		go server.MsgAccepter(conn)
	}
}

