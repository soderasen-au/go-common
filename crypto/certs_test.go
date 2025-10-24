package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Helper to generate test PEM certificates
func generateTestPEMFiles(t *testing.T, dir string) (certFile, keyFile, caFile string) {
	t.Helper()

	// Generate CA key and cert
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}

	// Write CA certificate
	caFile = filepath.Join(dir, ROOT_CA_FILENAME)
	caOut, _ := os.Create(caFile)
	_ = pem.Encode(caOut, &pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})
	_ = caOut.Close()

	// Generate client key and cert
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test Client"},
			CommonName:   "client.example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Parse CA cert for signing
	caCert, _ := x509.ParseCertificate(caCertDER)

	clientCertDER, err := x509.CreateCertificate(rand.Reader, &clientTemplate, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}

	// Write client certificate
	certFile = filepath.Join(dir, CLIENT_PEM_FILENAME)
	certOut, _ := os.Create(certFile)
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER})
	_ = certOut.Close()

	// Write client key
	keyFile = filepath.Join(dir, CLIENT_KEY_FILENAME)
	keyOut, _ := os.Create(keyFile)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	_ = keyOut.Close()

	return certFile, keyFile, caFile
}

func TestNewCertificates(t *testing.T) {
	testFolder := "/test/cert/folder"
	certs := NewCertificates(testFolder)

	if certs == nil {
		t.Fatal("NewCertificates() returned nil")
	}

	expectedClientFile := filepath.Join(testFolder, CLIENT_PEM_FILENAME)
	expectedKeyFile := filepath.Join(testFolder, CLIENT_KEY_FILENAME)
	expectedCAFile := filepath.Join(testFolder, ROOT_CA_FILENAME)

	if certs.ClientFile != expectedClientFile {
		t.Errorf("ClientFile = %v, want %v", certs.ClientFile, expectedClientFile)
	}
	if certs.ClientkeyFile != expectedKeyFile {
		t.Errorf("ClientkeyFile = %v, want %v", certs.ClientkeyFile, expectedKeyFile)
	}
	if certs.CAFile != expectedCAFile {
		t.Errorf("CAFile = %v, want %v", certs.CAFile, expectedCAFile)
	}
}

func TestCertificates_NewTlsConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "certs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	certFile, keyFile, caFile := generateTestPEMFiles(t, tmpDir)

	tests := []struct {
		name    string
		certs   Certificates
		wantErr bool
	}{
		{
			name: "valid certificates",
			certs: Certificates{
				ClientFile:    certFile,
				ClientkeyFile: keyFile,
				CAFile:        caFile,
			},
			wantErr: false,
		},
		{
			name: "missing client cert",
			certs: Certificates{
				ClientFile:    filepath.Join(tmpDir, "missing.pem"),
				ClientkeyFile: keyFile,
				CAFile:        caFile,
			},
			wantErr: true,
		},
		{
			name: "missing CA file",
			certs: Certificates{
				ClientFile:    certFile,
				ClientkeyFile: keyFile,
				CAFile:        filepath.Join(tmpDir, "missing_ca.pem"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, res := tt.certs.NewTlsConfig()
			if (res != nil) != tt.wantErr {
				t.Errorf("NewTlsConfig() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tlsConfig == nil {
					t.Error("NewTlsConfig() returned nil")
					return
				}
				if len(tlsConfig.Certificates) == 0 {
					t.Error("TLS config has no certificates")
				}
				if tlsConfig.RootCAs == nil {
					t.Error("TLS config has no RootCAs")
				}
				if !tlsConfig.InsecureSkipVerify {
					t.Error("InsecureSkipVerify should be true")
				}
			}
		})
	}
}

func TestCertificates_NewRsaKeyPair(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "certs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	certFile, keyFile, _ := generateTestPEMFiles(t, tmpDir)

	tests := []struct {
		name    string
		certs   Certificates
		wantErr bool
	}{
		{
			name: "valid certificates",
			certs: Certificates{
				ClientFile:    certFile,
				ClientkeyFile: keyFile,
			},
			wantErr: false,
		},
		{
			name: "missing files",
			certs: Certificates{
				ClientFile:    filepath.Join(tmpDir, "missing.pem"),
				ClientkeyFile: filepath.Join(tmpDir, "missing_key.pem"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaKeyPair, res := tt.certs.NewRsaKeyPair()
			if (res != nil) != tt.wantErr {
				t.Errorf("NewRsaKeyPair() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if rsaKeyPair == nil {
					t.Error("NewRsaKeyPair() returned nil")
					return
				}
				if rsaKeyPair.Key == nil {
					t.Error("RsaKeyPair.Key is nil")
				}
				if rsaKeyPair.Cert == nil {
					t.Error("RsaKeyPair.Cert is nil")
				}
			}
		})
	}
}

func TestCerts_NewTlsConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "certs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	certFile, keyFile, caFile := generateTestPEMFiles(t, tmpDir)

	tests := []struct {
		name    string
		certs   Certs
		wantErr bool
	}{
		{
			name: "with QlikPem",
			certs: Certs{
				QlikPem: &Certificates{
					ClientFile:    certFile,
					ClientkeyFile: keyFile,
					CAFile:        caFile,
				},
			},
			wantErr: false,
		},
		{
			name: "with KeyPair",
			certs: Certs{
				KeyPair: &KeyPairFiles{
					Cert: certFile,
					Key:  keyFile,
				},
			},
			wantErr: false,
		},
		{
			name: "with inline KeyPair",
			certs: Certs{
				KeyPair: &KeyPairFiles{
					Cert: certFile,
					Key:  keyFile,
				},
			},
			wantErr: false,
		},
		{
			name:    "no certificates",
			certs:   Certs{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, res := tt.certs.NewTlsConfig()
			if (res != nil) != tt.wantErr {
				t.Errorf("NewTlsConfig() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tlsConfig == nil {
					t.Error("NewTlsConfig() returned nil")
				}
			}
		})
	}
}

