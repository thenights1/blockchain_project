// consensus/pbft.go

package data

import (
	"fmt"
	"sync"
)

// PBFT 实现了 PBFT 共识算法
type PBFT struct {
	nodeID          string
	node            *Node
	prePrepareChan  chan *Block
	prepareChan     chan *Block
	commitChan      chan *Block
	committedBlocks map[int]*Block
	ifPropser       bool
	mutex           sync.Mutex
}

// NewPBFT 创建一个新的 PBFT 实例
func NewPBFT(n *Node, ifpropser bool) *PBFT {
	return &PBFT{
		node:            NewNode(n.Addr, n.ID),
		prePrepareChan:  make(chan *Block),
		prepareChan:     make(chan *Block),
		commitChan:      make(chan *Block),
		committedBlocks: make(map[int]*Block),
		ifPropser:       ifpropser,
	}
}

// StartConsensus 启动 PBFT 共识
//func (p *PBFT) StartConsensus() {
//	fmt.Println("PBFT consensus started")
//
//	go p.prePrepare()
//	go p.prepare()
//	go p.commit()
//}

// prePrepare 模拟 PBFT 中的预准备步骤
//func (p *PBFT) prePrepare() {
//	for {
//		select {
//		case block := <-p.prePrepareChan:
//			fmt.Println("Received pre-prepare for block", block.BlockNumber)
//			// Implement pre-prepare logic
//			// ...
//
//			// Broadcast prepare to other nodes
//			go p.broadcastPrepare(block)
//		}
//	}
//}
//
//// prepare 模拟 PBFT 中的准备步骤
//func (p *PBFT) prepare() {
//	for {
//		select {
//		case block := <-p.prepareChan:
//			fmt.Println("Received prepare for block", block.BlockNumber)
//			// Implement prepare logic
//			// ...
//
//			// Broadcast commit to other nodes
//			go p.broadcastCommit(block)
//		}
//	}
//}

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
//func (p *PBFT) broadcastPrepare(block *Block) {
//	for _, node := range NodeTable {
//		if p.ifPropser == true {
//			go func(n *Node) {
//				// Simulate network delay
//				time.Sleep(time.Millisecond * 100)
//				//n.ReceivePrepare(block)
//			}()
//		}
//	}
//}

// broadcastCommit 模拟广播提交消息给其他节点
//func (p *PBFT) broadcastCommit(block *Block) {
//	for _, node := range p.nodes {
//		if node != block.Proposer {
//			go func(n *Node) {
//				// Simulate network delay
//				time.Sleep(time.Millisecond * 100)
//				//n.ReceiveCommit(block)
//			}(node)
//		}
//	}
//}

// 主节点提出新区块的函数
func (p *PBFT) proposeNewBlock() {
	// 此处省略具体逻辑，可根据具体需求实现
	// 生成新的区块，将其广播到网络上
	// 此处需要注意轮换主节点的机制
}

// 提交交易到共识过程
func (p *PBFT) SubmitTransaction(transaction *Transaction) {
	// 此处省略具体逻辑，可根据具体需求实现
	// 将交易加入待处理池，并触发共识过程
}

// 主节点轮换机制
func (p *PBFT) rotatePrimary() {
	// 此处省略具体逻辑，可根据具体需求实现
	// 根据轮换算法选择下一个主节点
}

// 外部接口，模拟客户端向节点提交交易
func (p *PBFT) ClientSubmitTransaction(transaction *Transaction) {
	// 此处省略具体逻辑，可根据具体需求实现
	// 客户端将交易提交到某个节点
	//node := p.nodes[0] // 简单示例，选择第一个节点
	//node.ReceiveTransaction(transaction)
}

// ReceiveTransaction 接收来自网络的交易
func (p *PBFT) ReceiveTransaction(transaction *Transaction) {
	// 此处省略具体逻辑，可根据具体需求实现
	// 节点接收到来自网络的交易，并加入待处理池
	// 触发共识过程
	p.SubmitTransaction(transaction)
}

// ReceivePrepare 接收来自网络的准备消息
func (p *PBFT) ReceivePrepare(block *Block) {
	// 此处省略具体逻辑，可根据具体需求实现
	// 节点接收到来自网络的准备消息，并进行相应处理
}

// ReceiveCommit 接收来自网络的提交消息
func (p *PBFT) ReceiveCommit(block *Block) {
	// 此处省略具体逻辑，可根据具体需求实现
	// 节点接收到来自网络的提交消息，并进行相应处理
}

// 外部接口，模拟客户端查询区块信息
//func (p *PBFT) ClientQueryBlock(blockNumber int) *Block {
//	// 此处省略具体逻辑，可根据具体需求实现
//	// 客户端查询指定区块信息
//	// 可以选择任意一个节点进行查询
//	node := p.nodes[0] // 简单示例，选择第一个节点
//	return node.QueryBlock(blockNumber)
//}

// QueryBlock 查询区块信息
func (p *PBFT) QueryBlock(blockNumber int) *Block {
	// 此处省略具体逻辑，可根据具体需求实现
	// 节点查询指定区块信息
	// 可以从已经提交的区块中查找
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.committedBlocks[blockNumber]
}
