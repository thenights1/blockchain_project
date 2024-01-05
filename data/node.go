// network/node.go

package data

import (
	"blockchain/crypto11"
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Node 节点的数据结构
type Node struct {
	ID              string            //node的唯一标识符
	TransactionPool []*Transaction    //存储待处理交易的池。
	Blockchain      *Blockchain       //表示节点拥有的区块链
	PublicKey       *ecdsa.PublicKey  //存储公钥
	PrivateKey      *ecdsa.PrivateKey //存储私钥
	Consensus       *PBFT             //共识算法，这里是pbft
	tempBlocks      []*Block          //临时区块池，存储打包好还未提交到共识算法的区块
	tempblock       *Block            //临时块，便于快速加入区块链
	mutex           sync.Mutex        //用于在多协程中保护数据一致性的互斥锁
	Synchronized    bool              //用于判断节点是否与其他节点同步
	stopChan        chan struct{}     //用于通知节点停止的通道
	Addr            string
	View            int
}

// NewNode 创建一个新的节点实例
func NewNode(addr string, id string, bc *Blockchain) *Node {

	privateKey, publicKey, err := crypto11.GenerateKeyPair()
	if err != nil {
		fmt.Printf("Error generating key pair: %v", err)
	}
	node := &Node{
		ID:              id,
		TransactionPool: make([]*Transaction, 0),
		Blockchain:      bc,
		stopChan:        make(chan struct{}),
		Synchronized:    false,
		Addr:            addr,
		PrivateKey:      privateKey,
		PublicKey:       publicKey,
	}
	return node
}

// Start 启动节点
func (n *Node) Start() {
	fmt.Printf("Node %s started\n", n.ID)
	flag := 0
	for {
		select {
		case <-n.stopChan:
			fmt.Printf("Node %s stopped\n", n.ID)
			return
		default:
			if flag != 1 {
				go n.tcpListen()
				flag = 1
			}
			// 模拟节点的周期性活动
			time.Sleep(time.Second)
			var newblock *Block

			n.TransactionPool, newblock = n.PackTransactions(n.TransactionPool)
			if newblock != nil {
				n.tempBlocks = append(n.tempBlocks, newblock) //加入临时区块池，之后提交
			}
			go func() {
				n.mutex.Lock()
				var block *Block
				if len(n.tempBlocks) > 0 {
					block = n.tempBlocks[0]
					n.tempBlocks = n.tempBlocks[1:] //切去第一个元素，因为即将被提交到pbft
					n.tempblock = block
					n.Consensus.PrePrepare(block)
				}
				n.mutex.Unlock()
			}()

		}
	}
	//n.tcpListen()
}

// Stop 停止节点
func (n *Node) Stop() {
	close(n.stopChan)
}

func (n *Node) HandleRequest(request string) {
	flag := request[:4]
	remainingString := request[4:]
	if flag == "tran" { //处理提交交易请求
		if n.View != 3 {
			return
		}
		transaction_temp := &Transaction{}
		err := transaction_temp.FromJSON(remainingString)
		if err != nil {
			fmt.Println("Error restoring transaction from JSON in node.go: %v", err)
		}
		loadedUser := NewClient("d", transaction_temp.SenderAddress, 10.0, nil, nil)
		err = LoadKeysFromFile(loadedUser)
		if err != nil {
			fmt.Println("Error loading keys in node.go", err)
		}
		f := transaction_temp.VerifySignature(loadedUser.PublicKey)
		if !f {
			fmt.Printf("收到交易ID：%s,且验证有效，即将加入待处理池\n", transaction_temp.ID)
			n.AddTransaction(transaction_temp)
		} else {
			fmt.Println("发现违法交易")
		}
	} else if flag == "preM" { //处理提交预准备请求
		nodeid := remainingString[:5]
		remainingString := remainingString[4:]
		fmt.Printf("收到来自主节点 %s 的预准备请求\n", nodeid)
		nodex := NewNode(NodeTable[nodeid], nodeid, n.Blockchain) //这里newnode增加到区块链可以忽略，主要是想提取到密钥
		nodex.PrivateKey = nil
		nodex.PublicKey = nil
		NodeLoadKeysFromFile(nodex)
		result := strings.SplitN(remainingString, " ", 2)
		//fmt.Println(nodex.PublicKey)
		f, _ := crypto11.Verify(nodex.PublicKey, []byte(result[1]), result[0])
		if !f {
			fmt.Println("预准备信息验证通过，将向其他结点广播准备信息")
			n.Consensus.broadcast("prep" + "Come on !")
		} else {
			fmt.Println("预准备信息验证失败！")
		}
		//println(remainingString)
	} else if flag == "prep" {
		if n.View == 3 { //只需要主节点做出反馈
			if remainingString == "Come on !" {
				n.Consensus.mutex.Lock()
				n.Consensus.prePrepareConfirmCount++
				if n.Consensus.prePrepareConfirmCount >= 3 {
					fmt.Println("准备过程通过，进入提交阶段")
					n.Consensus.broadcast("comi" + "Final !")
					n.Consensus.prePrepareConfirmCount = 1
				}
				n.Consensus.mutex.Unlock()
			}
		}
	} else if flag == "comi" {
		if remainingString == "Final !" {
			fmt.Println("决定通过提交，允许主节点继续操作")
			n.Consensus.broadcast("fina" + "GOGOGO")
		}
	} else if flag == "fina" {
		if n.View == 3 {
			if remainingString == "GOGOGO" {
				n.Consensus.mutex.Lock()
				n.Consensus.commitConfirmCount++
				if n.Consensus.commitConfirmCount >= 3 {
					fmt.Println("即将提交区块")
					n.Consensus.commitConfirmCount = 1
					n.Blockchain.AddBlock(n.tempblock, n)
					time.Sleep(time.Second)
					fmt.Println("成功增加到区块链")
					n.Blockchain.SaveBlockchainToJSON("./blockchain.json")

				}
				n.Consensus.mutex.Unlock()
			}
		}
	}

}

// AddTransaction 添加新交易到待处理池
func (n *Node) AddTransaction(transaction *Transaction) {
	// 使用互斥锁，确保在多线程环境下的安全并发操作
	n.mutex.Lock()
	// 函数结束后解锁
	defer n.mutex.Unlock()

	fmt.Printf("Node %s 添加了交易，ID为: %s\n", n.ID, transaction.ID)
	n.TransactionPool = append(n.TransactionPool, transaction)
}

// PackTransactions 将待处理池中的交易打包成区块
func (n *Node) PackTransactions(transactionPool []*Transaction) ([]*Transaction, *Block) {

	const maxTransactionsPerBlock = 4
	const maxFeeThreshold = 2.0
	n.mutex.Lock() //上锁
	defer n.mutex.Unlock()
	if len(transactionPool) == 0 {
		return transactionPool, nil
	}
	fmt.Println("开始打包区块")

	// 按手续费大小排序
	SortTransactionsByFee(transactionPool)

	var selectedTransactions []*Transaction
	totalFee := 0.0

	for _, transaction := range transactionPool {
		if len(selectedTransactions) >= maxTransactionsPerBlock {
			break
		}

		// 如果加入该交易后手续费总额不超过阈值，则加入选中的交易列表
		if totalFee+transaction.Premium <= maxFeeThreshold {
			selectedTransactions = append(selectedTransactions, transaction)
			totalFee += transaction.Premium
		}
	}

	// 从待处理池中移除已选中的交易
	remainingTransactions := RemoveSelectedTransactions(transactionPool, selectedTransactions)

	// 打包交易成区块
	blockNumber := 0 //暂时变为0，确定加入区块链后再计算值
	prehash := "xxx" //暂定值，确定加入区块链后再计算值
	newBlock := NewBlock(blockNumber, selectedTransactions, prehash, n)

	// 返回剩余的交易和新创建的区块
	return remainingTransactions, newBlock
}

// SortTransactionsByFee 按手续费大小对交易进行排序
func SortTransactionsByFee(transactions []*Transaction) {
	// 实现交易按手续费排序的逻辑，可以使用标准库的排序方法
	// 这里简化为按手续费降序排序
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Premium > transactions[j].Premium
	})
}

