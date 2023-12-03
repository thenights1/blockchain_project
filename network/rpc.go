// network/rpc.go

package network

import (
	. "blockchain/data"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// RPCService 节点提供的 RPC 服务
type RPCService struct {
	Node *Node
}

// ReceivePrepare 节点接收准备消息的 RPC 方法
func (r *RPCService) ReceivePrepare(block *Block, _ *struct{}) error {
	fmt.Printf("[%s] Received prepare for block %d from %s\n", r.Node.Name, block.BlockNumber, block.Proposer)
	// 实际 PBFT 中的准备消息处理逻辑
	return nil
}

// ReceiveCommit 节点接收提交消息的 RPC 方法
func (r *RPCService) ReceiveCommit(block *Block, _ *struct{}) error {
	fmt.Printf("[%s] Received commit for block %d from %s\n", r.Node.Name, block.BlockNumber, block.Proposer)
	// 实际 PBFT 中的提交消息处理逻辑
	return nil
}

// StartRPCServer 启动节点的 RPC 服务器
func (n *Node) StartRPCServer() {
	rpc.Register(&RPCService{Node: n})
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", n.Addr)
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Printf("[%s] RPC server started at %s\n", n.Name, n.Addr)
	err = http.Serve(l, nil)
	if err != nil {
		fmt.Println("Error serving RPC server:", err)
		return
	}
}