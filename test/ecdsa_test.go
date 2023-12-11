package test // crypto/ecdsa_test.go

import (
	"blocktest/crypto"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %v", err)
	}
	print(privKey, "\n", pubKey)
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
	//print(message, "\n")
	// 签名
	signature, err := crypto.Sign(privKey, message)
	if err != nil {
		t.Fatalf("Error signing message: %v", err)
	}
	print(signature)
	// 确认签名是否有效
	valid, err := crypto.Verify(pubKey, message, signature)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}

	if !valid {
		t.Error("Signature verification failed")
	}
}
