package test

import (
	"blockchain/crypto11"
	"blockchain/data"
	"fmt"
	"testing"
	"time"
)

func TestNewTransaction(t *testing.T) {
	// 生成密钥对
	privateKey, publicKey, err := crypto11.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Error generating key pair: %v", err)
	}

	// 创建一个新交易
	sender := "0x145789324"
	receiver := "0x145552314"
	amount := 10.0
	premium := 100.00

	transaction1, err := data.NewTransaction(sender, receiver, amount, premium, privateKey)

	fmt.Println(transaction1.ID)
	fmt.Println(transaction1.Amount)
	fmt.Println(transaction1.Premium)
	fmt.Println(transaction1.SenderAddress)
	fmt.Println(transaction1.ReceiverAddress)
	fmt.Println(transaction1.Timestamp)

	if err != nil {
		t.Fatalf("Error creating new transaction: %v", err)
	}

	// 验证交易的基本信息
	if transaction1.SenderAddress != sender {
		t.Errorf("Sender address mismatch. Expected %s, got %s", sender, transaction1.SenderAddress)
	}

	if transaction1.ReceiverAddress != receiver {
		t.Errorf("Receiver address mismatch. Expected %s, got %s", receiver, transaction1.ReceiverAddress)
	}

	if transaction1.Amount != amount {
		t.Errorf("Amount mismatch. Expected %f, got %f", amount, transaction1.Amount)
	}

	// 验证交易的签名
	valid := transaction1.VerifySignature(publicKey)
	if !valid {
		t.Error("Signature verification failed")
	}
}

func TestToJSONAndFromJSON(t *testing.T) {
	// 创建一个交易
	transaction := &data.Transaction{
		ID:              "testID",
		SenderAddress:   "testSender",
		ReceiverAddress: "testReceiver",
		Amount:          5.0,
		Timestamp:       time.Now(),
		Signature:       "testSignature",
		Premium:         55.0,
	}
	fmt.Println(transaction.Timestamp)

	// 将交易转换为JSON字符串
	jsonData, err := transaction.ToJSON()
	if err != nil {
		t.Fatalf("Error converting transaction to JSON: %v", err)
	}

	// 从JSON字符串恢复交易
	transaction1 := &data.Transaction{}
	err = transaction1.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Error restoring transaction from JSON: %v", err)
	}
	// 验证恢复的交易是否与原始交易相同
	fmt.Println(transaction1.ID)
	fmt.Println(transaction1.Amount)
	fmt.Println(transaction1.Premium)
	fmt.Println(transaction1.SenderAddress)
	fmt.Println(transaction1.ReceiverAddress)
	fmt.Println(transaction1.Timestamp)
	fmt.Println(transaction1.Signature)
	if transaction1.Amount != transaction.Amount {
		t.Error("Restored transaction does not match original transaction")
	}
}
