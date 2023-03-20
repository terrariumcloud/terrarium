package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	privateKeyPath = "./key.pem"
	publicKeyPath  = "./key.pub"
	keySize        = 4096
)

type Token interface {
	Sign(secretPath string) (string, error)
	Signature() *JWTSignature
	Payload() *JWTPayload
	Header() *JWTHeader
	ToJSON() ([]byte, error)
}

type TokenComponent interface {
	Encode() (string, error)
}

type TokenSignature interface {
	Create(header string, payload string) (string, error)
	Verify() error
}

func NewJWT(requestedScopes []string) *JWTToken {
	if len(requestedScopes) == 0 {
		requestedScopes = append(requestedScopes, "read")
	}
	now := time.Now()
	expire := now.Add(time.Hour * 24).Unix()
	nowStr := strconv.Itoa(int(now.Unix()))
	expireStr := strconv.Itoa(int(expire))
	header := &JWTHeader{
		Algorithm: "RS256",
		Type:      "JWT",
	}
	payload := &JWTPayload{
		Issuer:   "terrarium",
		Subject:  "",
		Audience: "",
		Expiry:   expireStr,
		IssuedAt: nowStr,
		Scopes:   requestedScopes,
	}
	signature := &JWTSignature{}
	token := &JWTToken{
		header:    header,
		payload:   payload,
		signature: signature,
	}
	return token
}

func NewJWTFromString(jwt string) (*JWTToken, error) {
	jwtParts := strings.Split(jwt, ".")
	if len(jwtParts) != 3 {
		return nil, errors.New("invalid jwt")
	}
	rawHeader, err := base64.RawURLEncoding.DecodeString(jwtParts[0])
	if err != nil {
		return nil, err
	}
	rawPayload, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return nil, err
	}
	signature, err := JWTSignatureFromJWT(jwt)
	if err != nil {
		return nil, err
	}
	header := &JWTHeader{}
	payload := &JWTPayload{}
	err = json.Unmarshal(rawHeader, &header)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawPayload, &payload)
	if err != nil {
		return nil, err
	}
	return &JWTToken{
		header:    header,
		payload:   payload,
		signature: signature,
	}, nil
}

func createJWTHash(header, payload string) [32]byte {
	unsignedJWT := fmt.Sprintf("%s.%s", header, payload)
	return sha256.Sum256([]byte(unsignedJWT))
}

func encode(c interface{}) (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return encoded, nil
}
