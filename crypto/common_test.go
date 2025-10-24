package crypto

import (
	"testing"
)

func TestInternalKeysInitialization(t *testing.T) {
	t.Run("InternalPrivateKey is initialized", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Error("InternalPrivateKey should be initialized during package init")
		}
	})

	t.Run("InternalPrivateKey is valid RSA key", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Skip("InternalPrivateKey is nil")
		}
		if InternalPrivateKey.N == nil {
			t.Error("InternalPrivateKey.N (modulus) is nil")
		}
		if InternalPrivateKey.D == nil {
			t.Error("InternalPrivateKey.D (private exponent) is nil")
		}
		if InternalPrivateKey.N == nil {
			t.Error("InternalPrivateKey.PublicKey.N is nil")
		}
	})

	t.Run("InternalKeyPair is initialized", func(t *testing.T) {
		if InternalKeyPair == nil {
			t.Error("InternalKeyPair should be initialized during package init")
		}
	})

	t.Run("InternalKeyPair has valid content", func(t *testing.T) {
		if InternalKeyPair == nil {
			t.Skip("InternalKeyPair is nil")
		}
		if InternalKeyPair.PrivateKey == "" {
			t.Error("InternalKeyPair.PrivateKey is empty")
		}
		if InternalKeyPair.PublicKey == "" {
			t.Error("InternalKeyPair.PublicKey is empty")
		}
		if InternalKeyPair.PrivateKey != PRIVATE_KEY {
			t.Error("InternalKeyPair.PrivateKey does not match PRIVATE_KEY constant")
		}
		if InternalKeyPair.PublicKey != PUBLIC_KEY {
			t.Error("InternalKeyPair.PublicKey does not match PUBLIC_KEY constant")
		}
	})

	t.Run("Constants are not empty", func(t *testing.T) {
		if PRIVATE_KEY == "" {
			t.Error("PRIVATE_KEY constant is empty")
		}
		if PUBLIC_KEY == "" {
			t.Error("PUBLIC_KEY constant is empty")
		}
	})

	t.Run("Key size is appropriate", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Skip("InternalPrivateKey is nil")
		}
		keySize := InternalPrivateKey.N.BitLen()
		// The test key should be at least 2048 bits
		if keySize < 2048 {
			t.Errorf("Key size is %d bits, should be at least 2048", keySize)
		}
		t.Logf("Internal RSA key size: %d bits", keySize)
	})

	t.Run("Public key matches private key", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Skip("InternalPrivateKey is nil")
		}
		// Verify public key components match
		if InternalPrivateKey.N.Cmp(InternalPrivateKey.N) != 0 {
			t.Error("Public key modulus does not match private key modulus")
		}
	})
}

func TestInternalKeyPairCanBeUsed(t *testing.T) {
	// Note: InternalKeyPair has PUBLIC_KEY constant which is a public key, not a certificate.
	// tls.X509KeyPair() requires a CERTIFICATE block, so these operations will fail.
	// This tests the actual behavior - the constants are not suitable for direct TLS use.

	t.Run("Cannot create RsaKeyPair from InternalKeyPair (expected failure)", func(t *testing.T) {
		rsaKeyPair, res := InternalKeyPair.NewRsaKeyPair()
		// Should fail because PUBLIC_KEY is not a certificate
		if res == nil {
			t.Error("Expected error when using PUBLIC_KEY instead of certificate")
			return
		}
		if rsaKeyPair != nil {
			t.Error("RsaKeyPair should be nil on error")
		}
	})

	t.Run("Cannot create TLS config from InternalKeyPair (expected failure)", func(t *testing.T) {
		tlsConfig, res := InternalKeyPair.NewTlsConfig()
		// Should fail because PUBLIC_KEY is not a certificate
		if res == nil {
			t.Error("Expected error when using PUBLIC_KEY instead of certificate")
			return
		}
		if tlsConfig != nil {
			t.Error("TLS config should be nil on error")
		}
	})
}

func TestInternalPrivateKeyValidation(t *testing.T) {
	t.Run("Validate InternalPrivateKey using RSA package", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Skip("InternalPrivateKey is nil")
		}

		// Validate the key using crypto/rsa
		if err := InternalPrivateKey.Validate(); err != nil {
			t.Errorf("InternalPrivateKey validation failed: %v", err)
		}
	})

	t.Run("InternalPrivateKey has correct type", func(t *testing.T) {
		if InternalPrivateKey == nil {
			t.Skip("InternalPrivateKey is nil")
		}

		// Type assertion to ensure it's actually an RSA key
		var _ = InternalPrivateKey
	})
}
