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

type JWTSignature struct {
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	jwtHash      [32]byte
	rawSignature []byte
}

func (s *JWTSignature) Create(header string, payload string) (string, error) {
	unsignedJWT := fmt.Sprintf("%s.%s", header, payload)
	hash := sha256.Sum256([]byte(unsignedJWT))
	err := s.readPrivateKey()
	if err != nil {
		return "", err
	}
	signature, err := rsa.SignPKCS1v15(nil, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	s.rawSignature = signature
	s.jwtHash = hash
	return base64.RawURLEncoding.EncodeToString(signature), nil
}

func (s *JWTSignature) Verify() error {
	err := s.readPublicKey()
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, s.jwtHash[:], s.rawSignature)
}

func (s *JWTSignature) readPrivateKey() error {
	secret, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(secret)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	s.privateKey = key
	return nil
}

func (s *JWTSignature) readPublicKey() error {
	publicSecret, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}
	publicBlock, _ := pem.Decode(publicSecret)
	publicKey, _ := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	s.publicKey = publicKey
	return nil
}
