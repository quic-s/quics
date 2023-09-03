package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/quic-s/quics/config"
	"path/filepath"

	"log"
	"math/big"
	"os"
	"time"
)

func SecurityFiles() {

	// generate a new RSA key pair
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// create a template for the certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
	}

	// create a self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	// encode the certificate, key to PEM format
	certOut := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyOut := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// write the certificate and key to disk
	quicsDir := config.GetDirPath()
	certFile, err := os.Create(filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_CERT_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer func(certFile *os.File) {
		err := certFile.Close()
		if err != nil {
			log.Fatalf("Error while closing certFile: %s", err)
		}
	}(certFile)

	keyFile, err := os.Create(filepath.Join(quicsDir, config.GetViperEnvVariables("QUICS_KEY_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer func(keyFile *os.File) {
		err := keyFile.Close()
		if err != nil {
			log.Fatalf("Error while closing keyFile: %s", err)
		}
	}(keyFile)

	if _, err := certFile.Write(certOut); err != nil {
		log.Fatal(err)
	}

	if _, err := keyFile.Write(keyOut); err != nil {
		log.Fatal(err)
	}
}
