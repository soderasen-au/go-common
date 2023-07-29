package saml

import (
	"github.com/Click-CI/common/crypto"
	"github.com/Click-CI/common/util"
	dsig "github.com/russellhaering/goxmldsig"
)

type ACSBindingType int

const (
	ACS_HTTP_REDIRECT = 0
	ACS_HTTP_POST     = 1
)

type SpConfig struct {
	EntityID                  string               `json:"entity_id,omitempty" yaml:"entity_id,omitempty" bson:"entity_id,omitempty"`
	SignAuthnRequests         bool                 `json:"sign_authn_requests,omitempty" yaml:"sign_authn_requests,omitempty" bson:"sign_authn_requests,omitempty"`
	ACSServiceURL             string               `json:"acs_service_url,omitempty" yaml:"acs_service_url,omitempty" bson:"acs_service_url,omitempty"`
	ACSBinding                ACSBindingType       `json:"acs_binding,omitempty" yaml:"acs_binding,omitempty" bson:"acs_binding,omitempty"`
	CertFiles                 *crypto.KeyPairFiles `json:"cert_files,omitempty" yaml:"cert_files,omitempty" bson:"cert_files,omitempty"`
	CertStore                 dsig.X509KeyStore    `json:"cert_store,omitempty" yaml:"cert_store,omitempty" bson:"cert_store,omitempty"`
	ValidateResponseSignature bool                 `json:"validate_response_signature,omitempty" yaml:"validate_response_signature,omitempty" bson:"validate_response_signature,omitempty"`
}

func (c *SpConfig) GetCertFromFiles() *util.Result {
	if c.CertFiles == nil {
		return util.MsgError("GetCertFromFiles", "No Sp cert files")
	}
	rsaKeyPair, res := crypto.NewRsaKeyPair(*c.CertFiles)
	if res != nil {
		return res.With("NewRsaKeyPair")
	}

	c.CertStore = rsaKeyPair
	return nil
}
