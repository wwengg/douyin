package main

import (
	"crypto/tls"
	"crypto/x509"
	"douyin/config"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"os"
)

func initializeCA(certificateFile string, privateKeyFile string) (*tls.Certificate, error) {
	var err error
	certificate, err := ioutil.ReadFile(certificateFile)
	if err != nil {
		return nil, err
	}

	privateKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	ca, err := tls.X509KeyPair(certificate, privateKey)
	if err != nil {
		return nil, err
	}

	return &ca, nil
}

func configureCA() {
	var err error
	certificateFile := config.Get("CERTIFICATE_FILE", "./certificates/proxy-ca.crt")
	privateKeyFile := config.Get("PRIVATE_KEY_FILE", "./certificates/proxy-ca.key")

	if _, err := os.Stat(certificateFile); os.IsNotExist(err) {
		log.Printf("Certificate file not found, bypassing CA configuration. File: %s", certificateFile)
		return
	}

	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		log.Printf("Private key file not found, bypassing CA configuration. File: %s", privateKeyFile)
		return
	}

	goproxyCa, err := initializeCA(certificateFile, privateKeyFile)
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		log.Printf("Error loading certificate, bypassing CA configuration. Error: %s", err)
		return
	}

	goproxy.GoproxyCa = *goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(goproxyCa)}
}
