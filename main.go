package main

import (
	"blockchain/data"
	"fmt"
	"log"
	"os"
	"sync"
)

// 节点池，主要用来存储监听地址
//var NodeTable map[string]string
//var Users []*data.Client

func main() {
	// 创建两个节点
	//node1 := data.NewNode("127.0.0.1:9001", "Node1")
	//node2 := data.NewNode("127.0.0.1:9002", "Node2")
	//node3 := data.NewNode("127.0.0.1:9003", "Node3")
	//node4 := data.NewNode("127.0.0.1:9004", "Node4")
	//node5 := data.NewNode("127.0.0.1:9005", "Node5")

	data.NodeTable = map[string]string{
		"Node0": "127.0.0.1:9000",
		"Node1": "127.0.0.1:9001",
		"Node2": "127.0.0.1:9002",
		"Node3": "127.0.0.1:9003",
	}
	var bc *data.Blockchain
	No_usenode := data.NewNode("ddd", "ddd", bc) //虚假结点，用来当做创世区块的提出者
	blockchain := data.NewBlockchain(No_usenode)
	//blockchain.SaveBlockchainToJSON("./kk.json")
	//data.ClientTable = map[string]string{
	//	"User1": "0x145287",
	//	"User2": "0x124563",
	//	"User3": "0x145235",
	//	"User4": "0x147889",
	//}

	// 创建 PBFT 共识实例
	//pbft := data.NewPBFT(node1, node2)
	if len(os.Args) != 2 {
		log.Panic("输入的参数有误！")
	}
	nodeID := os.Args[1]

	var wg sync.WaitGroup
	wg.Add(1)

	if nodeID == "client" {
		data.ClientTcpListen("127.0.0.1:8001")
	} else if addr, ok := data.NodeTable[nodeID]; ok {
		node := data.NewNode(addr, nodeID, blockchain)
		if node.ID == "Node0" {
			node.Consensus = data.NewPBFT(node, true)
			node.View = 3
		} else {
			node.Consensus = data.NewPBFT(node, false)
			node.View = 0
		}
		err := data.NodeSaveKeysToFile(node)
		if err != nil {
			fmt.Println("Error saving keys for node %s: %v", node.ID, err)
		}
		go node.Start()
	}
	//if nodeID == "Node1" {
	//	go node1.Start()
	//	// 等待节点启动完成
	//	time.Sleep(time.Second)
	//	//ccc
	//} else {
	//	go node2.Start()
	//	// 等待节点启动完成
	//	time.Sleep(time.Second)
	//	go func() {
	//		defer wg.Done()
	//		context := []byte("Hello from Node2")
	//		addr := "localhost:9001"
	//		fmt.Printf("[%s] Sending message to %s\n", node2.ID, addr)
	//		node2.Sendmessage(context, addr)
	//		node2.Stop()
	//		time.Sleep(time.Second)
	//	}()
	//
	//}
	// 启动节点

	// 启动 PBFT 共识
	//pbft.StartConsensus()

	// 等待共识完成
	//time.Sleep(3 * time.Second)

	// 向节点发送消息

	wg.Wait()
}
