package jwt

type JWTPayload struct {
	Issuer   string   `json:"iss"`
	Subject  string   `json:"subject"`
	Audience string   `json:"aud"`
	Expiry   string   `json:"exp"`
	IssuedAt string   `json:"iat"`
	Scopes   []string `json:"scopes"`
}

func (p *JWTPayload) Encode() (string, error) {
	return encode(p)
}
