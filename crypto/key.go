package crypto

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/Click-CI/common/util"
)

type KeyPairFiles struct {
	Cert string `json:"cert,omitempty" yaml:"cert,omitempty" bson:"cert,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty" bson:"key,omitempty"`
}

type RsaKeyPair struct {
	Files    *KeyPairFiles     `json:"files,omitempty" yaml:"files,omitempty" bson:"files,omitempty"`
	Key      *rsa.PrivateKey   `json:"key,omitempty" yaml:"key,omitempty" bson:"key,omitempty"`
	Cert     *rsa.PublicKey    `json:"cert,omitempty" yaml:"cert,omitempty" bson:"cert,omitempty"`
	X509Cert *x509.Certificate `json:"x_509_cert,omitempty" yaml:"x_509_cert,omitempty" bson:"x_509_cert,omitempty"`
}

func NewRsaKeyPair(keyPairFiles KeyPairFiles) (*RsaKeyPair, *util.Result) {
	cert, err := tls.LoadX509KeyPair(keyPairFiles.Cert, keyPairFiles.Key)
	if err != nil {
		return nil, util.Error("LoadX509KeyPair", err)
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, util.Error("x509.ParseCertificate", err)
	}

	switch pub := x509Cert.PublicKey.(type) {
	case *rsa.PublicKey:
		priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, util.MsgError("ValidateKeyPair", "private key type does not match public key type")
		}
		if pub.N.Cmp(priv.N) != 0 {
			return nil, util.MsgError("ValidateKeyPair", "private key does not match public key")
		}
		return &RsaKeyPair{Files: &keyPairFiles, Key: priv, Cert: pub, X509Cert: x509Cert}, nil

	default:
		return nil, util.MsgError("ValidateKeyPair", "invalid public key algorithm")
	}
}

func (kp RsaKeyPair) GetKeyPair() (privateKey *rsa.PrivateKey, cert []byte, err error) {
	if kp.Key == nil {
		return nil, nil, fmt.Errorf("no private key in rsa key pair")
	}
	if kp.X509Cert == nil {
		return nil, nil, fmt.Errorf("no public key in rsa key pair")
	}
	return kp.Key, kp.X509Cert.Raw, nil
}

type KeyPair struct {
	PrivateKey string `json:"private_key,omitempty" yaml:"private_key,omitempty" bson:"private_key,omitempty"`
	PublicKey  string `json:"public_key,omitempty" yaml:"public_key,omitempty" bson:"public_key,omitempty"`
}

func (kp KeyPair) NewRsaKeyPair() (*RsaKeyPair, *util.Result) {
	cert, err := tls.X509KeyPair([]byte(kp.PublicKey), []byte(kp.PrivateKey))
	if err != nil {
		return nil, util.Error("X509KeyPair", err)
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, util.Error("x509.ParseCertificate", err)
	}

	switch pub := x509Cert.PublicKey.(type) {
	case *rsa.PublicKey:
		priv, ok := cert.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, util.MsgError("ValidateKeyPair", "private key type does not match public key type")
		}
		if pub.N.Cmp(priv.N) != 0 {
			return nil, util.MsgError("ValidateKeyPair", "private key does not match public key")
		}
		return &RsaKeyPair{Files: nil, Key: priv, Cert: pub, X509Cert: x509Cert}, nil

	default:
		return nil, util.MsgError("ValidateKeyPair", "invalid public key algorithm")
	}
}
