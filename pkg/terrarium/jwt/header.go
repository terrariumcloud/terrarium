package jwt

type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

func (h *JWTHeader) Encode() (string, error) {
	return encode(h)
}
