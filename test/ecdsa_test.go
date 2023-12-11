package test // crypto/ecdsa_test.go

import (
	"blockchain/crypto"
	"fmt"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %v", err)
	}
	fmt.Println("Private Key:", privKey)
	fmt.Println("Public Key:", pubKey)

	if privKey == nil || pubKey == nil {
		t.Error("Generated key pair is nil")
	}
}

func TestSignAndVerify(t *testing.T) {
	privKey, pubKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Error generating key pair: %v", err)
	}

	message := []byte("test message, llllll")

	//fmt.Println(message)

	// 签名
	signature, err := crypto.Sign(privKey, message)
	if err != nil {
		t.Fatalf("Error signing message: %v", err)
	}

	fmt.Println("Sign pass!")

	// 确认签名是否有效
	valid, err := crypto.Verify(pubKey, message, signature)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}
	if !valid {
		t.Error("Signature verification failed")
	}

	//下面进行测试签名无效的情况
	privKey1, _, err1 := crypto.GenerateKeyPair()

	message1 := []byte("This is another message")

	signature1, err1 := crypto.Sign(privKey1, message1)
	if err1 != nil {
		t.Fatalf("Error signing message: %v", err)
	}
	//用第一个用户的公钥来验证是否正确，显然输出应该为不正确则测试成功
	valid1, err1 := crypto.Verify(pubKey, message, signature1)
	if err1 != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}
	if valid1 {
		t.Error("Signature error verification failed")
	}
}
