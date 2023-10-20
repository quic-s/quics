package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"path/filepath"

	"github.com/quic-s/quics/pkg/utils"

	"log"
	"math/big"
	"os"
)

// SecurityFiles generates a certificate file and key pair
func CreateSecurityFiles() error {

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return err
	}

	// encode the certificate, key to PEM format
	certOut := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyOut := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	// write the certificate and key to disk
	quicsDir := utils.GetQuicsDirPath()
	certFile, err := os.Create(filepath.Join(quicsDir, GetViperEnvVariables("QUICS_CERT_NAME")))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(certFile *os.File) {
		err := certFile.Close()
		if err != nil {
			log.Fatalf("Error while closing certFile: %s", err)
		}
	}(certFile)

	keyFile, err := os.Create(filepath.Join(quicsDir, GetViperEnvVariables("QUICS_KEY_NAME")))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(keyFile *os.File) {
		err := keyFile.Close()
		if err != nil {
			log.Fatalf("Error while closing keyFile: %s", err)
		}
	}(keyFile)

	if _, err := certFile.Write(certOut); err != nil {
		log.Fatal(err)
		return err
	}

	if _, err := keyFile.Write(keyOut); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
