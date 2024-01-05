// network/rpc.go

package data

import (
	"fmt"
	"io"
	"log"
	"net"
)

func ClientTcpListen(clientAddr string) {
	listen, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)
	defer listen.Close()

	RunClients()

}

// 节点使用的tcp监听
func (p *Node) tcpListen() {
	listen, err := net.Listen("tcp", p.Addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("节点开启监听，地址：%s\n", p.Addr)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := io.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		p.HandleRequest(string(b))

	}

}

// 使用tcp发送消息
func Sendmessage(context []byte, addr string) {
	//准备连接端口
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}
	//发送信息
	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}

	conn.Close()
}
