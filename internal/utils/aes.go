package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// Mã hóa truyền []byte trả về string
func EncryptAES(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// number used once
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// Seal (mã hóa): cipherText = nonce + cipherText + tag
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}
