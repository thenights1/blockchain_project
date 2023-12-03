// crypto/ecdsa.go

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// GenerateKeyPair 生成 ECDSA 密钥对
func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return priv, &priv.PublicKey, nil
}

// Sign 使用私钥对消息进行签名
func Sign(privateKey *ecdsa.PrivateKey, message []byte) (string, error) {
	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	signature := r.String() + "," + s.String()
	return base64.StdEncoding.EncodeToString([]byte(signature)), nil
}

// Verify 验证签名是否有效
func Verify(publicKey *ecdsa.PublicKey, message []byte, signature string) (bool, error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(message)
	r := new(big.Int)
	s := new(big.Int)
	signatureParts := string(signatureBytes)
	fmt.Sscanf(signatureParts, "%s,%s", r, s)

	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}