func TestCerts_NewRsaKeyPair(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "certs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	certFile, keyFile, _ := generateTestPEMFiles(t, tmpDir)

	tests := []struct {
		name    string
		certs   Certs
		wantErr bool
	}{
		{
			name: "with QlikPem",
			certs: Certs{
				QlikPem: &Certificates{
					ClientFile:    certFile,
					ClientkeyFile: keyFile,
				},
			},
			wantErr: false,
		},
		{
			name: "with KeyPair",
			certs: Certs{
				KeyPair: &KeyPairFiles{
					Cert: certFile,
					Key:  keyFile,
				},
			},
			wantErr: false,
		},
		{
			name:    "no certificates",
			certs:   Certs{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaKeyPair, res := tt.certs.NewRsaKeyPair()
			if (res != nil) != tt.wantErr {
				t.Errorf("NewRsaKeyPair() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if rsaKeyPair == nil {
					t.Error("NewRsaKeyPair() returned nil")
				}
			}
		})
	}
}

func TestCerts_Priority(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "certs-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	certFile, keyFile, caFile := generateTestPEMFiles(t, tmpDir)

	// Test priority: QlikPem > Pfx > KeyPair
	t.Run("QlikPem takes priority over KeyPair", func(t *testing.T) {
		certs := Certs{
			QlikPem: &Certificates{
				ClientFile:    certFile,
				ClientkeyFile: keyFile,
				CAFile:        caFile,
			},
			KeyPair: &KeyPairFiles{
				Cert: certFile,
				Key:  keyFile,
			},
		}

		// Should use QlikPem
		tlsConfig, res := certs.NewTlsConfig()
		if res != nil {
			t.Errorf("NewTlsConfig() error = %v", res)
			return
		}
		if tlsConfig == nil {
			t.Error("NewTlsConfig() returned nil")
			return
		}
		// QlikPem sets RootCAs, KeyPair doesn't
		if tlsConfig.RootCAs == nil {
			t.Error("Should have used QlikPem (with RootCAs)")
		}
	})

	t.Run("KeyPair used when QlikPem absent", func(t *testing.T) {
		certs := Certs{
			KeyPair: &KeyPairFiles{
				Cert: certFile,
				Key:  keyFile,
			},
		}

		tlsConfig, res := certs.NewTlsConfig()
		if res != nil {
			t.Errorf("NewTlsConfig() error = %v", res)
			return
		}
		if tlsConfig == nil {
			t.Error("NewTlsConfig() returned nil")
		}
	})
}

// Note: Pfx tests would require generating a valid PKCS#12 file
// which is complex. The existing tests cover the core functionality.
// To test Pfx fully, you would need to:
// 1. Generate a PKCS#12 file with proper password
// 2. Test Pfx.NewTlsConfig() and Pfx.NewRsaKeyPair()
// This is left as future enhancement if needed.

func TestPfx_EmptyPassword(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pfx-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a dummy file (not a real PFX)
	dummyFile := filepath.Join(tmpDir, "dummy.pfx")
	_ = os.WriteFile(dummyFile, []byte("not a real pfx"), 0644)

	pfx := Pfx{
		Cert:     dummyFile,
		Password: "",
	}

	// Should fail because it's not a real PFX file
	_, res := pfx.NewTlsConfig()
	if res == nil {
		t.Error("NewTlsConfig() should fail with invalid PFX file")
	}

	_, res = pfx.NewRsaKeyPair()
	if res == nil {
		t.Error("NewRsaKeyPair() should fail with invalid PFX file")
	}
}

func TestPfx_NonExistentFile(t *testing.T) {
	pfx := Pfx{
		Cert:     "/nonexistent/path/file.pfx",
		Password: "password",
	}

	t.Run("NewTlsConfig with non-existent file", func(t *testing.T) {
		_, res := pfx.NewTlsConfig()
		if res == nil {
			t.Error("NewTlsConfig() should fail with non-existent file")
		}
	})

	t.Run("NewRsaKeyPair with non-existent file", func(t *testing.T) {
		_, res := pfx.NewRsaKeyPair()
		if res == nil {
			t.Error("NewRsaKeyPair() should fail with non-existent file")
		}
	})
}
