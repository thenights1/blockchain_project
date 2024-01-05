// // data/client.go
package data

import (
	"blockchain/crypto11"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Client 客户端的数据结构
type Client struct {
	ID            string
	Address       string
	Balance       float64
	TransactionID int
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
}

// NewClient 创建一个新的客户端实例
func NewClient(id, address string, balance float64, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *Client {
	return &Client{
		ID:            id,
		Address:       address,
		Balance:       balance,
		TransactionID: 0,
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
	}
}

// StartClient 模拟客户端交易行为
func (c *Client) StartClient() {
	fmt.Printf("Client %s started\n", c.ID)

	for c.Balance > 0 {
		// 模拟用户交易
		transaction := c.createTransaction()
		transactionJSON, err := transaction.ToJSON()

		if err != nil {
			fmt.Println("Error creating transaction JSON:", err)
			return
		}
		for id, addr := range NodeTable { //给每个结点都发信息，但只有主节点会进行处理
			Sendmessage([]byte("tran"+transactionJSON), addr)
			fmt.Printf("客户端%s发送了信息给结点%s\n", c.ID, id)
		}
		//Sendmessage([]byte("tran"+transactionJSON), NodeTable["Node1"])
		// 向节点提交交易
		//node.AddTransaction(transactionJSON)

		// 模拟用户等待
		time.Sleep(time.Second)
	}

	fmt.Printf("Client %s finished\n", c.ID)
}

// createTransaction 创建一个新的交易
func (c *Client) createTransaction() *Transaction {
	amount := rand.Float64()*20 + 5   // 5-25 的随机数
	fee := rand.Float64()*1.19 + 0.01 // 0.01-1.2 的随机数
	c.Balance -= (amount + fee)

	transaction, err := NewTransaction(c.Address, "otherReceiver", amount, fee, c.PrivateKey)
	if err != nil {
		fmt.Println("Error creating new transaction in client.go: %v", err)
	}

	c.TransactionID++
	return transaction
}

// RunClients 启动模拟客户端
func RunClients() {
	// 创建模拟用户
	//下面生成四对密钥给用户
	//分界线---------------------------
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
	//分界线-----------------------------
	Users = []*Client{
		NewClient("User1", "0x145287", 100.0, pri1, pub1),
		NewClient("User2", "0x124563", 110.0, pri2, pub2),
		NewClient("User3", "0x145235", 130.0, pri3, pub3),
		NewClient("User4", "0x147889", 120.0, pri4, pub4),
	}
	for _, user := range Users { //利用文件存储持久化密钥
		err := SaveKeysToFile(user)
		if err != nil {
			fmt.Println("Error saving keys for user %s: %v", user.ID, err)
		}
	}
	//err = SaveUsersToFile(users, "../user.json")
	//if err != nil {
	//	fmt.Println("error in creat file")
	//}
	//fmt.Println("145")
	// 启动模拟客户端
	for _, user := range Users {
		go user.StartClient()
	}
}
func SaveKeysToFile(client *Client) error {
	// 创建存储目录
	err := os.MkdirAll("./keys", os.ModePerm)
	if err != nil {
		return err
	}

	// 将私钥序列化为字符串
	privateKeyStr, err := serializeECDSAPrivateKey(client.PrivateKey)
	if err != nil {
		return err
	}

	// 将公钥序列化为字符串
	publicKeyStr, err := serializeECDSAPublicKey(client.PublicKey)
	if err != nil {
		return err
	}

	// 写入私钥到文件
	privateKeyFilePath := filepath.Join("./keys", fmt.Sprintf("%s_private_key.txt", client.Address))
	err = os.WriteFile(privateKeyFilePath, []byte(privateKeyStr), 0644)
	if err != nil {
		return err
	}

	// 写入公钥到文件
	publicKeyFilePath := filepath.Join("./keys", fmt.Sprintf("%s_public_key.txt", client.Address))
	err = os.WriteFile(publicKeyFilePath, []byte(publicKeyStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadKeysFromFile 从文件中读取字符串并转换为客户端的公私钥
func LoadKeysFromFile(client *Client) error {
	// 读取私钥文件
	privateKeyFilePath := filepath.Join("./keys", fmt.Sprintf("%s_private_key.txt", client.Address))
	privateKeyStr, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return err
	}

	// 读取公钥文件
	publicKeyFilePath := filepath.Join("./keys", fmt.Sprintf("%s_public_key.txt", client.Address))
	publicKeyStr, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		return err
	}

	// 反序列化私钥和公钥
	deserializedPrivateKey, err := deserializeECDSAPrivateKey(string(privateKeyStr))
	if err != nil {
		return err
	}

	deserializedPublicKey, err := deserializeECDSAPublicKey(string(publicKeyStr))
	if err != nil {
		return err
	}

	// 将反序列化的私钥和公钥设置到客户端
	client.PrivateKey = deserializedPrivateKey
	client.PublicKey = deserializedPublicKey

	return nil
}

// 将 ECDSA 私钥序列化为字符串
func serializeECDSAPrivateKey(key *ecdsa.PrivateKey) (string, error) {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return "", err
	}
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: bytes,
	}
	return string(pem.EncodeToMemory(block)), nil
}

// 将 ECDSA 公钥序列化为字符串
func serializeECDSAPublicKey(key *ecdsa.PublicKey) (string, error) {
	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}
	return string(pem.EncodeToMemory(block)), nil
}

