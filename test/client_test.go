package test

import (
	"blockchain/crypto11"
	"blockchain/data"
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {
	//data.RunClients()
	pri1, pub1, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair in client.go")
	}
	pri2, pub2, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair in client.go")
	}
	pri3, pub3, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair in client.go")
	}
	pri4, pub4, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair in client.go")
	}
	//data.RunClients()
	data.Users = []*data.Client{
		data.NewClient("User1", "0x145287", 100.0, pri1, pub1),
		data.NewClient("User2", "0x124563", 150.0, pri2, pub2),
		data.NewClient("User3", "0x145235", 200.0, pri3, pub3),
		data.NewClient("User4", "0x147889", 120.0, pri4, pub4),
	}
	// 存储用户数组到文件
	for _, user := range data.Users {
		err := data.SaveKeysToFile(user)
		if err != nil {
			t.Errorf("Error saving keys for user %s: %v", user.ID, err)
		}
	}

	// 从文件中读取用户数组
	var loadedUsers []*data.Client
	for _, user := range data.Users {
		loadedUser := data.NewClient(user.ID, user.Address, user.Balance, nil, nil)
		err := data.LoadKeysFromFile(loadedUser)
		if err != nil {
			t.Errorf("Error loading keys for user %s: %v", loadedUser.ID, err)
		}
		loadedUsers = append(loadedUsers, loadedUser)
	}
	// 比较原始用户数组和从文件中加载的用户数组
	for i, user := range data.Users {
		loadedUser := loadedUsers[i]

		// 比较私钥
		if user.PrivateKey.D.Cmp(loadedUser.PrivateKey.D) != 0 {
			t.Errorf("Private keys do not match for user %s", user.ID)
		}

		// 比较公钥
		if user.PublicKey.X.Cmp(loadedUser.PublicKey.X) != 0 || user.PublicKey.Y.Cmp(loadedUser.PublicKey.Y) != 0 {
			t.Errorf("Public keys do not match for user %s", user.ID)
		}
	}

	fmt.Println("Test completed.")

}
