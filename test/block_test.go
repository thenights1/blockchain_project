package test

import (
	"blockchain/data"
	"testing"
	"time"
)

func TestNewBlock(t *testing.T) {
	// 创建一些交易
	transactions := []*data.Transaction{
		{
			ID:              "testID1",
			SenderAddress:   "testSender1",
			ReceiverAddress: "testReceiver1",
			Amount:          5.0,
			Timestamp:       time.Now(),
			Signature:       "testSignature1",
			Premium:         4.0,
		},
		{
			ID:              "testID2",
			SenderAddress:   "testSender2",
			ReceiverAddress: "testReceiver2",
			Amount:          8.0,
			Timestamp:       time.Now(),
			Signature:       "testSignature2",
			Premium:         5.0,
		},
	}

	// 创建一个提出者结点
	proposer := &data.Node{ID: "testProposer"}

	// 创建一个新区块
	blockNumber := 1
	prevHash := "testPrevHash"
	block := data.NewBlock(blockNumber, transactions, prevHash, proposer)

	// 验证区块的基本信息
	if block.BlockNumber != blockNumber {
		t.Errorf("Block number mismatch. Expected %d, got %d", blockNumber, block.BlockNumber)
	}

	// 验证区块的哈希值
	expectedHash := block.CalculateHash()
	if block.Hash != expectedHash {
		t.Errorf("Block hash mismatch. Expected %s, got %s", expectedHash, block.Hash)
	}
}
func TestAddTrac(t *testing.T) { //测试添加交易到交易池
	//nodeaddr := "localhost:8888"
	//id := "NO.1"
	////node := data.NewNode(nodeaddr, id)
	//
	//sender := "0x145789324"
	//receiver := "0x145552314"
	//amount := 10.0
	//premium := 100.00
	//
	//transaction1, err := data.NewTransaction(sender, receiver, amount, premium, privateKey)
	//
	//if err != nil {
	//	t.Fatalf("Error creating new transaction: %v", err)
	//}
}
