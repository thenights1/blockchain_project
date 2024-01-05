package test

import (
	"blockchain/data"
	"fmt"
	"testing"
)

var bc *data.Blockchain

func TestCreate(t *testing.T) {
	nodeaddr := "8888"
	id := "NO.1"
	node := data.NewNode(nodeaddr, id, bc)
	fmt.Println(node)
	fmt.Println(node.ID)
	fmt.Println(node.Blockchain)
	//t.Fatalf("aaa")
}

func TestSendmessage(t *testing.T) {
	node1addr := "127.0.0.1:8898"
	id1 := "NO.1"
	node := data.NewNode(node1addr, id1, bc)
	go node.Start()
	node2addr := "127.0.0.1:8088"
	id2 := "NO.2"
	node2 := data.NewNode(node2addr, id2, bc)
	go node2.Start()
	message := []byte("test message, llllll")
	go data.Sendmessage(message, node2.Addr)

}
func TestHandlequest(t *testing.T) {
	request := "{\"ID\":\"J7P9QWqnbpAk8zvZSTutoDo4KgjkZqH0q3AeeCFDdc4=\",\"SenderAddress\":\"0x145235\",\"ReceiverAddress\":\"otherReceiver\",\"Amount\":23.176979628254415,\"Timestamp\":\"2024-01-05T00:46:45.7117438+08:00\",\"Signature\":\"TWpReU1qVXdPVEl4TmpZeU9UVXpOVEU1T0RRMU5EZzBOamczTXpZd09EY3dNVEF4TkRJeE1UY3dPVFExT0RjMU5UVXlOekkxT0RNNU56Y3lPREEyTlRFeE5UTTJNVFkwTWpFME1UWXNPRGN6T1RVNU1USXlNekEyTWpjMk1qRTNNekkzT1RrMU5qWXhPRGN5TkRrNE5qTTBOelEwTlRNME1EYzVOemMxTlRFMk5ERTRPVFU0TXpVek1USXpOREkyTmpBMU9UZ3dNVEF5TlRJPQ==\",\"Premium\":0.8142253782058623}"
	node1addr := "127.0.0.1:8898"
	id1 := "NO.1"
	node := data.NewNode(node1addr, id1, bc)
	node.HandleRequest(request)
	//fmt.Println(node.TransactionPool[0].Premium)
}
