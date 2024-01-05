package test

import (
	"blockchain/data"
	"testing"
)

func TestSaveBlockchainToJSON(t *testing.T) {
	// 创建一个区块链实例
	blockchain := data.NewBlockchain(&data.Node{})

	// 添加一些区块
	block1 := data.NewBlock(1, nil, "data1", &data.Node{})
	block2 := data.NewBlock(2, nil, "data2", &data.Node{})
	blockchain.AddBlock(block1, &data.Node{})
	blockchain.AddBlock(block2, &data.Node{})

	// 保存区块链到JSON文件
	filePath := "test_blockchain.json"
	if err := blockchain.SaveBlockchainToJSON(filePath); err != nil {
		t.Fatalf("Error saving blockchain to JSON: %v", err)
	}

}
