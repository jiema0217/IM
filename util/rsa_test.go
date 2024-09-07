package util

import (
	"fmt"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	// 能存放长度为53的密码，理论上是够用的
	ak, err := GenRsaKey(1 << 9)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("GenRsaKey OK, ak:", ak)
}
func TestRsa(t *testing.T) {
	InitRsaKey("3ef9a6948fb0413facaec9fbfac58809")
	key := "123456"
	encryptKey, err := RsaEncrypt(key)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("encrypt:", encryptKey)
	decryptKey, err := RsaDecrypt(encryptKey)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("decrypt:", decryptKey)
	fmt.Println(key == decryptKey)
}
