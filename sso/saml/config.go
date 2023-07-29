package saml

import "github.com/Click-CI/common/util"

type Config struct {
	Idp *IdpConfig `json:"idp,omitempty" yaml:"idp,omitempty" bson:"idp,omitempty"`
	Sp  *SpConfig  `json:"sp,omitempty" yaml:"sp,omitempty" bson:"sp,omitempty"`
}

func (c *Config) Init() *util.Result {
	var res *util.Result

	if c.Idp == nil {
		return util.MsgError("SAMLConfigInit", "No Idp config")
	}
	res = c.Idp.Init()
	if res != nil {
		return res.With("SAML Idp.Init()")
	}

	if c.Sp == nil {
		return util.MsgError("SAMLConfigInit", "No Sp config")
	}
	res = c.Sp.GetCertFromFiles()
	if res != nil {
		return res.With("SAML Sp.GetCertFromFiles()")
	}

	return nil
}
