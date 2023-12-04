// data/transaction.go

package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Transaction 交易的数据结构
type Transaction struct {
	ID        string
	Sender    string
	Recipient string
	Amount    int
}

// NewTransaction 创建一个新的交易实例
func NewTransaction(sender, recipient string, amount int) *Transaction {
	transaction := &Transaction{
		ID:        generateTransactionID(sender, recipient, amount),
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	return transaction
}

// generateTransactionID 生成交易的唯一标识符
func generateTransactionID(sender, recipient string, amount int) string {
	data := fmt.Sprintf("%s%s%d", sender, recipient, amount)
	hashInBytes := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashInBytes[:])
}
