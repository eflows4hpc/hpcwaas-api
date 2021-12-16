package rest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func (s *Server) createKey(gc *gin.Context) {
	userName := gc.Param("user_name")

	privateKey, publicKey, err := generateSSHKeyPair(2048, true)
	if err != nil {
		writeError(gc, newInternalServerError(err))
		return
	}

	err = s.vaultManager.StoreKV(fmt.Sprintf("/secret/data/ssh-credentials/%s", userName), map[string]interface{}{"data": map[string]interface{}{"userName": userName, "privateKey": string(privateKey), "publicKey": string(publicKey)}})
	if err != nil {
		log.Printf("Error storing key in vault: %+v", err)
		writeError(gc, newInternalServerError(err))
		return
	}

	gc.JSON(http.StatusCreated, api.SSHKey{PublicKey: string(publicKey)})
}

func generateSSHKeyPair(size int, useRSA bool) ([]byte, []byte, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		return nil, nil, err
	}
	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}

	publicRsaKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return nil, nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return encodePrivateKeyToPEM(privateKey), pubKeyBytes, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}
