package sutils

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// Des and Base64 encode
func DesEncrypt(src, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}

	return Base64Encode2(out), nil
}

func DesDecrypt(src, key []byte) (string, error) {
	src, err := Base64Decode2(string(src))
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}

	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out, err = PKCS5UnPadding(out)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// 3DES CBC base64
func TripleDesCBCEncrypt(origData, key []byte, iv []byte) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return Base64Encode2(crypted), nil
}

func TripleDesCBCDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS5UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

/*
	算法: DES
	模式: CBC
	填充: PKCS5(Block的大小是8位)
	偏移量: 默认使用key
	输出: 加密原数据
*/
func DesCBCEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS5Padding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key)

	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}
