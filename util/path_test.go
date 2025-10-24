package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() string
		cleanup  func(string)
		expected bool
		wantErr  bool
	}{
		{
			name: "existing file",
			setup: func() string {
				tmpFile, _ := os.CreateTemp("", "test-*.txt")
				tmpFile.Close()
				return tmpFile.Name()
			},
			cleanup: func(path string) {
				os.Remove(path)
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "existing directory",
			setup: func() string {
				tmpDir, _ := os.MkdirTemp("", "test-*")
				return tmpDir
			},
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "non-existing path",
			setup: func() string {
				return "/tmp/nonexistent-path-" + filepath.Base(os.TempDir())
			},
			cleanup:  func(path string) {},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			exists, err := Exists(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if exists != tt.expected {
				t.Errorf("Exists() = %v, want %v", exists, tt.expected)
			}
		})
	}
}

func TestListSubFolders(t *testing.T) {
	// Create test directory structure
	tmpDir, err := os.MkdirTemp("", "test-subfolders-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create subdirectories
	os.Mkdir(filepath.Join(tmpDir, "folder1"), 0755)
	os.Mkdir(filepath.Join(tmpDir, "folder2"), 0755)
	// Create a file (should be ignored)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644)

	tests := []struct {
		name     string
		path     string
		expected []string
		wantErr  bool
	}{
		{
			name:     "directory with subfolders",
			path:     tmpDir,
			expected: []string{"folder1", "folder2"},
			wantErr:  false,
		},
		{
			name:     "non-existing directory",
			path:     "/tmp/nonexistent-dir-xyz",
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folders, err := ListSubFolders(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSubFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expected == nil && folders == nil {
				return
			}
			if len(folders) != len(tt.expected) {
				t.Errorf("ListSubFolders() = %v, want %v", folders, tt.expected)
				return
			}
			// Check all expected folders exist
			for _, expected := range tt.expected {
				found := false
				for _, folder := range folders {
					if folder == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ListSubFolders() missing expected folder %s", expected)
				}
			}
		})
	}
}

func TestListFiles(t *testing.T) {
	// Create test directory with files
	tmpDir, err := os.MkdirTemp("", "test-listfiles-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test files
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)
	// Create a subdirectory (should be ignored)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	tests := []struct {
		name        string
		folder      string
		expectCount int
		wantErr     bool
	}{
		{
			name:        "directory with files",
			folder:      tmpDir,
			expectCount: 2,
			wantErr:     false,
		},
		{
			name:        "non-existing directory",
			folder:      "/tmp/nonexistent-dir-xyz",
			expectCount: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, result := ListFiles(tt.folder)
			if (result != nil) != tt.wantErr {
				t.Errorf("ListFiles() error = %v, wantErr %v", result, tt.wantErr)
				return
			}
			if result == nil && len(files) != tt.expectCount {
				t.Errorf("ListFiles() returned %d files, want %d", len(files), tt.expectCount)
			}
		})
	}
}

func TestFilterFiles(t *testing.T) {
	// Create test directory with various files
	tmpDir, err := os.MkdirTemp("", "test-filterfiles-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("test2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "data.csv"), []byte("data"), 0644)

	// Create subdirectory with files
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)

	tests := []struct {
		name        string
		folder      string
		pattern     string
		expectCount int
		wantErr     bool
	}{
		{
			name:        "filter txt files",
			folder:      tmpDir,
			pattern:     "*.txt",
			expectCount: 3, // includes nested.txt
			wantErr:     false,
		},
		{
			name:        "filter csv files",
			folder:      tmpDir,
			pattern:     "*.csv",
			expectCount: 1,
			wantErr:     false,
		},
		{
			name:        "filter no matches",
			folder:      tmpDir,
			pattern:     "*.pdf",
			expectCount: 0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, result := FilterFiles(tt.folder, tt.pattern)
			if (result != nil) != tt.wantErr {
				t.Errorf("FilterFiles() error = %v, wantErr %v", result, tt.wantErr)
				return
			}
			if len(files) != tt.expectCount {
				t.Errorf("FilterFiles() returned %d files, want %d", len(files), tt.expectCount)
			}
		})
	}
}

