// consensus/pbft.go

package data

import (
	"blockchain/crypto11"
	"fmt"
	"math/rand"
	"sync"
)

type Message struct {
	ID      string
	Sign    string
	Content string
}

// PBFT 实现了 PBFT 共识算法
type PBFT struct {
	node                   *Node
	messagePool            map[string]Message //临时消息池，用id来索引消息
	prePrepareConfirmCount int                ////第一个string为消息id，第二个为结点名，用来代表该预准备消息已经被多少结点确认
	PrePareConfirmCount    int                //第一个string为消息id，第二个为结点名，用来代表该准备消息已经被多少结点确认
	commitConfirmCount     int                //第一个string为消息id，第二个为结点名，用来代表该提交消息已经被多少结点确认
	//该笔消息是否已进行Commit广播
	isCommitBordcast map[string]bool
	//该笔消息是否已对客户端进行Reply
	isReply   map[string]bool
	ifPropser bool
	mutex     sync.Mutex
}

// NewPBFT 创建一个新的 PBFT 实例
func NewPBFT(n *Node, ifpropser bool) *PBFT {
	return &PBFT{
		node:                   NewNode(n.Addr, n.ID),
		messagePool:            make(map[string]Message),
		prePrepareConfirmCount: 1,
		PrePareConfirmCount:    1,
		commitConfirmCount:     1,
		isCommitBordcast:       make(map[string]bool),
		isReply:                make(map[string]bool),
		ifPropser:              ifpropser,
	}
}
func GenerateMessageID() string {
	// 定义字符集
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机字符串
	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(result)
}

//StartConsensus 启动 PBFT 共识
//func (p *PBFT) StartConsensus() {
//	fmt.Println("PBFT consensus started")
//
//	go p.prePrepare()
//	go p.prepare()
//	go p.commit()
//}

// prePrepare 模拟 PBFT 中的预准备步骤
func (pbft *PBFT) PrePrepare(block *Block) {
	pbft.mutex.Lock()
	defer pbft.mutex.Unlock()

	if !pbft.ifPropser {
		fmt.Println("Error: Only proposer can send PrePrepare messages")
		return
	}

	message := Message{
		ID: GenerateMessageID(),
	}

	message.Sign, _ = crypto11.Sign(pbft.node.PrivateKey, []byte(block.Hash))
	message.Content = "preM" + pbft.node.ID + message.Sign + " " + block.Hash
	//fmt.Println(block.Hash)
	//fmt.Println(pbft.node.PublicKey)

	// 将消息加入消息池
	pbft.messagePool[message.ID] = message

	// 向其他节点发送预准备消息
	pbft.broadcast(message.Content)
}

//func (pbft *PBFT) Prepare(messageID string) {
//	pbft.mutex.Lock()
//	defer pbft.mutex.Unlock()
//
//	if pbft.ifPropser {
//		fmt.Println("Error: Proposer should not send Prepare messages")
//		return
//	}
//
//	// 检查消息是否存在于消息池中
//	_, ok := pbft.messagePool[messageID]
//	if !ok {
//		fmt.Println("Error: Message not found in the pool")
//		return
//	}
//
//	// 将节点对预准备消息的确认计数加一
//	pbft.prePrepareConfirmCount[messageID][pbft.node.ID] = true
//
//	// 如果满足某个条件（例如，超过一定数量的确认），则广播准备消息
//	if pbft.checkPrePareConfirm(messageID) {
//		pbft.broadcastPrepare(messageID)
//	}
//}

// Commit 发送提交消息
//func (pbft *PBFT) Commit(messageID string) {
//	pbft.mutex.Lock()
//	defer pbft.mutex.Unlock()
//
//	if !pbft.ifPropser {
//		fmt.Println("Error: Only proposer can send Commit messages")
//		return
//	}
//
//	// 检查消息是否存在于消息池中
//	_, ok := pbft.messagePool[messageID]
//	if !ok {
//		fmt.Println("Error: Message not found in the pool")
//		return
//	}
//
//	// 将节点对准备消息的确认计数加一
//	pbft.commitConfirmCount[messageID][pbft.node.ID] = true
//
//	// 如果满足某个条件（例如，超过一定数量的确认），则广播提交消息
//	if pbft.checkCommitConfirm(messageID) {
//		pbft.broadcastCommit(messageID)
//		// 这里可以添加处理提交消息后的逻辑，例如向客户端发送 Reply
//	}
//}

// broadcast 向其他节点广播消息
func (pbft *PBFT) broadcast(message string) {
	fmt.Printf("Node %s broadcasts message\n", pbft.node.ID)
	for id, addr := range NodeTable {
		if id == pbft.node.ID {
			continue
		}
		go Sendmessage([]byte(message), addr)
	}
}

//// broadcastPrepare 向其他节点广播准备消息
//func (pbft *PBFT) broadcastPrepare(messageID string) {
//	fmt.Printf("Node %s broadcasts Prepare for message %s\n", pbft.node.ID, messageID)
//	// 在实际应用中，这里应该向其他节点发送消息
//}
//
//// broadcastCommit 向其他节点广播提交消息
//func (pbft *PBFT) broadcastCommit(messageID string) {
//	fmt.Printf("Node %s broadcasts Commit for message %s\n", pbft.node.ID, messageID)
//
//	// 在实际应用中，这里应该向其他节点发送消息
//}

// checkPrePareConfirm 检查是否满足预准备消息的确认条件
//func (pbft *PBFT) checkPrePareConfirm(messageID string) bool {
//	// 在实际应用中，可以根据具体的确认条件来判断
//	// 这里简单地假设超过半数节点确认即可
//	count := 0
//	for _, confirmed := range pbft.prePrepareConfirmCount[messageID] {
//		if confirmed {
//			count++
//		}
//	}
//	return count > len(pbft.prePrepareConfirmCount[messageID])/2
//}
//
//// checkCommitConfirm 检查是否满足提交消息的确认条件
//func (pbft *PBFT) checkCommitConfirm(messageID string) bool {
//	// 在实际应用中，可以根据具体的确认条件来判断
//	// 这里简单地假设超过半数节点确认即可
//	count := 0
//	for _, confirmed := range pbft.commitConfirmCount[messageID] {
//		if confirmed {
//			count++
//		}
//	}
//	return count > len(pbft.commitConfirmCount[messageID])/2
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
//func (p *PBFT) commit() {
//	for {
//		select {
//		case block := <-p.commitChan:
//			fmt.Println("Received commit for block", block.BlockNumber)
//			// Implement commit logic
//			// ...
//
//			// Add the block to the committed blocks map
//			p.mutex.Lock()
//			p.committedBlocks[block.BlockNumber] = block
//			p.mutex.Unlock()
//		}
//	}
//}

//broadcastPrepare 模拟广播准备消息给其他节点
//func (p *PBFT) broadcastPrepare(block *Block) {
//	for id, _ := range NodeTable {
//		if p.node.ID == id{
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
//func (p *PBFT) QueryBlock(blockNumber int) *Block {
//	// 此处省略具体逻辑，可根据具体需求实现
//	// 节点查询指定区块信息
//	// 可以从已经提交的区块中查找
//	p.mutex.Lock()
//	defer p.mutex.Unlock()
//	return p.committedBlocks[blockNumber]
//}
