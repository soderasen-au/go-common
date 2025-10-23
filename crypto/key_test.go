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

// Helper function to generate test certificate and key files
func generateTestCertFiles(t *testing.T, dir string) (certFile, keyFile string) {
	t.Helper()

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			CommonName:   "test.example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Write certificate to file
	certFile = filepath.Join(dir, "test.pem")
	certOut, err := os.Create(certFile)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	certOut.Close()

	// Write private key to file
	keyFile = filepath.Join(dir, "test_key.pem")
	keyOut, err := os.Create(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyOut.Close()

	return certFile, keyFile
}

func TestKeyPairFiles_NewRsaKeyPair(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "crypto-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certFile, keyFile := generateTestCertFiles(t, tmpDir)

	tests := []struct {
		name    string
		kp      KeyPairFiles
		wantErr bool
	}{
		{
			name: "valid key pair files",
			kp: KeyPairFiles{
				Cert: certFile,
				Key:  keyFile,
			},
			wantErr: false,
		},
		{
			name: "non-existent cert file",
			kp: KeyPairFiles{
				Cert: filepath.Join(tmpDir, "nonexistent.pem"),
				Key:  keyFile,
			},
			wantErr: true,
		},
		{
			name: "non-existent key file",
			kp: KeyPairFiles{
				Cert: certFile,
				Key:  filepath.Join(tmpDir, "nonexistent_key.pem"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaKeyPair, res := tt.kp.NewRsaKeyPair()
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
				if rsaKeyPair.X509Cert == nil {
					t.Error("RsaKeyPair.X509Cert is nil")
				}
			}
		})
	}
}

func TestKeyPairFiles_NewTlsConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "crypto-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certFile, keyFile := generateTestCertFiles(t, tmpDir)

	tests := []struct {
		name    string
		kp      KeyPairFiles
		wantErr bool
	}{
		{
			name: "valid key pair files",
			kp: KeyPairFiles{
				Cert: certFile,
				Key:  keyFile,
			},
			wantErr: false,
		},
		{
			name: "invalid cert file",
			kp: KeyPairFiles{
				Cert: filepath.Join(tmpDir, "nonexistent.pem"),
				Key:  keyFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, res := tt.kp.NewTlsConfig()
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
				if !tlsConfig.InsecureSkipVerify {
					t.Error("InsecureSkipVerify should be true")
				}
			}
		})
	}
}

func TestNewRsaKeyPair(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "crypto-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certFile, keyFile := generateTestCertFiles(t, tmpDir)

	tests := []struct {
		name    string
		kpFiles KeyPairFiles
		wantErr bool
	}{
		{
			name: "valid files",
			kpFiles: KeyPairFiles{
				Cert: certFile,
				Key:  keyFile,
			},
			wantErr: false,
		},
		{
			name: "missing cert",
			kpFiles: KeyPairFiles{
				Cert: filepath.Join(tmpDir, "missing.pem"),
				Key:  keyFile,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaKeyPair, res := NewRsaKeyPair(tt.kpFiles)
			if (res != nil) != tt.wantErr {
				t.Errorf("NewRsaKeyPair() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if rsaKeyPair == nil {
					t.Error("NewRsaKeyPair() returned nil")
					return
				}
				if rsaKeyPair.Files == nil {
					t.Error("RsaKeyPair.Files should not be nil")
				}
			}
		})
	}
}

func TestRsaKeyPair_GetKeyPair(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "crypto-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	certFile, keyFile := generateTestCertFiles(t, tmpDir)
	kpFiles := KeyPairFiles{Cert: certFile, Key: keyFile}
	rsaKeyPair, res := NewRsaKeyPair(kpFiles)
	if res != nil {
		t.Fatalf("Failed to create RsaKeyPair: %v", res)
	}

	tests := []struct {
		name    string
		kp      RsaKeyPair
		wantErr bool
	}{
		{
			name:    "valid key pair",
			kp:      *rsaKeyPair,
			wantErr: false,
		},
		{
			name: "missing private key",
			kp: RsaKeyPair{
				Key:      nil,
				Cert:     rsaKeyPair.Cert,
				X509Cert: rsaKeyPair.X509Cert,
			},
			wantErr: true,
		},
		{
			name: "missing X509 cert",
			kp: RsaKeyPair{
				Key:      rsaKeyPair.Key,
				Cert:     rsaKeyPair.Cert,
				X509Cert: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, cert, err := tt.kp.GetKeyPair()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if privateKey == nil {
					t.Error("GetKeyPair() privateKey is nil")
				}
				if cert == nil {
					t.Error("GetKeyPair() cert is nil")
				}
			}
		})
	}
}

func TestKeyPair_NewRsaKeyPair(t *testing.T) {
	// Generate a valid certificate for testing
	tmpDir, _ := os.MkdirTemp("", "keypair-test-*")
	defer os.RemoveAll(tmpDir)
	certFile, keyFile := generateTestCertFiles(t, tmpDir)

	// Read the cert and key as strings
	certPEM, _ := os.ReadFile(certFile)
	keyPEM, _ := os.ReadFile(keyFile)

	tests := []struct {
		name    string
		kp      KeyPair
		wantErr bool
	}{
		{
			name: "valid key pair with certificate",
			kp: KeyPair{
				PrivateKey: string(keyPEM),
				PublicKey:  string(certPEM),
			},
			wantErr: false,
		},
		{
			name: "internal key pair fails - PUBLIC_KEY is not a certificate",
			kp: KeyPair{
				PrivateKey: PRIVATE_KEY,
				PublicKey:  PUBLIC_KEY,
			},
			wantErr: true, // PUBLIC_KEY is not a CERTIFICATE, so this fails
		},
		{
			name: "empty private key",
			kp: KeyPair{
				PrivateKey: "",
				PublicKey:  PUBLIC_KEY,
			},
			wantErr: true,
		},
		{
			name: "empty public key",
			kp: KeyPair{
				PrivateKey: PRIVATE_KEY,
				PublicKey:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid private key",
			kp: KeyPair{
				PrivateKey: "invalid key",
				PublicKey:  PUBLIC_KEY,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsaKeyPair, res := tt.kp.NewRsaKeyPair()
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
				if rsaKeyPair.Files != nil {
					t.Error("RsaKeyPair.Files should be nil for inline KeyPair")
				}
			}
		})
	}
}

func TestKeyPair_NewTlsConfig(t *testing.T) {
	// Generate a valid certificate for testing
	tmpDir, _ := os.MkdirTemp("", "keypair-test-*")
	defer os.RemoveAll(tmpDir)
	certFile, keyFile := generateTestCertFiles(t, tmpDir)

	// Read the cert and key as strings
	certPEM, _ := os.ReadFile(certFile)
	keyPEM, _ := os.ReadFile(keyFile)

	tests := []struct {
		name    string
		kp      KeyPair
		wantErr bool
	}{
		{
			name: "valid key pair with certificate",
			kp: KeyPair{
				PrivateKey: string(keyPEM),
				PublicKey:  string(certPEM),
			},
			wantErr: false,
		},
		{
			name: "internal key pair fails - PUBLIC_KEY is not a certificate",
			kp: KeyPair{
				PrivateKey: PRIVATE_KEY,
				PublicKey:  PUBLIC_KEY,
			},
			wantErr: true, // PUBLIC_KEY is not a CERTIFICATE
		},
		{
			name: "invalid key pair",
			kp: KeyPair{
				PrivateKey: "invalid",
				PublicKey:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, res := tt.kp.NewTlsConfig()
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
			}
		})
	}
}
