package test

import (
	"blockchain/data"
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	nodeaddr := "8888"
	id := "NO.1"
	node := data.NewNode(nodeaddr, id)
	fmt.Println(node)
	fmt.Println(node.ID)
	fmt.Println(node.Blockchain)
	//t.Fatalf("aaa")
}

//func TestNodeActivity(t *testing.T) {
//	// 创建一个新节点
//	nodeaddr := 8888
//	node := data.NewNode(nodeaddr)
//
//	// 启动节点
//	go node.Start()
//
//	// 模拟一些节点活动
//	node.AddTransaction("testTransaction")
//	node.SubmitBlockForConsensus()
//	node.GetTransactionInfo("testTransaction")
//	node.GetBlockInfo(1)
//	node.SyncNodes()
//
//	// 等待一段时间，让节点完成周期性活动
//	time.Sleep(time.Second * 2)
//
//	// 停止节点
//	node.Stop()
//
//	// 在实际测试中，你可能需要验证节点的输出是否符合预期
//	// 以及进行更详细的测试
//
//	// 这里演示一些基本的断言，你可以根据实际情况扩展测试
//	if len(node.TransactionPool) != 1 {
//		t.Errorf("Expected transaction pool length to be 1, got %d", len(node.TransactionPool))
//	}
//
//	if len(node.Blockchain.Blocks) != 1 {
//		t.Errorf("Expected blockchain length to be 1, got %d", len(node.Blockchain.Blocks))
//	}
//
//	if !node.Synchronized {
//		t.Error("Expected node to be synchronized, but it's not")
//	}
//
//	t.Logf("Node %s activity test completed", nodeID)
//}
