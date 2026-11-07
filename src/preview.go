package src

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func GetDirPreviewLines(info DirPreviewInfo, parentTotalSize int64, width, maxLines int, lang Language) []string {
	d := T(lang)
	var lines []string

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d0d0d0"))
	sizeStyle := lipgloss.NewStyle().Bold(true).Foreground(GetSizeColor(info.TotalSize))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b0b0b0"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	lines = append(lines, titleStyle.Render(TruncateString(filepath.Base(info.Path)+"/", width)))
	lines = append(lines, labelStyle.Render("Path: ")+valStyle.Render(TruncateString(info.Path, width-6)))

	meta := fmt.Sprintf("%s %s  %s %s  %s %s",
		labelStyle.Render("Size:"), sizeStyle.Render(FormatBytes(info.TotalSize)),
		labelStyle.Render("Files:"), valStyle.Render(FormatCount(info.FileCount)),
		labelStyle.Render("Dirs:"), valStyle.Render(FormatCount(info.DirCount)),
	)
	lines = append(lines, meta)
	lines = append(lines, "")

	if len(info.TopSubdirs) == 0 && len(info.TopFiles) == 0 {
		lines = append(lines, dimStyle.Render(d.EmptyDir))
		return lines
	}

	if len(info.TopSubdirs) > 0 {
		lines = append(lines, sectionStyle.Render(d.TopSubdirs))
		maxSubSize := info.TopSubdirs[0].Size
		for _, sub := range info.TopSubdirs {
			if len(lines) >= maxLines-4 {
				break
			}
			ratio := 0.0
			if maxSubSize > 0 {
				ratio = float64(sub.Size) / float64(maxSubSize)
			}
			bar := MakeProgressBar(ratio, 8)
			sizeS := lipgloss.NewStyle().Foreground(GetSizeColor(sub.Size)).Render(fmt.Sprintf("%9s", FormatBytes(sub.Size)))
			barS := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(bar)
			nameS := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).Render(TruncateString(sub.Name+"/", width-24))
			lines = append(lines, fmt.Sprintf("  %s %s %s", sizeS, barS, nameS))
		}
		lines = append(lines, "")
	}

	if len(info.TopFiles) > 0 {
		lines = append(lines, sectionStyle.Render(d.TopFiles))
		for _, f := range info.TopFiles {
			if len(lines) >= maxLines {
				break
			}
			sizeS := lipgloss.NewStyle().Foreground(GetSizeColor(f.Size)).Render(fmt.Sprintf("%9s", FormatBytes(f.Size)))
			nameS := lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc")).Render(TruncateString(f.Name, width-14))
			lines = append(lines, fmt.Sprintf("  %s %s", sizeS, nameS))
		}
	}

	return lines
}

func GetFilePreviewLines(info FilePreviewInfo, width, maxLines int, lang Language) []string {
	var lines []string

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d0d0d0"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b0b0b0"))
	lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	lines = append(lines, titleStyle.Render(TruncateString(filepath.Base(info.Path), width)))
	lines = append(lines, fmt.Sprintf("%s %s  %s %s",
		labelStyle.Render("Size:"), valStyle.Render(FormatBytes(info.Size)),
		labelStyle.Render("Perms:"), valStyle.Render(info.Mode.String()),
	))
	lines = append(lines, labelStyle.Render("Modified: ")+valStyle.Render(info.ModTime.Format("2006-01-02 15:04:05")))
	if info.MimeType != "" {
		lines = append(lines, labelStyle.Render("MIME: ")+valStyle.Render(info.MimeType))
	}
	lines = append(lines, "")

	if info.IsBinary {
		lines = append(lines, sectionStyle.Render("Binary / Hex:"))
		if info.HexDump != "" {
			hexLines := strings.Split(info.HexDump, "\n")
			for _, hl := range hexLines {
				if len(lines) >= maxLines {
					break
				}
				if hl != "" {
					lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Render(TruncateString(hl, width)))
				}
			}
		}
		return lines
	}

	lines = append(lines, sectionStyle.Render("Preview:"))
	for i, l := range info.Lines {
		if len(lines) >= maxLines {
			break
		}
		num := lineNumStyle.Render(fmt.Sprintf("%3d ", i+1))
		content := valStyle.Render(TruncateString(l, width-5))
		lines = append(lines, num+content)
	}

	return lines
}

func RenderDirPreview(info DirPreviewInfo, parentTotalSize int64, width, height int, lang Language) string {
	lines := GetDirPreviewLines(info, parentTotalSize, width, height, lang)
	return strings.Join(lines, "\n")
}

func RenderFilePreview(info FilePreviewInfo, width, height int, lang Language) string {
	lines := GetFilePreviewLines(info, width, height, lang)
	return strings.Join(lines, "\n")
}
