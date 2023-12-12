package blockchain

import (
	"blockchain/data"
	"fmt"
)

func main() {
	// 初始化节点
	nodeA := data.NewNode("NodeA")
	nodeB := data.NewNode("NodeB")
	// ...

	// 启动节点
	go nodeA.Start()
	go nodeB.Start()
	// ...

	// 初始化共识
	pbft := data.NewPBFT(nodeA, nodeB /* add other nodes */)

	// 启动共识
	pbft.StartConsensus()

	// 添加交易到待处理池
	nodeA.AddTransaction("Transaction1")
	// ...

	// 提交区块到共识过程
	nodeA.SubmitBlockForConsensus()

	// 获取特定交易或区块的信息
	txInfo := nodeA.GetTransactionInfo("Transaction1")
	fmt.Println(txInfo)

	blockInfo := nodeA.GetBlockInfo(1)
	fmt.Println(blockInfo)

	// 执行节点之间的同步操作
	nodeA.SyncNodes()
	// ...

	// 等待节点停止
	nodeA.Stop()
	nodeB.Stop()
	// ...

}