// RemoveSelectedTransactions 从交易池中移除已选中的交易
func RemoveSelectedTransactions(transactionPool, selectedTransactions []*Transaction) []*Transaction {

	var remainingTransactions []*Transaction

	for _, transaction := range transactionPool {
		if !contains(selectedTransactions, transaction) {
			remainingTransactions = append(remainingTransactions, transaction)
		}
	}

	return remainingTransactions
}

// contains 检查交易是否包含在已选中的交易列表中
func contains(transactions []*Transaction, target *Transaction) bool {
	for _, transaction := range transactions {
		if transaction == target {
			return true
		}
	}
	return false
}

// GetTransactionInfo 获取特定交易的信息
func (n *Node) GetTransactionInfo(transaction string) string {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return fmt.Sprintf("Transaction Info for %s from Node %s", transaction, n.ID)
}

// GetBlockInfo 获取特定区块的信息
func (n *Node) GetBlockInfo(blockNumber int) string {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return fmt.Sprintf("Block Info for Block %d from Node %s", blockNumber, n.ID)
}

// 存储和读取密钥，和客户端代码基本一致
func NodeSaveKeysToFile(n *Node) error {
	// 创建存储目录
	err := os.MkdirAll("./nodekeys", os.ModePerm)
	if err != nil {
		return err
	}

	// 将私钥序列化为字符串
	privateKeyStr, err := serializeECDSAPrivateKey(n.PrivateKey)
	if err != nil {
		return err
	}

	// 将公钥序列化为字符串
	publicKeyStr, err := serializeECDSAPublicKey(n.PublicKey)
	if err != nil {
		return err
	}

	// 写入私钥到文件
	privateKeyFilePath := filepath.Join("./nodekeys", fmt.Sprintf("%s_private_key.txt", n.ID))
	err = os.WriteFile(privateKeyFilePath, []byte(privateKeyStr), 0644)
	if err != nil {
		return err
	}

	// 写入公钥到文件
	publicKeyFilePath := filepath.Join("./nodekeys", fmt.Sprintf("%s_public_key.txt", n.ID))
	err = os.WriteFile(publicKeyFilePath, []byte(publicKeyStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadKeysFromFile 从文件中读取字符串并转换为结点的公私钥
func NodeLoadKeysFromFile(n *Node) error {
	// 读取私钥文件
	privateKeyFilePath := filepath.Join("./nodekeys", fmt.Sprintf("%s_private_key.txt", n.ID))
	privateKeyStr, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return err
	}

	// 读取公钥文件
	publicKeyFilePath := filepath.Join("./nodekeys", fmt.Sprintf("%s_public_key.txt", n.ID))
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

	// 将反序列化的私钥和公钥设置到结点
	n.PrivateKey = deserializedPrivateKey
	n.PublicKey = deserializedPublicKey

	return nil
}
