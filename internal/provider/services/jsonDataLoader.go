package services

import (
	"encoding/json"
	"log"
	"os"
)

type LoadJSONData interface {
	LoadData() (map[string]ProviderData, error)
}

// Structs to load data into (from a JSON file for now, will be from DB later)

type GPGPublicKey struct {
	KeyID          string `json:"key_id"`
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceURL      string `json:"source_url"`
}

type SigningKeys struct {
	GPGPublicKeys []GPGPublicKey `json:"gpg_public_keys"`
}

type ProviderMetadata struct {
	OS                  string      `json:"os"`
	Arch                string      `json:"arch"`
	Filename            string      `json:"filename"`
	DownloadURL         string      `json:"download_url"`
	ShasumsURL          string      `json:"shasums_url"`
	ShasumsSignatureURL string      `json:"shasums_signature_url"`
	Shasum              string      `json:"shasum"`
	SigningKeys         SigningKeys `json:"signing_keys"`
}

type VersionData struct {
	Protocols []string           `json:"protocols"`
	Platforms []ProviderMetadata `json:"platforms"`
}

type ProviderData map[string]VersionData

var providerObj map[string]ProviderData

func LoadData() (map[string]ProviderData, error) {
	log.Printf("loadJSONDataIntoStructs")

	filePath := "./providers.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &providerObj); err != nil {
		return nil, err
	}
	return providerObj, nil
}
