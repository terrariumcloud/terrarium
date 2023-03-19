package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

type JWTSignature struct{}

func (s *JWTSignature) Create(header string, payload string, secretPath string) (string, error) {
	unsignedJWT := fmt.Sprintf("%s.%s", header, payload)
	hash := sha256.Sum256([]byte(unsignedJWT))
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(secret)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	signature, err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(signature), nil
}
