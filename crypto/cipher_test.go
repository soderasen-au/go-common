package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInternalEncrypt(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		wantErr bool
	}{
		{
			name:    "empty text",
			text:    []byte{},
			wantErr: false,
		},
		{
			name:    "simple text",
			text:    []byte("hello world"),
			wantErr: false,
		},
		{
			name:    "longer text",
			text:    []byte("This is a longer text to test encryption functionality"),
			wantErr: false,
		},
		{
			name:    "binary data",
			text:    []byte{0x00, 0x01, 0x02, 0xff, 0xfe},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := InternalEncypt(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("InternalEncypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cipher == nil {
					t.Error("InternalEncypt() returned nil cipher")
				}
				// Cipher should be different from plaintext
				if len(tt.text) > 0 && string(cipher) == string(tt.text) {
					t.Error("Cipher should be different from plaintext")
				}
			}
		})
	}
}

func TestInternalDecrypt(t *testing.T) {
	plaintext := []byte("test message")

	// First encrypt
	cipher, err := InternalEncypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	tests := []struct {
		name    string
		cipher  []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid cipher",
			cipher:  cipher,
			want:    plaintext,
			wantErr: false,
		},
		{
			name:    "invalid cipher - random bytes",
			cipher:  []byte("invalid cipher text"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty cipher",
			cipher:  []byte{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := InternalDecrypt(tt.cipher)
			if (err != nil) != tt.wantErr {
				t.Errorf("InternalDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if string(text) != string(tt.want) {
					t.Errorf("InternalDecrypt() = %v, want %v", string(text), string(tt.want))
				}
			}
		})
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "simple text",
			plaintext: []byte("hello"),
		},
		{
			name:      "longer text",
			plaintext: []byte("The quick brown fox jumps over the lazy dog"),
		},
		{
			name:      "unicode text",
			plaintext: []byte("你好世界 🌍"),
		},
		{
			name:      "binary data",
			plaintext: []byte{0x00, 0x01, 0x02, 0x03, 0xff, 0xfe, 0xfd},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			cipher, err := InternalEncypt(tt.plaintext)
			if err != nil {
				t.Fatalf("InternalEncypt() error = %v", err)
			}

			// Decrypt
			decrypted, err := InternalDecrypt(cipher)
			if err != nil {
				t.Fatalf("InternalDecrypt() error = %v", err)
			}

			// Verify
			if string(decrypted) != string(tt.plaintext) {
				t.Errorf("Round trip failed: got %v, want %v", string(decrypted), string(tt.plaintext))
			}
		})
	}
}

func TestNewCipherFile(t *testing.T) {
	file := "test.enc"
	cf := NewCipherFile(file)

	if cf == nil {
		t.Error("NewCipherFile() returned nil")
		return
	}
	if cf.File != file {
		t.Errorf("NewCipherFile().File = %v, want %v", cf.File, file)
	}
	if cf.Text != "" {
		t.Errorf("NewCipherFile().Text = %v, want empty string", cf.Text)
	}
}

func TestCipherFile_WriteToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cipher-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		cf      CipherFile
		wantErr bool
	}{
		{
			name: "valid text and file",
			cf: CipherFile{
				File: filepath.Join(tmpDir, "test1.enc"),
				Text: "secret password",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			cf: CipherFile{
				File: filepath.Join(tmpDir, "test2.enc"),
				Text: "",
			},
			wantErr: true,
		},
		{
			name: "empty filename",
			cf: CipherFile{
				File: "",
				Text: "secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.cf.WriteToFile()
			if (res != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify file exists
				if _, err := os.Stat(tt.cf.File); os.IsNotExist(err) {
					t.Errorf("WriteToFile() did not create file %s", tt.cf.File)
				}
			}
		})
	}
}

func TestCipherFile_ReadFromFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cipher-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test encrypted file
	testFile := filepath.Join(tmpDir, "test.enc")
	testText := "my secret password"
	cf := CipherFile{File: testFile, Text: testText}
	if res := cf.WriteToFile(); res != nil {
		t.Fatalf("Failed to create test file: %v", res)
	}

	tests := []struct {
		name     string
		cf       *CipherFile
		wantText string
		wantErr  bool
	}{
		{
			name:     "valid encrypted file",
			cf:       &CipherFile{File: testFile},
			wantText: testText,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			cf:       &CipherFile{File: filepath.Join(tmpDir, "nonexistent.enc")},
			wantText: "",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			cf:       &CipherFile{File: ""},
			wantText: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.cf.ReadFromFile()
			if (res != nil) != tt.wantErr {
				t.Errorf("ReadFromFile() error = %v, wantErr %v", res, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.cf.Text != tt.wantText {
					t.Errorf("ReadFromFile() Text = %v, want %v", tt.cf.Text, tt.wantText)
				}
			}
		})
	}
}

func TestCipherFile_WriteReadRoundTrip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cipher-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		text string
	}{
		{
			name: "simple password",
			text: "password123",
		},
		{
			name: "long secret",
			text: "This is a very long secret with special chars !@#$%^&*()",
		},
		{
			name: "unicode text",
			text: "密码 パスワード مرور",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(tmpDir, tt.name+".enc")

			// Write
			writeCf := CipherFile{File: file, Text: tt.text}
			if res := writeCf.WriteToFile(); res != nil {
				t.Fatalf("WriteToFile() error = %v", res)
			}

			// Read
			readCf := CipherFile{File: file}
			if res := readCf.ReadFromFile(); res != nil {
				t.Fatalf("ReadFromFile() error = %v", res)
			}

			// Verify
			if readCf.Text != tt.text {
				t.Errorf("Round trip failed: got %v, want %v", readCf.Text, tt.text)
			}
		})
	}
}

func TestCipherFile_InvalidCipherData(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cipher-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with invalid cipher data
	invalidFile := filepath.Join(tmpDir, "invalid.enc")
	err = os.WriteFile(invalidFile, []byte("this is not encrypted data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cf := CipherFile{File: invalidFile}
	res := cf.ReadFromFile()
	if res == nil {
		t.Error("ReadFromFile() should fail with invalid cipher data")
	}
}
