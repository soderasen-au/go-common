package saml

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Click-CI/common/util"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

type IdpMetaDataType int

const (
	IDP_METADATA_NONE IdpMetaDataType = 0
	IDP_METADATA_URL  IdpMetaDataType = 1
	IDP_METADATA_FILE IdpMetaDataType = 2
)

type IdpMetaData struct {
	Type    IdpMetaDataType         `json:"type,omitempty" yaml:"type,omitempty" bson:"type,omitempty"`
	File    *string                 `json:"file,omitempty" yaml:"file,omitempty" bson:"file,omitempty"`
	Url     *string                 `json:"url,omitempty" yaml:"url,omitempty" bson:"url,omitempty"`
	RawData []byte                  `json:"-" yaml:"-" bson:"-"`
	Data    *types.EntityDescriptor `json:"-" yaml:"-" bson:"-"`
}

func (meta IdpMetaData) HasData() bool {
	return meta.Type != IDP_METADATA_NONE
}

func (meta *IdpMetaData) parseRawData() (*types.EntityDescriptor, dsig.X509CertificateStore, *util.Result) {
	metadata := &types.EntityDescriptor{}
	err := xml.Unmarshal(meta.RawData, metadata)
	if err != nil {
		return nil, nil, util.Error("ParseRawData", err)
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				return nil, nil, util.MsgError("", fmt.Sprintf("metadata certificate(%d) must not be empty", idx))
			}
			certStr := strings.Join(strings.Fields(xcert.Data), "")
			certData, err := base64.StdEncoding.DecodeString(certStr)
			if err != nil {
				return nil, nil, util.Error(fmt.Sprintf("DecodeCertificate[%d]", idx), err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, nil, util.Error(fmt.Sprintf("ParseCertificate[%d]", idx), err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	return metadata, &certStore, nil
}

func (meta *IdpMetaData) GetMetaData() (*types.EntityDescriptor, dsig.X509CertificateStore, *util.Result) {
	var err error

	switch meta.Type {
	case IDP_METADATA_NONE:
		return nil, nil, nil
	case IDP_METADATA_FILE:
		if meta.File == nil {
			return nil, nil, util.MsgError("GetFileMetaData", "No filename")
		}

		meta.RawData, err = ioutil.ReadFile(*meta.File)
		if err != nil {
			return nil, nil, util.Error("ReadFileMetaData", err)
		}
		return meta.parseRawData()
	case IDP_METADATA_URL:
		if meta.Url == nil {
			return nil, nil, util.MsgError("GetUrlMetaData", "No URL")
		}
		response, err := http.Get(*meta.Url)
		if err != nil {
			return nil, nil, util.Error("GetHttpMetaData", err)
		}
		meta.RawData, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, util.Error("ReadUrlMetaData", err)
		}
		return meta.parseRawData()
	default:
		return nil, nil, util.MsgError("GetMetaData", "Unsupported metadata type")
	}
}
