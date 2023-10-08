package main

import (
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// getHttp3Client returns created HTTP/3 client
func getHttp3Client() *http.Client {
	quicConfig := &quic.Config{}

	client := &http.Client{
		Transport: &http3.RoundTripper{
			QuicConfig: quicConfig,
		},
	}

	return client
}
