package authorization

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"

	"github.com/terrariumcloud/terrarium/internal/oauth"
	"github.com/terrariumcloud/terrarium/internal/oauth/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/jwt"
	"google.golang.org/grpc"
)

const (
	privateKeyPath = "./key.pem"
	publicKeyPath  = "./key.pub"
	keySize        = 4096
)

type AuthorizationServer struct {
	services.UnimplementedAuthorizationServer
}

func (a *AuthorizationServer) CreatePKI() error {
	_, err := os.Stat(privateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return a.generateRSAKeys()
		}
		return err
	} else {
		log.Println("RSA keys already exist. Skipping")
	}
	var token jwt.Token
	token = jwt.NewJWT([]string{}, privateKeyPath)
	signed, err := token.Sign(privateKeyPath)
	if err != nil {
		return err
	}
	log.Printf("%v", signed)
	return nil
}

// RegisterWithServer Registers AuthorizationServer with grpc server
func (a *AuthorizationServer) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	services.RegisterAuthorizationServer(grpcServer, a)
	return nil
}

func (a *AuthorizationServer) CreateApplication(ctx context.Context, req *oauth.CreateApplicationRequest) (*oauth.ApplicationResponse, error) {
	return nil, nil
}

func (a *AuthorizationServer) UpdateApplication(ctx context.Context, req *oauth.UpdateApplicationRequest) (*oauth.ApplicationResponse, error) {
	return nil, nil
}

func (a *AuthorizationServer) DeleteApplication(ctx context.Context, req *oauth.DeleteApplicationRequest) (*oauth.ApplicationResponse, error) {
	return nil, nil
}

func (a *AuthorizationServer) RotateApplicationSecrets(ctx context.Context, req *oauth.RotateApplicationSecretsRequest) (*oauth.RotateApplicationSecretsResponse, error) {
	return nil, nil
}

func (a *AuthorizationServer) IssueJWTToken(ctx context.Context, req *oauth.IssueJWTTokenRequest) (*oauth.IssueJWTTokenResponse, error) {
	return nil, nil
}

func (a *AuthorizationServer) generateRSAKeys() error {
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