func TestMaybeCreate(t *testing.T) {
	tmpBase, err := os.MkdirTemp("", "test-maybecreate-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpBase)

	tests := []struct {
		name    string
		folder  string
		wantErr bool
	}{
		{
			name:    "create single directory",
			folder:  filepath.Join(tmpBase, "newdir"),
			wantErr: false,
		},
		{
			name:    "create nested directories",
			folder:  filepath.Join(tmpBase, "level1", "level2", "level3"),
			wantErr: false,
		},
		{
			name:    "create already existing directory",
			folder:  tmpBase,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MaybeCreate(tt.folder)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaybeCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Verify directory was created
			if err == nil {
				if exists, _ := Exists(tt.folder); !exists {
					t.Errorf("MaybeCreate() did not create directory %s", tt.folder)
				}
			}
		})
	}
}

func TestMoveFile(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (src, dst string)
		cleanup func(src, dst string)
		wantErr bool
	}{
		{
			name: "move file to new location",
			setup: func() (string, string) {
				tmpFile, _ := os.CreateTemp("", "test-src-*.txt")
				tmpFile.WriteString("test content")
				tmpFile.Close()
				tmpDir, _ := os.MkdirTemp("", "test-dst-*")
				return tmpFile.Name(), filepath.Join(tmpDir, "moved.txt")
			},
			cleanup: func(src, dst string) {
				os.RemoveAll(filepath.Dir(dst))
			},
			wantErr: false,
		},
		{
			name: "move file to directory",
			setup: func() (string, string) {
				tmpFile, _ := os.CreateTemp("", "test-src-*.txt")
				tmpFile.WriteString("test content")
				tmpFile.Close()
				tmpDir, _ := os.MkdirTemp("", "test-dst-*")
				return tmpFile.Name(), tmpDir
			},
			cleanup: func(src, dst string) {
				os.RemoveAll(dst)
			},
			wantErr: false,
		},
		{
			name: "move file creates parent directory",
			setup: func() (string, string) {
				tmpFile, _ := os.CreateTemp("", "test-src-*.txt")
				tmpFile.WriteString("test content")
				tmpFile.Close()
				tmpBase, _ := os.MkdirTemp("", "test-dst-*")
				return tmpFile.Name(), filepath.Join(tmpBase, "newdir", "moved.txt")
			},
			cleanup: func(src, dst string) {
				os.RemoveAll(filepath.Dir(filepath.Dir(dst)))
			},
			wantErr: false,
		},
		{
			name: "error on empty source",
			setup: func() (string, string) {
				tmpDir, _ := os.MkdirTemp("", "test-dst-*")
				return "", tmpDir
			},
			cleanup: func(src, dst string) {
				os.RemoveAll(dst)
			},
			wantErr: true,
		},
		{
			name: "error on non-existent source",
			setup: func() (string, string) {
				tmpDir, _ := os.MkdirTemp("", "test-dst-*")
				return "/tmp/nonexistent-file-xyz.txt", tmpDir
			},
			cleanup: func(src, dst string) {
				os.RemoveAll(dst)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst := tt.setup()
			defer tt.cleanup(src, dst)

			result := MoveFile(src, dst)
			if (result != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", result, tt.wantErr)
				return
			}

			if result == nil {
				// Verify source no longer exists
				if exists, _ := Exists(src); exists {
					t.Errorf("MoveFile() did not remove source file %s", src)
				}
				// Verify destination exists
				dstPath := dst
				if info, err := os.Stat(dst); err == nil && info.IsDir() {
					dstPath = filepath.Join(dst, filepath.Base(src))
				}
				if exists, _ := Exists(dstPath); !exists {
					t.Errorf("MoveFile() did not create destination file %s", dstPath)
				}
			}
		})
	}
}
