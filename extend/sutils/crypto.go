package sutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func MD5(in string) string {

	alg := md5.New()
	alg.Write([]byte(in))

	return hex.EncodeToString(alg.Sum(nil))
}

func SHA1(in string) string {
	alg := sha1.New()
	alg.Write([]byte(in))

	return hex.EncodeToString(alg.Sum(nil))
}

func SHA256(in string) string {
	alg := sha256.New()
	alg.Write([]byte(in))
	return hex.EncodeToString(alg.Sum(nil))
}

func SHA512(in string) string {
	alg := sha512.New()
	alg.Write([]byte(in))
	return hex.EncodeToString(alg.Sum(nil))
}

func MD5ToBase64(in string) string {

	alg := md5.New()
	alg.Write([]byte(in))

	return Base64Encode2(alg.Sum(nil))
}

func HmacSha1ToBase64(data, key string) string {
	hmac := hmac.New(sha1.New, []byte(key))
	hmac.Write([]byte(data))

	return Base64Encode2(hmac.Sum(nil))
}

func HmacSha256ToHex(data, key string) string {
	hmac := hmac.New(sha256.New, []byte(key))
	hmac.Write([]byte(data))

	return hex.EncodeToString(hmac.Sum(nil))
}
