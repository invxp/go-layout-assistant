package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
)

/*
工具包
常规Crypto实现
主要包括AES-PKCS7加解密和MD5校验
*/

//Base64Decode 解Base64字符串(尝试4种方式)
func Base64Decode(plantText string) ([]byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(plantText)
	if err == nil {
		return data, nil
	}
	data, err = base64.RawStdEncoding.DecodeString(plantText)
	if err == nil {
		return data, err
	}
	data, err = base64.URLEncoding.DecodeString(plantText)
	if err == nil {
		return data, err
	}
	return base64.StdEncoding.DecodeString(plantText)
}

//AESPKCS7Encode AES-PKCS7加密
func AESPKCS7Encode(originBytes, key, iv []byte) (string, error) {
	if len(originBytes) == 0 {
		return "", fmt.Errorf("empty string")
	}
	keyLength := 16 - len(key)
	for i := 0; i < keyLength; i++ {
		key = append(key, 0)
	}

	ivLength := 16 - len(iv)
	for i := 0; i < ivLength; i++ {
		iv = append(iv, 0)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	tmp := pkcs7Padding(originBytes, block.BlockSize())
	cypherText := make([]byte, len(tmp))
	blockModel := cipher.NewCBCEncrypter(block, iv)
	blockModel.CryptBlocks(cypherText, tmp)
	return base64.StdEncoding.EncodeToString(cypherText), nil
}

//AESPKCS7Decode AES-PKCS7解密
func AESPKCS7Decode(encodedString string, key, iv []byte) ([]byte, error) {
	if encodedString == "" {
		return nil, fmt.Errorf("empty string")
	}

	lenKey := 16 - len(key)
	for i := 0; i < lenKey; i++ {
		key = append(key, 0)
	}

	lenIv := 16 - len(iv)
	for i := 0; i < lenIv; i++ {
		iv = append(iv, 0)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	b64Str, err := Base64Decode(encodedString)
	if err != nil {
		return nil, err
	}

	cypherText := b64Str
	blockModel := cipher.NewCBCDecrypter(block, iv)
	finalStr := make([]byte, len(cypherText))
	blockModel.CryptBlocks(finalStr, cypherText)
	finalStr = pkcs7UnPadding(finalStr)
	return finalStr, nil
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}

//FileMD5 计算文件MD5
func FileMD5(fullPath string, upper bool) (string, error) {
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return MD5(data, upper), nil
}

//MD5 计算数据MD5
func MD5(input []byte, upper bool) string {
	md5Ctx := md5.New()
	md5Ctx.Write(input)
	if upper {
		return strings.ToUpper(hex.EncodeToString(md5Ctx.Sum(nil)))
	}
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
