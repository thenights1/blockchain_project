// network/node.go

package data

import (
	"blockchain/crypto11"
	"crypto/ecdsa"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Node 节点的数据结构
type Node struct {
	ID              string            //node的唯一标识符
	TransactionPool []Transaction     //存储待处理交易的池。
	Blockchain      *Blockchain       //表示节点拥有的区块链
	PublicKey       *ecdsa.PublicKey  //存储公钥
	PrivateKey      *ecdsa.PrivateKey //存储私钥
	//Consensus       *PBFT             //表示节点使用的共识机制，这里是PBFT。
	mutex        sync.Mutex    //用于在多协程中保护数据一致性的互斥锁
	Synchronized bool          //用于判断节点是否与其他节点同步
	stopChan     chan struct{} //用于通知节点停止的通道
	Addr         string
	view         int
}

// NewNode 创建一个新的节点实例
func NewNode(addr string, id string) *Node {

	privateKey, publicKey, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Printf("Error generating key pair: %v", err)
	}
	node := &Node{
		ID:              id,
		TransactionPool: make([]Transaction, 0),
		Blockchain:      NewBlockchain(),
		stopChan:        make(chan struct{}),
		Synchronized:    false,
		Addr:            addr,
		PrivateKey:      privateKey,
		PublicKey:       publicKey,
	}

	return node
}

// Start 启动节点
func (n *Node) Start() {
	fmt.Printf("Node %s started\n", n.ID)
	flag := 0
	for {
		select {
		case <-n.stopChan:
			fmt.Printf("Node %s stopped\n", n.ID)
			return
		default:
			if flag != 1 {
				go n.tcpListen()
				flag = 1
			}
			// 模拟节点的周期性活动
			//time.Sleep(3*time.Second)
			//n.PackTransactions(*TransactionPool)
		}
	}
	//n.tcpListen()
}

// Stop 停止节点
func (n *Node) Stop() {
	close(n.stopChan)
}

func (n *Node) HandleRequest(request string) {
	flag := request[:4]
	if flag == "tran" {
		transaction_temp := &Transaction{}
		err := transaction_temp.FromJSON(request)
		if err != nil {
			fmt.Println("Error restoring transaction from JSON in node.go: %v", err)
		}
		n.AddTransaction(*transaction_temp)
	}

}

// AddTransaction 添加新交易到待处理池
func (n *Node) AddTransaction(transaction Transaction) {
	// 使用互斥锁，确保在多线程环境下的安全并发操作
	n.mutex.Lock()
	// 函数结束后解锁
	defer n.mutex.Unlock()

	fmt.Printf("Node %s added transaction: %s\n", n.ID, transaction)
	n.TransactionPool = append(n.TransactionPool, transaction)
}

// PackTransactions 将待处理池中的交易打包成区块
func (n *Node) PackTransactions(transactionPool []*Transaction) ([]*Transaction, *Block) {
	const maxTransactionsPerBlock = 3
	const maxFeeThreshold = 1.0

	// 按手续费大小排序
	SortTransactionsByFee(transactionPool)

	var selectedTransactions []*Transaction
	totalFee := 0.0

	for _, transaction := range transactionPool {
		if len(selectedTransactions) >= maxTransactionsPerBlock {
			break
		}

		// 如果加入该交易后手续费总额不超过阈值，则加入选中的交易列表
		if totalFee+transaction.Premium <= maxFeeThreshold {
			selectedTransactions = append(selectedTransactions, transaction)
			totalFee += transaction.Premium
		}
	}

	// 从待处理池中移除已选中的交易
	n.mutex.Lock()
	remainingTransactions := RemoveSelectedTransactions(transactionPool, selectedTransactions)
	n.mutex.Unlock()

	// 打包交易成区块
	blockNumber := 0 //暂时变为0，确定加入区块链后再计算值
	prehash := "xxx" //暂定值，确定加入区块链后再计算值
	newBlock := NewBlock(blockNumber, selectedTransactions, prehash, n)

	// 返回剩余的交易和新创建的区块
	return remainingTransactions, newBlock
}

// SortTransactionsByFee 按手续费大小对交易进行排序
func SortTransactionsByFee(transactions []*Transaction) {
	// 实现交易按手续费排序的逻辑，可以使用标准库的排序方法
	// 这里简化为按手续费降序排序
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Premium > transactions[j].Premium
	})
}

// RemoveSelectedTransactions 从交易池中移除已选中的交易
func RemoveSelectedTransactions(transactionPool, selectedTransactions []*Transaction) []*Transaction {

	var remainingTransactions []*Transaction

	for _, transaction := range transactionPool {
		if !contains(selectedTransactions, transaction) {
			remainingTransactions = append(remainingTransactions, transaction)
		}
	}

	return remainingTransactions
}

// contains 检查交易是否包含在已选中的交易列表中
func contains(transactions []*Transaction, target *Transaction) bool {
	for _, transaction := range transactions {
		if transaction == target {
			return true
		}
	}
	return false
}

// SubmitBlockForConsensus 提交区块到共识过程
func (n *Node) SubmitBlockForConsensus() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	fmt.Printf("Node %s submitted block for consensus\n", n.ID)
	//n.Consensus.StartConsensus()
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
