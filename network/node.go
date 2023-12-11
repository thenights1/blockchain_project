// network/node.go

package network

import (
	. "blockchain/consensus"
	. "blockchain/data"
	"fmt"
	"sync"
	"time"
)

// Node 节点的数据结构
type Node struct {
	ID              string        //node的唯一标识符
	TransactionPool []string      //存储待处理交易的池。
	Blockchain      *Blockchain   //表示节点拥有的区块链
	PublicKey       string        //存储公钥
	PrivateKey      string        //存储私钥
	Consensus       *PBFT         //表示节点使用的共识机制，这里是PBFT。
	mutex           sync.Mutex    //用于在多协程中保护数据一致性的互斥锁
	Synchronized    bool          //用于判断节点是否与其他节点同步
	stopChan        chan struct{} //用于通知节点停止的通道
}

// NewNode 创建一个新的节点实例
func NewNode(id string) *Node {
	return &Node{
		ID:              id,
		TransactionPool: make([]string, 0),
		Blockchain:      NewBlockchain(),
		Consensus:       NewPBFT(),
		stopChan:        make(chan struct{}),
	}
}

// Start 启动节点
func (n *Node) Start() {
	fmt.Printf("Node %s started\n", n.ID)

	for {
		select {
		case <-n.stopChan:
			fmt.Printf("Node %s stopped\n", n.ID)
			return
		default:
			// 模拟节点的周期性活动
			n.periodicActivity()
		}
	}
}

// Stop 停止节点
func (n *Node) Stop() {
	close(n.stopChan)
}

// AddTransaction 添加新交易到待处理池
func (n *Node) AddTransaction(transaction string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s added transaction: %s\n", n.ID, transaction)
	n.TransactionPool = append(n.TransactionPool, transaction)
}

// SubmitBlockForConsensus 提交区块到共识过程
func (n *Node) SubmitBlockForConsensus() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s submitted block for consensus\n", n.ID)
	n.Consensus.StartConsensus()
}

// GetTransactionInfo 获取特定交易的信息
func (n *Node) GetTransactionInfo(transaction string) string {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s retrieved transaction info for: %s\n", n.ID, transaction)
	// 实际情况中，你可能会根据交易 ID 从数据库或区块链中检索相关信息
	return fmt.Sprintf("Transaction Info for %s from Node %s", transaction, n.ID)
}

// GetBlockInfo 获取特定区块的信息
func (n *Node) GetBlockInfo(blockNumber int) string {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s retrieved block info for block %d\n", n.ID, blockNumber)
	// 实际情况中，你可能会根据区块编号从数据库或区块链中检索相关信息
	return fmt.Sprintf("Block Info for Block %d from Node %s", blockNumber, n.ID)
}

// SyncNodes 执行节点之间的同步操作
func (n *Node) SyncNodes() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s syncing nodes\n", n.ID)
	// 实际情况中，你可能会与其他节点进行数据同步操作
}

// periodicActivity 模拟节点的周期性活动
func (n *Node) periodicActivity() {
	// 模拟节点的周期性活动，例如清理过期交易、定期提交区块等
	time.Sleep(time.Second)
}
