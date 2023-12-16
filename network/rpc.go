// network/rpc.go

package network

import (
	"blockchain/data"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// RPCService 节点提供的 RPC 服务
type RPCService struct {
	Node *data.Node
}

// ReceiveTransaction 节点接收交易的 RPC 方法
func (r *RPCService) ReceiveTransaction(transaction data.Transaction) error {
	fmt.Printf("[%s] Received transaction: %s\n", r.Node.ID, transaction)

	// 处理交易，将其添加到节点的交易池中

	r.Node.AddTransaction(transaction)

	return nil
}

// SubmitBlockForConsensus 节点提交区块到共识过程的 RPC 方法
func (r *RPCService) SubmitBlockForConsensus(_ *struct{}, _ *struct{}) error {
	fmt.Printf("[%s] Submitted block for consensus\n", r.Node.ID)

	// 提交区块到共识过程
	r.Node.SubmitBlockForConsensus()

	return nil
}

// GetTransactionInfo 节点获取特定交易信息的 RPC 方法
func (r *RPCService) GetTransactionInfo(transaction string, reply *string) error {
	fmt.Printf("[%s] Retrieving transaction info for: %s\n", r.Node.ID, transaction)

	// 获取特定交易的信息
	info := r.Node.GetTransactionInfo(transaction)
	*reply = info

	return nil
}

// GetBlockInfo 节点获取特定区块信息的 RPC 方法
func (r *RPCService) GetBlockInfo(blockNumber int, reply *string) error {
	fmt.Printf("[%s] Retrieving block info for block %d\n", r.Node.ID, blockNumber)

	// 获取特定区块的信息
	info := r.Node.GetBlockInfo(blockNumber)
	*reply = info

	return nil
}

// StartRPCServer 启动节点的 RPC 服务器
func (n *RPCService) StartRPCServer() {
	rpc.Register(&RPCService{Node: n.Node})
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", n.Node.Addr)
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Printf("[%s] RPC server started at %s\n", n.Node.ID, n.Node.Addr)
	err = http.Serve(l, nil)
	if err != nil {
		fmt.Println("Error serving RPC server:", err)
		return
	}
}
