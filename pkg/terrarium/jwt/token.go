package jwt

import (
	"fmt"
)

type JWTToken struct {
	header    TokenComponent
	payload   TokenComponent
	signature TokenSignature
}

func (j *JWTToken) Sign() (string, error) {
	header, err := j.header.Encode()
	if err != nil {
		return "", err
	}
	payload, err := j.payload.Encode()
	if err != nil {
		return "", err
	}
	signature, err := j.signature.Create(header, payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s.%s", header, payload, signature), nil
}

func (j *JWTToken) Verify() error {
	return nil
}
