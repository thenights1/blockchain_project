// consensus/pbft.go

package data

import (
	"blockchain/crypto11"
	"fmt"
	"math/rand"
	"sync"
	"time"
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
		node:                   NewNode(n.Addr, n.ID, n.Blockchain),
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

// broadcast 向其他节点广播消息
func (pbft *PBFT) broadcast(message string) {
	fmt.Printf("Node %s broadcasts message\n", pbft.node.ID)
	time.Sleep(time.Second)
	for id, addr := range NodeTable {
		if id == pbft.node.ID {
			continue
		}
		go Sendmessage([]byte(message), addr)
	}
}
