package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/fs"
	"io/ioutil"

	"github.com/Click-CI/common/util"
)

func InternalEncypt(text []byte) ([]byte, error) {
	cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &InternalPrivateKey.PublicKey, text, nil)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func InternalDecrypt(cipher []byte) ([]byte, error) {
	text, err := InternalPrivateKey.Decrypt(nil, cipher, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}
	return text, nil
}

type CipherFile struct {
	File string `json:"file" yaml:"file" bson:"file" csv:"file"`
	Text string `json:"-" yaml:"-" bson:"-" csv:"-"`
}

func NewCipherFile(file string) *CipherFile {
	return &CipherFile{
		File: file,
		Text: "",
	}
}

func (c CipherFile) WriteToFile() *util.Result {
	if c.Text == "" {
		return util.MsgError("Check", "empty text")
	}
	if c.File == "" {
		return util.MsgError("Check", "empty file name")
	}

	cipher, err := InternalEncypt([]byte(c.Text))
	if err != nil {
		return util.Error("Encrypt", err)
	}

	err = ioutil.WriteFile(c.File, cipher, fs.ModePerm)
	if err != nil {
		return util.Error("WriteFile", err)
	}

	return nil
}

func (c *CipherFile) ReadFromFile() *util.Result {
	if c.File == "" {
		return util.MsgError("Check", "empty file name")
	}

	pwbytes, err := ioutil.ReadFile(c.File)
	if err != nil {
		return util.Error("ReadFile", err)
	}
	pw, err := InternalDecrypt(pwbytes)
	if err != nil {
		return util.Error("InternalDecrypt", err)
	}

	c.Text = string(pw)
	return nil
}
