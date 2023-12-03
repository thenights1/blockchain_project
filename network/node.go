// network/node.go

package network

import (
	. "blockchain/data"
	"fmt"
	"net/rpc"
	"sync"
)

// Node 节点的数据结构
type Node struct {
	Name       string
	Addr       string
	rpcClient  *rpc.Client
	PublicKey  string // 在实际系统中，这可能是节点的公钥
	PrivateKey string // 在实际系统中，这可能是节点的私钥
	mutex      sync.Mutex
}

// NewNode 创建一个新的节点
func NewNode(name, addr, publicKey, privateKey string) *Node {
	return &Node{
		Name:       name,
		Addr:       addr,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// StartRPCClient 启动节点的 RPC 客户端
func (n *Node) StartRPCClient() {
	client, err := rpc.DialHTTP("tcp", n.Addr)
	if err != nil {
		fmt.Println("Error starting RPC client:", err)
		return
	}

	n.mutex.Lock()
	n.rpcClient = client
	n.mutex.Unlock()
	fmt.Printf("[%s] RPC client started at %s\n", n.Name, n.Addr)
}

// RPCClient 返回节点的 RPC 客户端
func (n *Node) RPCClient() *rpc.Client {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.rpcClient
}

// SendPrepare 向指定节点发送准备消息
func (n *Node) SendPrepare(target *Node, block *Block) {
	client := target.RPCClient()
	if client != nil {
		go func() {
			var reply struct{}
			err := client.Call("RPCService.ReceivePrepare", block, &reply)
			if err != nil {
				fmt.Printf("[%s] Error sending prepare to %s: %s\n", n.Name, target.Name, err)
			}
		}()
	}
}

// SendCommit 向指定节点发送提交消息
func (n *Node) SendCommit(target *Node, block *Block) {
	client := target.RPCClient()
	if client != nil {
		go func() {
			var reply struct{}
			err := client.Call("RPCService.ReceiveCommit", block, &reply)
			if err != nil {
				fmt.Printf("[%s] Error sending commit to %s: %s\n", n.Name, target.Name, err)
			}
		}()
	}
}

// SomeFunction 在 node.go 中定义的一个函数
func (n *Node) SomeFunction() {
	fmt.Printf("[%s] This is a function from node.go\n", n.Name)
}

// Start 模拟节点启动
func (n *Node) Start() {
	fmt.Printf("[%s] Node started\n", n.Name)
}

// Stop 模拟节点停止
func (n *Node) Stop() {
	fmt.Printf("[%s] Node stopped\n", n.Name)
}

// AddTransaction 模拟添加交
