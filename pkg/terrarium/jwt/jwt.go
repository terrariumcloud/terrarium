package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

const (
	privateKeyPath = "./key.pem"
	publicKeyPath  = "./key.pub"
	keySize        = 4096
)

type Token interface {
	Verify() error
	Sign(secretPath string) (string, error)
}

type TokenComponent interface {
	Encode() (string, error)
}

type TokenSignature interface {
	Create(header string, payload string) (string, error)
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

func encode(c interface{}) (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return encoded, nil
}
