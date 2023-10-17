package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	http3 "github.com/quic-go/quic-go/http3"
	"github.com/quic-s/quics/pkg/config"
)

type RestClient struct {
	qconf        *quic.Config
	roundTripper *http3.RoundTripper
	hclient      *http.Client
}

func NewRestClient() *RestClient {
	quicConfig := &quic.Config{
		KeepAlivePeriod: 60,
	}

	restClient := &RestClient{
		qconf: quicConfig,
	}

	restClient.roundTripper = &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QuicConfig: restClient.qconf,
	}

	restClient.hclient = &http.Client{
		Transport: restClient.roundTripper,
	}

	return restClient
}

func (r *RestClient) GetRequest(path string) (*bytes.Buffer, error) {
	url := "https://" + config.GetViperEnvVariables("REST_SERVER_ADDR") + path

	rsp, err := r.hclient.Get(url)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	return body, nil
}

func (r *RestClient) PostRequest(path string, contentType string, content []byte) (*bytes.Buffer, error) {
	url := "https://" + config.GetViperEnvVariables("REST_SERVER_ADDR") + path

	contentReader := bytes.NewReader(content)
	rsp, err := r.hclient.Post(url, contentType, contentReader)
	if err != nil {
		log.Println("quics: ", err)
		return nil, err
	}

	buf := make([]byte, rsp.ContentLength)
	_, err = rsp.Body.Read(buf)
	if err != nil && err != io.EOF {
		log.Println("quics: ", err)
		return nil, err
	}

	if string(buf) != "" {
		log.Println("quis: ", string(buf))
	} else {
		log.Println("quis: ", "Success")
	}

	return nil, nil
}

func (r *RestClient) Close() error {
	r.hclient.CloseIdleConnections()

	err := r.roundTripper.Close()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}
