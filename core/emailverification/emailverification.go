package emailverification

import (
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
)

type EmailVerificationStrategy interface {
	GenerateCode() (rawCode string, protectedCode string, strategyName string)
}

func NewEmailVerificationStrategy() EmailVerificationStrategy {
	return &defaultStrategy{}
}

type defaultStrategy struct{}

func (*defaultStrategy) GenerateCode() (rawCode string, protectedCode string, strategyName string) {
	var (
		code       = make([]byte, 8)
		charset    = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		lenCharset = byte(len(charset))
	)

	rand.Read(code)

	for i := range code {
		code[i] = charset[code[i]%lenCharset]
	}
	rawCode = string(code)

	digest := sha3.Sum512([]byte(code))
	protectedCode = hex.EncodeToString(digest[:])
	strategyName = "sha3-512"

	return
}
