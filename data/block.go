// data/block.go

package data

// Block 区块结构
type Block struct {
	BlockNumber  int
	Transactions []Transaction
	Proposer     string // 在实际系统中，这可能是节点的公钥
}
