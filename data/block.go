// data/block.go

package data

import (
	. "../network"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	block.Hash = block.calculateHash()
	return block
}

// calculateHash 计算区块的哈希值
func (b *Block) calculateHash() string {
	data := fmt.Sprintf("%d%s%s%s%s", b.BlockNumber, b.Transactions, b.Timestamp, b.PrevHash, b.Proposer)
	hashInBytes := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashInBytes[:])
}
