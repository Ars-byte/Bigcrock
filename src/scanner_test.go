package src

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateDirSize(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bigcrock_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	f1 := filepath.Join(tempDir, "file1.txt")
	if err := os.WriteFile(f1, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}

	f2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(f2, []byte("testing bigcrock"), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	size, fileCount, dirCount, err := CalculateDirSize(context.Background(), tempDir)
	if err != nil {
		t.Fatalf("CalculateDirSize returned error: %v", err)
	}

	expectedSize := int64(11 + 16)
	if size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}
	if fileCount != 2 {
		t.Errorf("expected 2 files, got %d", fileCount)
	}
	if dirCount != 1 {
		t.Errorf("expected 1 subdir, got %d", dirCount)
	}

	val, ok := sizeCache.Load(tempDir)
	if !ok {
		t.Fatalf("expected dir to be cached")
	}
	cached := val.(CachedDir)
	if cached.Size != expectedSize {
		t.Errorf("expected cached size %d, got %d", expectedSize, cached.Size)
	}
}

func TestPseudoFS(t *testing.T) {
	if !IsPseudoFS("/proc") {
		t.Errorf("expected /proc to be pseudo FS")
	}
	if !IsPseudoFS("/sys/kernel") {
		t.Errorf("expected /sys/kernel to be pseudo FS")
	}
	if !IsPseudoFS("/dev/null") {
		t.Errorf("expected /dev/null to be pseudo FS")
	}
	if IsPseudoFS("/home/ars") {
		t.Errorf("expected /home/ars not to be pseudo FS")
	}
}
