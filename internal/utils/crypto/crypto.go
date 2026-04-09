package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AESEncrypt 使用 AES-GCM 加密数据
func AESEncrypt(plaintext []byte, key string) (string, error) {
	if key == "" {
		return "", errors.New("encryption key is empty")
	}

	// 解码 base64 密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		// 如果解码失败，尝试直接使用原始字符串
		keyBytes = []byte(key)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt 解密数据
func AESDecrypt(ciphertextStr string, key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("decryption key is empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return nil, err
	}

	// 解码 base64 密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		// 如果解码失败，尝试直接使用原始字符串
		keyBytes = []byte(key)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptURL 加密 URL
func EncryptURL(url string, key string) (string, error) {
	return AESEncrypt([]byte(url), key)
}

// DecryptURL 解密 URL
func DecryptURL(encryptedURL string, key string) (string, error) {
	decrypted, err := AESDecrypt(encryptedURL, key)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
