package http3

import (
	"crypto/tls"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics/pkg/config"
	"github.com/quic-s/quics/pkg/utils"
)

// connectRestHandler creates mux router and connect handler to router
func connectRestHandler() *mux.Router {
	r := mux.NewRouter()

	// add handler to router

	return r
}

// startHttp3Server starts HTTP/3 server
func startHttp3Server(r *mux.Router) {
	quicConfig := &quic.Config{}

	server := &http3.Server{
		Addr:       config.GetViperEnvVariables("REST_SERVER_ADDR"),
		QuicConfig: quicConfig,
		Handler:    r,
	}

	// get directory path for certification
	quicsDirPath := utils.GetQuicsDirPath()
	certFileDir := filepath.Join(quicsDirPath, config.GetViperEnvVariables("QUICS_CERT_NAME"))
	keyFileDir := filepath.Join(quicsDirPath, config.GetViperEnvVariables("QUICS_KEY_NAME"))

	// load the certificate and the key from the files
	_, err := tls.LoadX509KeyPair(certFileDir, keyFileDir)
	if err != nil {
		config.SecurityFiles()
	}

	go func() {
		log.Fatal(server.ListenAndServeTLS(certFileDir, keyFileDir))
	}()

	fmt.Println("HTTP/3 server started successfully.")
}
