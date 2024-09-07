package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	prvkey []byte
	pubkey []byte
	once   sync.Once
)

// RsaEncrypt RSA加密
func RsaEncrypt(key string) (string, error) {
	block, _ := pem.Decode(pubkey)
	if block == nil {
		slog.Error("decode private key is nil")
		return "", errors.New("decode public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		slog.Error("parse public key error", "err", err)
		return "", err
	}
	encryKey, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(key))
	if err != nil {
		slog.Error("encrypt key error", "err", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryKey), nil
}

// RsaDecrypt RSA解密
func RsaDecrypt(key string) (string, error) {
	block, _ := pem.Decode(prvkey)
	if block == nil {
		slog.Error("decode private key is nil")
		return "", errors.New("decode private key error")
	}
	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("parse private key error", "err", err)
		return "", err
	}
	decodeKey, _ := base64.StdEncoding.DecodeString(key)
	decryKey, err := rsa.DecryptPKCS1v15(rand.Reader, prv, decodeKey)
	if err != nil {
		slog.Error("decrypt key error", "err", err)
		return "", err
	}
	return string(decryKey), nil
}

// GenRsaKey 生成RSA密钥对
// 正常来说，应该有密钥对管理机制，但为了简单，在项目生成密钥对，用uuid作为密钥对文件的绑定
func GenRsaKey(bits int) (string, error) {
	rsaAk := strings.ReplaceAll(uuid.New().String(), "-", "")
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		slog.Error("generates rsa private key fail", "err", err)
		return "", err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	prvBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateFile, err := os.Create(fmt.Sprintf("./secret_key/%s/private_key.pem", rsaAk))
	if err != nil {
		slog.Error("create rsa private key file fail", "err", err)
		return "", err
	}
	defer privateFile.Close()
	err = pem.Encode(privateFile, prvBlock)
	if err != nil {
		slog.Error("pem encode private key file fail", "err", err)
		return "", err
	}

	// 生成公钥
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		slog.Error("converts a public key to PKIX fail", "err", err)
		return "", err
	}
	pubBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	publicFile, err := os.Create(fmt.Sprintf("./secret_key/%s/public_key.pem", rsaAk))
	if err != nil {
		slog.Error("create rsa public file fail", "err", err)
		return "", err
	}
	defer publicFile.Close()
	err = pem.Encode(publicFile, pubBlock)
	if err != nil {
		slog.Error("pem encode public key file fail", "err", err)
		return "", err
	}
	return rsaAk, nil
}

func InitRsaKey(ak string) {
	once.Do(func() {
		privateData, err := os.ReadFile(fmt.Sprintf("./secret_key/%s/private_key.pem", ak))
		if err != nil {
			panic(fmt.Sprintf("read private key file fail, err = %v", err))
		}
		publicData, err := os.ReadFile(fmt.Sprintf("./secret_key/%s/public_key.pem", ak))
		if err != nil {
			panic(fmt.Sprintf("read public key file fail, err = %v", err))
		}
		prvkey = privateData
		pubkey = publicData
	})
}
