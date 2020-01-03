package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// https://golang.org/src/crypto/cipher/example_test.go
func AESDecrypt(key, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	blockSize := block.BlockSize()

	if len(ciphertext) < blockSize {
		err = errors.New("ciphertext too short")
		return
	}

	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	if len(ciphertext)%blockSize != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	//  解填充
	plaintext = PKCS7UnPadding(ciphertext)

	return
}

// https://golang.org/src/crypto/cipher/example_test.go
// /aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AESEecrypt(key, plaintext []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)

	ciphertext = make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize] //初始向量IV必须是唯一，但不需要保密
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[blockSize:], plaintext)

	return
}
