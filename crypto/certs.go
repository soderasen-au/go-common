package crypto

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"software.sslmate.com/src/go-pkcs12"

	"github.com/soderasen-au/go-common/util"
)

const (
	CLIENT_KEY_FILENAME string = "client_key.pem"
	CLIENT_PEM_FILENAME string = "client.pem"
	ROOT_CA_FILENAME    string = "root.pem"
)

type CertReader interface {
	NewTlsConfig() (*tls.Config, *util.Result)
	NewRsaKeyPair() (*RsaKeyPair, *util.Result)
}

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

func (certs Certificates) NewTlsConfig() (*tls.Config, *util.Result) {
	cert, err := tls.LoadX509KeyPair(certs.ClientFile, certs.ClientkeyFile)
	if err != nil {
		return nil, util.Error("LoadX509KeyPair", err)
	}

	caCert, err := os.ReadFile(certs.CAFile)
	if err != nil {
		return nil, util.Error("ReadCAFile", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, util.Error("AppendCertsFromPEM", err)
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

type Pfx struct {
	Cert     string `json:"cert" yaml:"cert"`
	Password string `json:"password" yaml:"password"`
}

func (p Pfx) NewTlsConfig() (*tls.Config, *util.Result) {
	pfxData, err := os.ReadFile(p.Cert)
	if err != nil {
		return nil, util.Error("ReadPfxFile", err)
	}

	key, cert, caCerts, err := pkcs12.DecodeChain(pfxData, p.Password)
	if err != nil {
		return nil, util.Error("DecodeChain", err)
	}

	caCertPool := x509.NewCertPool()
	chain := [][]byte{cert.Raw}
	for _, caCert := range caCerts {
		caCertPool.AddCert(caCert)
		chain = append(chain, caCert.Raw)
	}

	tlsCert := tls.Certificate{
		Certificate: chain,
		PrivateKey:  key,
		Leaf:        cert,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{tlsCert},
		RootCAs:            caCertPool,
	}

	return tlsConfig, nil
}

func (p Pfx) NewRsaKeyPair() (*RsaKeyPair, *util.Result) {
	pfxData, err := os.ReadFile(p.Cert)
	if err != nil {
		return nil, util.Error("ReadPfxFile", err)
	}

	key, cert, _, err := pkcs12.DecodeChain(pfxData, p.Password)
	if err != nil {
		return nil, util.Error("DecodeChain", err)
	}

	var rsaKey *rsa.PrivateKey
	switch key.(type) {
	case *rsa.PrivateKey:
		rsaKey = key.(*rsa.PrivateKey)
	default:
		return nil, util.MsgError("ParsePrivateKey", "private key is not rsa.PrivateKey")
	}
	var rsaPubKey *rsa.PublicKey
	switch cert.PublicKey.(type) {
	case *rsa.PublicKey:
		rsaPubKey = cert.PublicKey.(*rsa.PublicKey)
	default:
		return nil, util.MsgError("ParsePublicKey", "public key is not rsa.PublicKey")
	}

	rsaKeyPair := &RsaKeyPair{
		Files:    nil,
		Key:      rsaKey,
		Cert:     rsaPubKey,
		X509Cert: cert,
	}
	return rsaKeyPair, nil
}

type Certs struct {
	Pfx     *Pfx          `json:"pfx,omitempty" yaml:"pfx,omitempty"`
	QlikPem *Certificates `json:"qlik_pem,omitempty" yaml:"qlik_pem,omitempty"`
	KeyPair *KeyPairFiles `json:"key_pair,omitempty" yaml:"key_pair,omitempty"`
}

func (c Certs) NewTlsConfig() (*tls.Config, *util.Result) {
	if c.QlikPem != nil {
		return c.QlikPem.NewTlsConfig()
	}
	if c.Pfx != nil {
		return c.Pfx.NewTlsConfig()
	}
	if c.KeyPair != nil {
		return c.KeyPair.NewTlsConfig()
	}
	return nil, util.MsgError("", "there's no cert")
}

func (c Certs) NewRsaKeyPair() (*RsaKeyPair, *util.Result) {
	if c.QlikPem != nil {
		return c.QlikPem.NewRsaKeyPair()
	}
	if c.Pfx != nil {
		return c.Pfx.NewRsaKeyPair()
	}
	if c.KeyPair != nil {
		return c.KeyPair.NewRsaKeyPair()
	}
	return nil, util.MsgError("", "there's no cert")
}
