package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

func CreatePKI() error {
	_, err := os.Stat(privateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return generateRSAKeys()
		}
		return err
	} else {
		log.Println("RSA keys already exist. Skipping")
	}
	return nil
}

func generateRSAKeys() error {
	// https://stackoverflow.com/questions/64104586/use-golang-to-get-rsa-key-the-same-way-openssl-genrsa
	log.Println("Creating RSA keys...")
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}
	privatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	publicPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
		},
	)
	err = ioutil.WriteFile(privateKeyPath, privatePEM, 0700)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(publicKeyPath, publicPEM, 0755)
	if err != nil {
		return err
	}
	return nil
}
