package main

import (
	"net"
	"bufio"
	"os"
	"strings"
	"fmt"
	"sync"
)
/*
type Cilent stuct {

}
*/
func Process(conn net.Conn, wg *sync.WaitGroup)  {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		// 输入 "QUIT" 退出
		if strings.ToUpper(input) == "QUIT" {
			conn.Write([]byte(input))
			conn.Close()
			break
		}

	    // 消息格式 : 地址+"#"+内容 IP#xxxxxxxx
		_, err := conn.Write([]byte(input))
		if err!=nil {
			conn.Close()
			fmt.Printf("error while connecting : %s", err.Error())
			break
		}
	}
	wg.Done()
}

func HandleReceived (conn net.Conn, wg *sync.WaitGroup){
	remotemsg := make([]byte, 512)
	for {
		_,err := conn.Read(remotemsg)
		if err!=nil {
			fmt.Println("客户端已退出！")
			break
		}
		fmt.Printf("Remote msg: %s", remotemsg)
	}
	wg.Done()
}

func main()  {
	var wg sync.WaitGroup
	wg.Add(2)

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	HandleError(err)
	defer conn.Close()

	go Process(conn, &wg)   //发消息
	go HandleReceived(conn, &wg) //读消息

	wg.Wait()
}
