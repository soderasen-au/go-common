package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Click-CI/common/util"
)

const (
	CLIENT_KEY_FILENAME string = "client_key.pem"
	CLIENT_PEM_FILENAME string = "client.pem"
	ROOT_CA_FILENAME    string = "root.pem"
)

// Certificates stores where to find certificates used to connect to Engine.
type Certificates struct {
	ClientFile    string `json:"client" yaml:"client"`         // "/client.pem"
	ClientkeyFile string `json:"client_key" yaml:"client_key"` //"/client_key.pem"
	CAFile        string `json:"root_ca" yaml:"root_ca"`       // "/root.pem"
}

// NewCertificates creates Certificates from folder which contains
// "client.pem", "client_key.pem" and "root.pem"
func NewCertificates(certFolder string) *Certificates {
	certs := Certificates{}
	certs.ClientFile = filepath.Join(certFolder, CLIENT_PEM_FILENAME)
	certs.ClientkeyFile = filepath.Join(certFolder, CLIENT_KEY_FILENAME)
	certs.CAFile = filepath.Join(certFolder, ROOT_CA_FILENAME)
	return &certs
}

func (certs Certificates) NewTlsConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certs.ClientFile, certs.ClientkeyFile)
	if err != nil {
		return nil, fmt.Errorf("can't load client certs: %s", err.Error())
	}

	caCert, err := ioutil.ReadFile(certs.CAFile)
	if err != nil {
		return nil, fmt.Errorf("can't read ca certs: %s", err.Error())
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf("can't load ca certs: %s", err.Error())
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
	}

	return tlsConfig, nil
}

func (certs Certificates) NewRsaKeyPair() (*RsaKeyPair, *util.Result) {
	return NewRsaKeyPair(KeyPairFiles{
		Cert: certs.ClientFile,
		Key:  certs.ClientkeyFile,
	})
}
