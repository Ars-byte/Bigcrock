package src

import (
	"os"
	"time"
)

type SortMode int

const (
	SortSizeDesc SortMode = iota
	SortSizeAsc
	SortNameAsc
	SortNameDesc
	SortCountDesc
)

type Item struct {
	Name       string
	Path       string
	IsDir      bool
	IsSymlink  bool
	Size       int64
	SizeKnown  bool
	IsScanning bool
	FileCount  int
	DirCount   int
	ModTime    time.Time
	Mode       os.FileMode
	ScanErr    error
}

type DirPreviewInfo struct {
	Path       string
	TotalSize  int64
	FileCount  int
	DirCount   int
	TopSubdirs []Item
	TopFiles   []Item
	Ready      bool
	Err        error
}

type FilePreviewInfo struct {
	Path       string
	Size       int64
	Mode       os.FileMode
	ModTime    time.Time
	Lines      []string
	TotalLines int
	IsBinary   bool
	MimeType   string
	HexDump    string
}

type ScanProgressMsg struct {
	Path      string
	ItemPath  string
	Size      int64
	FileCount int
	DirCount  int
	Err       error
}

type DirPreviewReadyMsg struct {
	Path string
	Info DirPreviewInfo
}

type FilePreviewReadyMsg struct {
	Path string
	Info FilePreviewInfo
}

type SpinnerTickMsg struct{}
