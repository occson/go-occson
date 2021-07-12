package occson

import (
	"io"
	"crypto/md5"
	"bytes"
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
)

func kdf(password, salt []byte) ([]byte, []byte) {
	hasher := md5.New()
	derivedKey := []byte{}
	block := []byte{}

	for len(derivedKey) < 48 {
		if len(block) != 0 {
			io.Copy(hasher, bytes.NewBuffer(block))
		}
		io.Copy(hasher, bytes.NewBuffer(password))
		io.Copy(hasher, bytes.NewBuffer(salt))
		block = hasher.Sum(nil)
		hasher.Reset()

		derivedKey = append(derivedKey, block...)
	}
	return derivedKey[0:32], derivedKey[32:]
}

func pkcs7Unpad(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func pkcs7Pad(ciphertext []byte, blockSize int) []byte {
   padding := blockSize - len(ciphertext) % blockSize
   padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   return append(ciphertext, padtext...)
}

func ccsDecrypt(encryptedContent, passphrase string) []byte {
	data, _ := base64.StdEncoding.Strict().DecodeString(encryptedContent)
	salt := data[8:16]
	ciphertext := data[16:]

	key, iv := kdf([]byte(passphrase), salt)

	block, _ := aes.NewCipher(key)
	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	ecb.CryptBlocks(decrypted, ciphertext)

	return pkcs7Unpad(decrypted)
}

func ccsEncrypt(plaintext, passphrase, salt string) string {
	key, iv := kdf([]byte(passphrase), []byte(salt))

	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
    padded := pkcs7Pad([]byte(plaintext), blockSize)

	ecb := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padded))
	ecb.CryptBlocks(ciphertext, padded)

	output := append([]byte("Salted__"), []byte(salt)...)
	output = append(output, ciphertext...)

	return base64.StdEncoding.Strict().EncodeToString(output)
}
