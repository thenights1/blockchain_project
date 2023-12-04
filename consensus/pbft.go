package consensus

import (
	. "../data"
	. "../network"
	"fmt"
	"sync"
	"time"
)

// PBFT 实现了 PBFT 共识算法
type PBFT struct {
	nodes           []*Node
	prePrepareChan  chan *Block
	prepareChan     chan *Block
	commitChan      chan *Block
	committedBlocks map[int]*Block
	mutex           sync.Mutex
}

// NewPBFT 创建一个新的 PBFT 实例
func NewPBFT(nodes ...*Node) *PBFT {
	return &PBFT{
		nodes:           nodes,
		prePrepareChan:  make(chan *Block),
		prepareChan:     make(chan *Block),
		commitChan:      make(chan *Block),
		committedBlocks: make(map[int]*Block),
	}
}

// StartConsensus 启动 PBFT 共识
func (p *PBFT) StartConsensus() {
	fmt.Println("PBFT consensus started")

	// Implement PBFT consensus initialization
	// ...

	// Simulate some consensus steps
	go p.prePrepare()
	go p.prepare()
	go p.commit()
}

// prePrepare 模拟 PBFT 中的预准备步骤
func (p *PBFT) prePrepare() {
	for {
		select {
		case block := <-p.prePrepareChan:
			fmt.Println("Received pre-prepare for block", block.BlockNumber)
			// Implement pre-prepare logic
			// ...

			// Broadcast prepare to other nodes
			go p.broadcastPrepare(block)
		}
	}
}

// prepare 模拟 PBFT 中的准备步骤
func (p *PBFT) prepare() {
	for {
		select {
		case block := <-p.prepareChan:
			fmt.Println("Received prepare for block", block.BlockNumber)
			// Implement prepare logic
			// ...

			// Broadcast commit to other nodes
			go p.broadcastCommit(block)
		}
	}
}

// commit 模拟 PBFT 中的提交步骤
func (p *PBFT) commit() {
	for {
		select {
		case block := <-p.commitChan:
			fmt.Println("Received commit for block", block.BlockNumber)
			// Implement commit logic
			// ...

			// Add the block to the committed blocks map
			p.mutex.Lock()
			p.committedBlocks[block.BlockNumber] = block
			p.mutex.Unlock()
		}
	}
}

// broadcastPrepare 模拟广播准备消息给其他节点
func (p *PBFT) broadcastPrepare(block *Block) {
	for _, node := range p.nodes {
		if node != block.Proposer {
			go func(n *Node) {
				// Simulate network delay
				time.Sleep(time.Millisecond * 100)
				n.ReceivePrepare(block)
			}(node)
		}
	}
}

// broadcastCommit 模拟广播提交消息给其他节点
func (p *PBFT) broadcastCommit(block *Block) {
	for _, node := range p.nodes {
		if node != block.Proposer {
			go func(n *Node) {
				// Simulate network delay
				time.Sleep(time.Millisecond * 100)
				n.ReceiveCommit(block)
			}(node)
		}
	}
}

// 其他 PBFT 相关方法和结构体可以在这里添加
