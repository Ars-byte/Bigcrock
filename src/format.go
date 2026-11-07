package src

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func FitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	w := ansi.StringWidth(s)
	if w > width {
		return ansi.Truncate(s, width, "")
	}
	if w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}

func FormatBytes(b int64) string {
	const unit = 1024
	if b < 0 {
		return "0 B"
	}
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	val := float64(b) / float64(div)
	if val >= 100 {
		return fmt.Sprintf("%.1f %s", val, units[exp])
	}
	return fmt.Sprintf("%.2f %s", val, units[exp])
}

func FormatBytesPadded(b int64) string {
	raw := FormatBytes(b)
	return fmt.Sprintf("%9s", raw)
}

func FormatCount(n int) string {
	if n < 0 {
		return "0"
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var res strings.Builder
	remainder := len(s) % 3
	if remainder > 0 {
		res.WriteString(s[:remainder])
		if len(s) > remainder {
			res.WriteString(",")
		}
	}
	for i := remainder; i < len(s); i += 3 {
		res.WriteString(s[i : i+3])
		if i+3 < len(s) {
			res.WriteString(",")
		}
	}
	return res.String()
}

func MakeProgressBar(ratio float64, barWidth int) string {
	if barWidth <= 0 {
		return ""
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	totalBlocks := float64(barWidth) * ratio
	fullBlocks := int(totalBlocks)
	if fullBlocks > barWidth {
		fullBlocks = barWidth
	}
	emptyBlocks := barWidth - fullBlocks

	var sb strings.Builder
	for i := 0; i < fullBlocks; i++ {
		sb.WriteString("█")
	}
	for i := 0; i < emptyBlocks; i++ {
		sb.WriteString("░")
	}
	return sb.String()
}

func GetSizeColor(size int64) lipgloss.Color {
	const (
		GB10  = 10 * 1024 * 1024 * 1024
		GB1   = 1 * 1024 * 1024 * 1024
		MB100 = 100 * 1024 * 1024
		MB10  = 10 * 1024 * 1024
	)

	switch {
	case size >= GB10:
		return lipgloss.Color("#ffffff")
	case size >= GB1:
		return lipgloss.Color("#e0e0e0")
	case size >= MB100:
		return lipgloss.Color("#b8b8b8")
	case size >= MB10:
		return lipgloss.Color("#909090")
	default:
		return lipgloss.Color("#686868")
	}
}

func GetItemBadge(item Item) (string, lipgloss.Style) {
	if item.IsDir {
		if strings.HasPrefix(item.Name, ".") {
			return "[HIDE]", lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
		}
		return "[DIR ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)
	}

	if item.IsSymlink {
		return "[LINK]", lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa"))
	}

	ext := strings.ToLower(filepath.Ext(item.Name))
	switch ext {
	case ".tar", ".gz", ".xz", ".zip", ".bz2", ".7z", ".rar", ".iso":
		return "[ARCH]", lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0")).Bold(true)
	case ".go", ".c", ".cpp", ".rs", ".py", ".js", ".ts", ".sh", ".bash":
		return "[CODE]", lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)
	case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".mp3", ".flac", ".wav":
		return "[MEDIA]", lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return "[IMG ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	case ".pdf", ".epub", ".mobi", ".docx", ".odt":
		return "[DOC ]", lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0"))
	case ".json", ".yaml", ".yml", ".toml", ".xml", ".csv":
		return "[DATA]", lipgloss.NewStyle().Foreground(lipgloss.Color("#b8b8b8"))
	default:
		return "[FILE]", lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	}
}

func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}
