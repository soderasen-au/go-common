package crypto

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/soderasen-au/go-common/util"
)

func SHA2656Hex(buf []byte) (string, *util.Result) {
	hasher := sha256.New()
	_, err := hasher.Write(buf)
	if err != nil {
		return "", util.Error("SHA256Write", err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
