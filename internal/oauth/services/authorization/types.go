package authorization

type JWTComponent interface {
	Header() *JWTHeader
	Payload() *JWTPayload
	EncodedHeader() ([]byte, error)
	EncodedPayload() ([]byte, error)
	EncodedJWT() ([]byte, error)
}

type JWTToken struct {
	header  *JWTHeader
	payload *JWTPayload
}

func (j *JWTToken) Header() *JWTHeader {
	return j.header
}

func (j *JWTToken) EncodedHeader() ([]byte, error) {
	return j.encode(j.header)
}

func (j *JWTToken) EncodedPayload() ([]byte, error) {
	return j.encode(j.payload)
}

func (j *JWTToken) Payload() *JWTPayload {
	return j.payload
}

func (j *JWTToken) encode(c interface{}) ([]byte, error) {
	return nil, nil
}

type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type JWTPayload struct{}

func NewJWT() *JWTToken {
	return &JWTToken{
		header: &JWTHeader{
			Algorithm: "RSA",
			Type:      "JWT",
		},
		payload: &JWTPayload{},
	}
}
