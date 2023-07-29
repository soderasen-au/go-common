package oauth

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/Click-CI/common/crypto"
	"github.com/Click-CI/common/util"
	"github.com/dghubble/oauth1"
	"golang.org/x/oauth2"
)

type AuthMethod string

const (
	AM_OAuth_AuthCode  AuthMethod = "OAuthCode"
	AM_OAuth_CredToken AuthMethod = "OAuthCredToken"
	AM_AppPassword     AuthMethod = "AppPassword"
	AM_OAuth1_RSA      AuthMethod = "OAuth1RSA"
)

const AuthCodeStateName = "click-ci-core"

type User struct {
	UserName    *string `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	AccountID   *string `json:"account_id,omitempty"`
	UUID        *string `json:"uuid,omitempty"`
	ID          *int64  `json:"id,omitempty"`
}

type AuthInfo struct {
	Method          AuthMethod      `json:"method" yaml:"method"`
	ID              *string         `json:"id,omitempty" yaml:"id,omitempty"`
	Secret          *string         `json:"secret,omitempty" yaml:"secret,omitempty"`
	User            *string         `json:"user,omitempty" yaml:"user,omitempty"`
	Password        *string         `json:"password,omitempty" yaml:"password,omitempty"`
	RsaPEMKeyPair   *crypto.KeyPair `json:"rsa_pem_key_pair,omitempty" yaml:"rsa_pem_key_pair,omitempty"`
	RequestTokenURL *string         `json:"request_token_url,omitempty" yaml:"request_token_url,omitempty"`
	AuthorizeURL    *string         `json:"authorize_url,omitempty" yaml:"authorize_url,omitempty"`
	RedirectURL     *string         `json:"redirect_url,omitempty" yaml:"redirect_url,omitempty"`
	AccessTokenURL  *string         `json:"access_token_url,omitempty" yaml:"access_token_url,omitempty"`
}

func (a AuthInfo) IsTokenAuth() bool {
	return a.Method == AM_OAuth_AuthCode || a.Method == AM_OAuth_CredToken || a.Method == AM_OAuth1_RSA
}

type Config struct {
	Ctx           context.Context `json:"-" yaml:"-"`
	Cfg           *oauth2.Config  `json:"-" yaml:"-"`
	Token         *oauth2.Token   `json:"oauth_token,omitempty" yaml:"oauth_token,omitempty"`
	CfgVer1       *oauth1.Config  `json:"-" yaml:"-"`
	TokenVer1     *oauth1.Token   `json:"oauth_token_ver1,omitempty" yaml:"oauth_token_ver1,omitempty"`
	Auth          AuthInfo        `json:"auth_info" yaml:"auth_info"`
	RsaPEMKeyPair *crypto.KeyPair `json:"rsa_pem_key_pair,omitempty" yaml:"rsa_pem_key_pair,omitempty"`
}

func NewConfig(auth AuthInfo) *Config {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)
	return &Config{
		Ctx:  ctx,
		Auth: auth,
	}
}

func (c *Config) AuthRequest(r *http.Request) {
	if c.Auth.IsTokenAuth() {
		c.Token.SetAuthHeader(r)
	} else {
		r.SetBasicAuth(*c.Auth.User, *c.Auth.Password)
	}
}

func (c *Config) AuthCodeURL() (string, *util.Result) {
	if c.Auth.Method == AM_OAuth_AuthCode {
		return c.Cfg.AuthCodeURL(AuthCodeStateName), nil
	} else if c.Auth.Method == AM_OAuth1_RSA {
		requestToken, _, err := c.CfgVer1.RequestToken() //RSA Signer doesn't need requestSecret;
		if err != nil {
			return "", util.Error("OAuth1: RequestToken", err)
		}
		authorizationURL, err := c.CfgVer1.AuthorizationURL(requestToken)
		if err != nil {
			return "", util.Error("OAuth1: AuthorizationURL", err)
		}
		return authorizationURL.String(), nil
	}
	return "", util.MsgError("AuthCodeURL", "Unsupported auth method")
}

func (c *Config) Exchange(code string) (*oauth2.Token, error) {
	return c.Cfg.Exchange(c.Ctx, code)
}

func (c *Config) CredentialToken() (*oauth2.Token, error) {
	return c.Cfg.PasswordCredentialsToken(c.Ctx, *c.Auth.User, *c.Auth.Password)
}
