// data/blockchain.go

package data

import "sync"

// Blockchain 区块链结构
type Blockchain struct {
	Chain []Block
	Lock  sync.Mutex
}
