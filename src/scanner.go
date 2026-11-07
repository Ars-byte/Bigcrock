package src

import (
	"bufio"
	"context"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type CachedDir struct {
	Size      int64
	FileCount int
	DirCount  int
	ModTime   time.Time
}

var (
	sizeCache sync.Map
	scanSem   = make(chan struct{}, 2)
)

func InvalidateCache(dirPath string) {
	abs, err := filepath.Abs(dirPath)
	if err == nil {
		sizeCache.Delete(abs)
	} else {
		sizeCache.Delete(dirPath)
	}
}

func IsPseudoFS(path string) bool {
	clean := filepath.Clean(path)
	if clean == "/proc" || strings.HasPrefix(clean, "/proc/") {
		return true
	}
	if clean == "/sys" || strings.HasPrefix(clean, "/sys/") {
		return true
	}
	if clean == "/dev" || strings.HasPrefix(clean, "/dev/") {
		return true
	}
	return false
}

func ReadDirEntries(dirPath string, showHidden bool) ([]Item, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, entry := range entries {
		name := entry.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)
		info, err := entry.Info()
		if err != nil {
			continue
		}

		isSymlink := (entry.Type() & os.ModeSymlink) != 0
		isDir := entry.IsDir()

		item := Item{
			Name:      name,
			Path:      fullPath,
			IsDir:     isDir,
			IsSymlink: isSymlink,
			ModTime:   info.ModTime(),
			Mode:      info.Mode(),
		}

		if !isDir {
			item.Size = info.Size()
			item.SizeKnown = true
		} else {
			if val, ok := sizeCache.Load(fullPath); ok {
				cached := val.(CachedDir)
				item.Size = cached.Size
				item.FileCount = cached.FileCount
				item.DirCount = cached.DirCount
				item.SizeKnown = true
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func CalculateDirSize(ctx context.Context, dirPath string) (int64, int, int, error) {
	if IsPseudoFS(dirPath) {
		return 0, 0, 0, nil
	}

	select {
	case <-ctx.Done():
		return 0, 0, 0, ctx.Err()
	case scanSem <- struct{}{}:
		defer func() { <-scanSem }()
	}

	var totalSize int64
	var fileCount int
	var dirCount int

	checkCount := 0

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		checkCount++
		if checkCount%1024 == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		if path == dirPath {
			return nil
		}

		if d.Type().IsRegular() {
			info, err := d.Info()
			if err == nil {
				totalSize += info.Size()
				fileCount++
			}
			return nil
		}

		if (d.Type() & os.ModeSymlink) != 0 {
			info, err := d.Info()
			if err == nil && !info.IsDir() {
				totalSize += info.Size()
				fileCount++
			}
			return nil
		}

		if d.IsDir() {
			if IsPseudoFS(path) {
				return filepath.SkipDir
			}
			dirCount++
			return nil
		}

		return nil
	})

	if err != nil && err != context.Canceled {
		return totalSize, fileCount, dirCount, err
	}

	if err == nil {
		info, statErr := os.Stat(dirPath)
		var mtime time.Time
		if statErr == nil {
			mtime = info.ModTime()
		}
		sizeCache.Store(dirPath, CachedDir{
			Size:      totalSize,
			FileCount: fileCount,
			DirCount:  dirCount,
			ModTime:   mtime,
		})
	}

	return totalSize, fileCount, dirCount, err
}

func GenerateDirPreview(dirPath string, maxSubdirs, maxFiles int) DirPreviewInfo {
	info := DirPreviewInfo{
		Path: dirPath,
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		info.Err = err
		return info
	}

	var subdirs []Item
	var files []Item

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}

		isDir := entry.IsDir()
		isSymlink := (entry.Type() & os.ModeSymlink) != 0

		item := Item{
			Name:      entry.Name(),
			Path:      fullPath,
			IsDir:     isDir,
			IsSymlink: isSymlink,
			ModTime:   entryInfo.ModTime(),
			Mode:      entryInfo.Mode(),
		}

		if isDir {
			if val, ok := sizeCache.Load(fullPath); ok {
				cached := val.(CachedDir)
				item.Size = cached.Size
				item.SizeKnown = true
			} else {
				item.Size = entryInfo.Size()
			}
			subdirs = append(subdirs, item)
		} else {
			item.Size = entryInfo.Size()
			item.SizeKnown = true
			files = append(files, item)
		}
	}

	sort.Slice(subdirs, func(i, j int) bool {
		return subdirs[i].Size > subdirs[j].Size
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(subdirs) > maxSubdirs {
		info.TopSubdirs = subdirs[:maxSubdirs]
	} else {
		info.TopSubdirs = subdirs
	}

	if len(files) > maxFiles {
		info.TopFiles = files[:maxFiles]
	} else {
		info.TopFiles = files
	}

	if val, ok := sizeCache.Load(dirPath); ok {
		cached := val.(CachedDir)
		info.TotalSize = cached.Size
		info.FileCount = cached.FileCount
		info.DirCount = cached.DirCount
	} else {
		var s int64
		for _, f := range files {
			s += f.Size
		}
		for _, d := range subdirs {
			s += d.Size
		}
		info.TotalSize = s
		info.FileCount = len(files)
		info.DirCount = len(subdirs)
	}

	info.Ready = true
	return info
}

func GenerateFilePreview(filePath string, maxLines int) FilePreviewInfo {
	info := FilePreviewInfo{
		Path: filePath,
	}

	fileStat, err := os.Stat(filePath)
	if err != nil {
		return info
	}

	info.Size = fileStat.Size()
	info.Mode = fileStat.Mode()
	info.ModTime = fileStat.ModTime()

	f, err := os.Open(filePath)
	if err != nil {
		return info
	}
	defer f.Close()

	header := make([]byte, 512)
	n, _ := f.Read(header)
	if n > 0 {
		info.MimeType = http.DetectContentType(header[:n])
		for i := 0; i < n; i++ {
			if header[i] == 0 {
				info.IsBinary = true
				break
			}
		}
	}

	if info.IsBinary {
		limit := n
		if limit > 128 {
			limit = 128
		}
		info.HexDump = hex.Dump(header[:limit])
		return info
	}

	if _, err := f.Seek(0, io.SeekStart); err == nil {
		scanner := bufio.NewScanner(f)
		count := 0
		for scanner.Scan() {
			count++
			if len(info.Lines) < maxLines {
				info.Lines = append(info.Lines, scanner.Text())
			}
		}
		info.TotalLines = count
	}

	return info
}