// 从字符串反序列化 ECDSA 私钥
func deserializeECDSAPrivateKey(str string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 从字符串反序列化 ECDSA 公钥
func deserializeECDSAPublicKey(str string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch key := key.(type) {
	case *ecdsa.PublicKey:
		return key, nil
	default:
		return nil, fmt.Errorf("unexpected key type %T", key)
	}
}

//func SaveUsersToFile(users []*Client, filePath string) error {
//	jsonData, err := json.Marshal(users)
//	if err != nil {
//		return err
//	}
//
//	err = os.WriteFile(filePath, jsonData, 0644)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// // LoadUsersFromFile 从文件中加载用户数据
//
//	func LoadUsersFromFile(filePath string) ([]*Client, error) {
//		fileData, err := os.ReadFile(filePath)
//		if err != nil {
//			return nil, err
//		}
//
//		var users []*Client
//		err = json.Unmarshal(fileData, &users)
//		if err != nil {
//			fmt.Println("Error unmarshalling JSON:", err)
//			return nil, err
//		}
//
//		return users, nil
//	}
//func LoadUsersFromFile(filePath string) ([]*Client, error) {
//	fileData, err := os.ReadFile(filePath)
//	if err != nil {
//		return nil, err
//	}
//
//	// 定义结构体用于解析JSON数据
//	var jsonData []map[string]interface{}
//	err = json.Unmarshal(fileData, &jsonData)
//	if err != nil {
//		fmt.Println("Error unmarshalling JSON:", err)
//		return nil, err
//	}
//
//	// 创建用户列表
//	var users []*Client
//
//	// 遍历JSON数据，提取并创建Client实例
//	for _, userJSON := range jsonData {
//		id := userJSON["ID"].(string)
//		address := userJSON["Address"].(string)
//		balance, _ := strconv.ParseFloat(fmt.Sprint(userJSON["Balance"]), 64)
//
//		// 解析PrivateKey
//		var d, x, y *big.Int
//
//		if dJSON, ok := userJSON["D"]; ok && dJSON != nil {
//			if d, ok = dJSON.(*big.Int); !ok {
//				return nil, fmt.Errorf("Invalid D field in JSON")
//			}
//		}
//
//		if xJSON, ok := userJSON["X"]; ok && xJSON != nil {
//			if x, ok = xJSON.(*big.Int); !ok {
//				return nil, fmt.Errorf("Invalid X field in JSON")
//			}
//		}
//
//		if yJSON, ok := userJSON["Y"]; ok && yJSON != nil {
//			if y, ok = yJSON.(*big.Int); !ok {
//				return nil, fmt.Errorf("Invalid Y field in JSON")
//			}
//		}
//
//		privateKey := &ecdsa.PrivateKey{
//			D: d,
//			PublicKey: ecdsa.PublicKey{
//				X: x,
//				Y: y,
//			},
//		}
//
//		publicKey := &ecdsa.PublicKey{
//			X: x,
//			Y: y,
//		}
//
//		// 创建Client实例
//		client := NewClient(id, address, balance, privateKey, publicKey)
//		users = append(users, client)
//	}
//
//	return users, nil
//}
