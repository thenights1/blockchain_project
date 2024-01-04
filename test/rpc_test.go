package test

import (
	"blockchain/data"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRpctest(t *testing.T) {
	// 创建两个节点
	node1 := data.NewNode("Node1", "localhost:9001")
	node2 := data.NewNode("Node2", "localhost:9002")

	// 创建 PBFT 共识实例
	pbft := data.NewPBFT(node1, node2)

	// 启动节点
	go node1.Start()
	go node2.Start()

	// 等待节点启动完成
	time.Sleep(time.Second)

	// 启动 PBFT 共识
	pbft.StartConsensus()

	// 等待共识完成
	time.Sleep(3 * time.Second)

	// 向节点发送消息
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		context := []byte("Hello from Node1")
		addr := "localhost:9002"
		fmt.Printf("[%s] Sending message to %s\n", node1.ID, addr)
		data.Sendmessage(context, addr)
	}()

	go func() {
		defer wg.Done()
		context := []byte("Hello from Node2")
		addr := "localhost:9001"
		fmt.Printf("[%s] Sending message to %s\n", node2.ID, addr)
		data.Sendmessage(context, addr)
	}()

	wg.Wait()
}
