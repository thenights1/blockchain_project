// data/transaction.go

package data

import (
	"blockchain/crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Transaction 交易的数据结构
type Transaction struct {
	ID              string    // 交易ID
	SenderAddress   string    // 发送者地址
	ReceiverAddress string    // 接收者地址
	Amount          float64   // 交易金额
	Timestamp       time.Time // 时间戳
	Signature       string    // 数字签名
	Premium         float64   //交易手续费
}

// NewTransaction 创建一个新的交易
func NewTransaction(sender string, receiver string, amount float64, premium float64, privateKey *ecdsa.PrivateKey) (*Transaction, error) {
	transaction := &Transaction{
		SenderAddress:   sender,
		ReceiverAddress: receiver,
		Amount:          amount,
		Premium:         premium,
		Timestamp:       time.Now(),
	}

	// 生成交易ID
	transaction.generateID()

	// 对交易进行签名
	err := transaction.sign(privateKey)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// 生成交易ID
func (t *Transaction) generateID() {
	data := fmt.Sprintf("%s%s%s%f%s", t.SenderAddress, t.ReceiverAddress, t.Timestamp.String(), t.Amount, t.Signature)
	hash := sha256.New()
	hash.Write([]byte(data))
	t.ID = base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// 对交易进行签名
func (t *Transaction) sign(privateKey *ecdsa.PrivateKey) error {
	data := fmt.Sprintf("%s%s%s%f%s", t.SenderAddress, t.ReceiverAddress, t.Timestamp.String(), t.Amount, t.ID)

	hash := sha256.New()
	hash.Write([]byte(data))

	signature, err := crypto.Sign(privateKey, hash.Sum(nil))
	signature1 := []byte(signature)
	if err != nil {
		return err
	}

	t.Signature = base64.URLEncoding.EncodeToString(signature1)
	return nil
}

// VerifySignature 验证交易的签名
func (t *Transaction) VerifySignature(publicKey *ecdsa.PublicKey) bool {
	data := fmt.Sprintf("%s%s%s%f%s", t.SenderAddress, t.ReceiverAddress, t.Timestamp.String(), t.Amount, t.ID)

	hash := sha256.New()
	hash.Write([]byte(data))

	signature, err := base64.URLEncoding.DecodeString(t.Signature)
	if err != nil {
		return false
	}
	sign := ""
	sign = string(signature)
	valid, err := crypto.Verify(publicKey, hash.Sum(nil), sign)
	return valid
}

// ToJSON 将交易转换为JSON格式的字符串
func (t *Transaction) ToJSON() (string, error) {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// FromJSON 从JSON格式的字符串恢复交易
func (t *Transaction) FromJSON(jsonData string) error {
	err := json.Unmarshal([]byte(jsonData), &t)
	if err != nil {
		return err
	}
	return nil
}
