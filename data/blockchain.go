// data/blockchain.go

package data

import (
	"sync"
)

// Block 区块的数据结构

// Blockchain 区块链的数据结构
type Blockchain struct {
	chain  []*Block
	mutex  sync.RWMutex
	height int
}

// NewBlockchain 创建一个新的区块链实例
func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		BlockNumber: 0,
		Proposer:    nil,
	}
	return &Blockchain{
		chain:  []*Block{genesisBlock},
		mutex:  sync.RWMutex{},
		height: 1,
	}
}

// AddBlock 将新区块添加到区块链
func (bc *Blockchain) AddBlock(data string, proposer *Node) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	block := &Block{
		BlockNumber: bc.height,
		Proposer:    proposer,
	}

	bc.chain = append(bc.chain, block)
	bc.height++
}

// GetBlockInfo 获取指定区块的信息
func (bc *Blockchain) GetBlockInfo(blockNumber int) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if blockNumber > 0 && blockNumber < bc.height {
		return bc.chain[blockNumber]
	}

	return nil
}

// GetBlockchainHeight 获取区块链的高度
func (bc *Blockchain) GetBlockchainHeight() int {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.height
}
