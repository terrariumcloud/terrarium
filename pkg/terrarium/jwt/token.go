package jwt

import (
	"encoding/json"
	"fmt"
)

type JWTToken struct {
	header    TokenComponent
	payload   TokenComponent
	signature TokenSignature
}

func (j *JWTToken) Header() *JWTHeader {
	h, ok := j.header.(*JWTHeader)
	if ok {
		return h
	}
	return nil
}

func (j *JWTToken) Payload() *JWTPayload {
	p, ok := j.payload.(*JWTPayload)
	if ok {
		return p
	}
	return nil
}

func (j *JWTToken) Signature() *JWTSignature {
	s, ok := j.signature.(*JWTSignature)
	if ok {
		return s
	}
	return nil
}

func (j *JWTToken) ToJSON() ([]byte, error) {
	var data map[string]interface{} = map[string]interface{}{
		"header":  j.Header(),
		"payload": j.Payload(),
	}
	return json.MarshalIndent(data, "", "   ")
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
