package saml

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Click-CI/common/loggers"
	"os"

	"github.com/Click-CI/common/util"
	saml2 "github.com/russellhaering/gosaml2"
	dsig "github.com/russellhaering/goxmldsig"
)

// for manual config of SAML
type IdpConfig struct {
	MetaData       *IdpMetaData              `json:"meta_data,omitempty" yaml:"meta_data,omitempty" bson:"meta_data,omitempty"`
	ServiceURL     string                    `json:"service_url,omitempty" yaml:"service_url,omitempty" bson:"service_url,omitempty"`
	ServiceBinding string                    `json:"service_binding,omitempty" yaml:"service_binding,omitempty" bson:"service_binding,omitempty"`
	SloUrl         string                    `json:"slo_url,omitempty" yaml:"slo_url,omitempty" bson:"slo_url,omitempty"`
	SloBinding     string                    `json:"slo_binding,omitempty" yaml:"slo_binding,omitempty" bson:"slo_binding,omitempty"`
	EntityID       string                    `json:"entity_id,omitempty" yaml:"entity_id,omitempty" bson:"entity_id,omitempty"`
	CertFile       *string                   `json:"cert_file,omitempty" yaml:"cert_file,omitempty" bson:"cert_file,omitempty"`
	CertStore      dsig.X509CertificateStore `json:"cert_store,omitempty" yaml:"cert_store,omitempty" bson:"cert_store,omitempty"`
}

func (c *IdpConfig) GetCertFromFile() *util.Result {
	if c.CertFile == nil {
		return util.MsgError("GetCertFromFile", "No idp cert file")
	}
	certPEMBlock, err := os.ReadFile(*c.CertFile)
	if err != nil {
		return util.Error("ReadCertFile", err)
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}
	idx := 0
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				return util.Error(fmt.Sprintf("ParseCertificate[%d]", idx), err)
			}
			certStore.Roots = append(certStore.Roots, x509Cert)
		} else {
			loggers.CoreDebugLogger.Info().Msgf("ParseCertificate[%d]: cert type (%s) is not unsupported", idx, certDERBlock.Type)
		}
		idx++
	}

	c.CertStore = &certStore
	return nil
}
func (c *IdpConfig) Init() *util.Result {
	if c.MetaData == nil {
		c.MetaData = &IdpMetaData{
			Type: IDP_METADATA_NONE,
		}
	}

	if c.MetaData.HasData() {
		metadata, certStore, res := c.MetaData.GetMetaData()
		if res != nil {
			return res.With("InitIdpConfig")
		}
		c.ServiceURL = metadata.IDPSSODescriptor.SingleSignOnServices[0].Location
		c.ServiceBinding = metadata.IDPSSODescriptor.SingleSignOnServices[0].Binding
		if len(metadata.IDPSSODescriptor.SingleLogoutServices) > 0 {
			c.SloUrl = metadata.IDPSSODescriptor.SingleLogoutServices[0].Location
			c.SloBinding = metadata.IDPSSODescriptor.SingleLogoutServices[0].Binding
		}
		c.EntityID = metadata.EntityID
		c.CertStore = certStore
		return nil
	}

	if len(c.ServiceBinding) == 0 || (c.ServiceBinding != saml2.BindingHttpRedirect && c.ServiceBinding != saml2.BindingHttpPost) {
		c.ServiceBinding = saml2.BindingHttpRedirect
	}
	if len(c.SloBinding) == 0 || (c.SloBinding != saml2.BindingHttpRedirect && c.SloBinding != saml2.BindingHttpPost) {
		c.SloBinding = saml2.BindingHttpRedirect
	}

	return c.GetCertFromFile()
}
