// data/block.go

package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// Block 区块的数据结构
type Block struct {
	BlockNumber  int
	Transactions []*Transaction
	Timestamp    time.Time
	PrevHash     string
	Hash         string
	Proposer     *Node
}

// NewBlock 创建一个新的区块实例
func NewBlock(blockNumber int, transactions []*Transaction, prevHash string, proposer *Node) *Block {
	block := &Block{
		BlockNumber:  blockNumber,
		Transactions: transactions,
		Timestamp:    time.Now(),
		PrevHash:     prevHash,
		Proposer:     proposer,
	}
	block.Hash = block.CalculateHash()
	return block
}

// calculateHash 计算区块的哈希值
func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%s%s%s%s", b.BlockNumber, b.Transactions, b.Timestamp, b.PrevHash, b.Proposer)
	hashInBytes := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashInBytes[:])
}

//处理区块内的交易

// PackTransactions 将待处理池中的交易打包成区块
func (b *Block) PackTransactions(transactionPool []*Transaction) ([]*Transaction, *Block) {
	const maxTransactionsPerBlock = 10
	const maxFeeThreshold = 5.0

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
	remainingTransactions := RemoveSelectedTransactions(transactionPool, selectedTransactions)

	// 打包交易成区块
	blockNumber := b.BlockNumber + 1
	newBlock := NewBlock(blockNumber, selectedTransactions, b.Hash, b.Proposer)

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
