package src

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var spinnerFrames = []string{"-", "\\", "|", "/"}

type Model struct {
	currDir      string
	parentDir    string
	parentItems  []Item
	parentCursor int

	items  []Item
	cursor int
	scroll int

	dirPreview  DirPreviewInfo
	filePreview FilePreviewInfo

	sortMode   SortMode
	showHidden bool
	lang       Language
	showAbout  bool

	scanning   bool
	scanTotal  int
	scanDone   int
	scanCtx    context.Context
	scanCancel context.CancelFunc

	spinnerIdx int
	width      int
	height     int
	err        string
}

func NewModel(initialDir string) Model {
	absDir, err := filepath.Abs(initialDir)
	if err != nil {
		absDir = initialDir
	}

	return Model{
		currDir:    absDir,
		sortMode:   SortSizeDesc,
		showHidden: true,
		lang:       LangEN,
		width:      80,
		height:     24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadDirectoryCmd(m.currDir),
		tickSpinner(),
	)
}

func tickSpinner() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

func (m *Model) loadDirectoryCmd(targetDir string) tea.Cmd {
	m.currDir = targetDir
	m.parentDir = filepath.Dir(targetDir)
	m.err = ""
	m.cursor = 0
	m.scroll = 0

	if m.scanCancel != nil {
		m.scanCancel()
	}
	m.scanCtx, m.scanCancel = context.WithCancel(context.Background())

	m.loadParentItems()

	items, err := ReadDirEntries(m.currDir, m.showHidden)
	if err != nil {
		m.err = fmt.Sprintf(T(m.lang).ErrReadDir, err)
		m.items = nil
		m.scanning = false
		return nil
	}

	m.items = items
	m.sortItems()

	var cmds []tea.Cmd
	uncalculated := 0
	for _, it := range m.items {
		if it.IsDir && !it.SizeKnown {
			uncalculated++
			cmds = append(cmds, scanItemCmd(m.scanCtx, m.currDir, it.Path))
		}
	}

	m.scanTotal = uncalculated
	m.scanDone = 0
	m.scanning = uncalculated > 0

	if m.scanning {
		cmds = append(cmds, tickSpinner())
	}

	m.updatePreview()

	return tea.Batch(cmds...)
}

func scanItemCmd(ctx context.Context, dirPath, itemPath string) tea.Cmd {
	return func() tea.Msg {
		size, files, dirs, err := CalculateDirSize(ctx, itemPath)
		return ScanProgressMsg{
			Path:      dirPath,
			ItemPath:  itemPath,
			Size:      size,
			FileCount: files,
			DirCount:  dirs,
			Err:       err,
		}
	}
}

func (m *Model) loadParentItems() {
	m.parentItems = nil
	m.parentCursor = -1

	if m.currDir == "/" {
		return
	}

	items, err := ReadDirEntries(m.parentDir, m.showHidden)
	if err != nil {
		return
	}

	m.parentItems = items
	for i, it := range m.parentItems {
		if it.Path == m.currDir {
			m.parentCursor = i
			break
		}
	}
}

func (m *Model) sortItems() {
	if len(m.items) == 0 {
		return
	}

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		switch m.sortMode {
		case SortSizeDesc:
			if a.Size != b.Size {
				return a.Size > b.Size
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)

		case SortSizeAsc:
			if a.Size != b.Size {
				return a.Size < b.Size
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)

		case SortNameAsc:
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)

		case SortNameDesc:
			return strings.ToLower(a.Name) > strings.ToLower(b.Name)

		case SortCountDesc:
			countA := a.FileCount + a.DirCount
			countB := b.FileCount + b.DirCount
			if countA != countB {
				return countA > countB
			}
			return a.Size > b.Size

		default:
			return a.Size > b.Size
		}
	})
}

func (m *Model) adjustScroll() {
	contentHeight := m.height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}

	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+contentHeight {
		m.scroll = m.cursor - contentHeight + 1
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
}

func (m *Model) updatePreview() {
	if len(m.items) == 0 || m.cursor < 0 || m.cursor >= len(m.items) {
		m.dirPreview = DirPreviewInfo{}
		m.filePreview = FilePreviewInfo{}
		return
	}

	sel := m.items[m.cursor]
	if sel.IsDir {
		m.dirPreview = GenerateDirPreview(sel.Path, 8, 8)
	} else {
		m.filePreview = GenerateFilePreview(sel.Path, 40)
	}
}

func (m *Model) totalActiveSize() int64 {
	var total int64
	for _, it := range m.items {
		total += it.Size
	}
	return total
}

func (m *Model) maxActiveSize() int64 {
	if len(m.items) == 0 {
		return 1
	}
	max := int64(0)
	for _, it := range m.items {
		if it.Size > max {
			max = it.Size
		}
	}
	if max <= 0 {
		return 1
	}
	return max
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adjustScroll()
		return m, nil

	case SpinnerTickMsg:
		if !m.scanning {
			return m, nil
		}
		m.spinnerIdx = (m.spinnerIdx + 1) % len(spinnerFrames)
		return m, tickSpinner()

	case ScanProgressMsg:
		if msg.Path != m.currDir {
			return m, nil
		}

		var selPath string
		if m.cursor >= 0 && m.cursor < len(m.items) {
			selPath = m.items[m.cursor].Path
		}

		for i := range m.items {
			if m.items[i].Path == msg.ItemPath {
				m.items[i].Size = msg.Size
				m.items[i].FileCount = msg.FileCount
				m.items[i].DirCount = msg.DirCount
				m.items[i].SizeKnown = true
				m.items[i].ScanErr = msg.Err
				m.scanDone++
				break
			}
		}

		if m.scanDone >= m.scanTotal {
			m.scanning = false
		}

		if m.sortMode == SortSizeDesc || m.sortMode == SortSizeAsc || m.sortMode == SortCountDesc {
			m.sortItems()
			if selPath != "" {
				for i, it := range m.items {
					if it.Path == selPath {
						m.cursor = i
						break
					}
				}
			}
		}

		m.adjustScroll()
		m.updatePreview()
		return m, nil

	case tea.KeyMsg:
		if m.showAbout {
			m.showAbout = false
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			if m.scanCancel != nil {
				m.scanCancel()
			}
			return m, tea.Quit

		case "tab", "t":
			m.lang = m.lang.Next()
			return m, nil

		case "?", "a":
			m.showAbout = true
			return m, nil

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
				m.updatePreview()
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.adjustScroll()
				m.updatePreview()
			}
			return m, nil

		case "home", "g":
			m.cursor = 0
			m.adjustScroll()
			m.updatePreview()
			return m, nil

		case "end", "G":
			if len(m.items) > 0 {
				m.cursor = len(m.items) - 1
				m.adjustScroll()
				m.updatePreview()
			}
			return m, nil

		case "pageup":
			contentHeight := m.height - 5
			if contentHeight < 1 {
				contentHeight = 1
			}
			m.cursor -= contentHeight
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.adjustScroll()
			m.updatePreview()
			return m, nil

		case "pagedown":
			contentHeight := m.height - 5
			if contentHeight < 1 {
				contentHeight = 1
			}
			m.cursor += contentHeight
			if m.cursor >= len(m.items) {
				m.cursor = len(m.items) - 1
			}
			m.adjustScroll()
			m.updatePreview()
			return m, nil

		case "enter", "right", "l":
			if len(m.items) > 0 && m.cursor >= 0 && m.cursor < len(m.items) {
				selected := m.items[m.cursor]
				if selected.IsDir {
					cmd := m.loadDirectoryCmd(selected.Path)
					return m, cmd
				}
			}
			return m, nil

		case "backspace", "left", "h":
			if m.currDir != "/" {
				cmd := m.loadDirectoryCmd(filepath.Dir(m.currDir))
				return m, cmd
			}
			return m, nil

		case "s":
			m.sortMode = (m.sortMode + 1) % 5
			var selPath string
			if m.cursor >= 0 && m.cursor < len(m.items) {
				selPath = m.items[m.cursor].Path
			}
			m.sortItems()
			if selPath != "" {
				for i, it := range m.items {
					if it.Path == selPath {
						m.cursor = i
						break
					}
				}
			}
			m.adjustScroll()
			m.updatePreview()
			return m, nil

		case ".", "H":
			m.showHidden = !m.showHidden
			cmd := m.loadDirectoryCmd(m.currDir)
			return m, cmd

		case "r":
			InvalidateCache(m.currDir)
			cmd := m.loadDirectoryCmd(m.currDir)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.showAbout {
		return m.renderAboutModal()
	}

	d := T(m.lang)
	totalSize := m.totalActiveSize()
	maxSize := m.maxActiveSize()

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	headerPathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#bbbbbb"))
	statusScanningStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#999999"))
	statusReadyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	sortBadgeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd"))
	hiddenBadgeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#383838"))

	var statusText string
	if m.scanning {
		spin := spinnerFrames[m.spinnerIdx]
		statusText = statusScanningStyle.Render(fmt.Sprintf("%s "+d.Scanning, spin, m.scanDone, m.scanTotal))
	} else {
		statusText = statusReadyStyle.Render(fmt.Sprintf(d.ScanDone, len(m.items), FormatBytes(totalSize)))
	}

	var sortStr string
	switch m.sortMode {
	case SortSizeDesc:
		sortStr = d.SortSizeDesc
	case SortSizeAsc:
		sortStr = d.SortSizeAsc
	case SortNameAsc:
		sortStr = d.SortNameAsc
	case SortNameDesc:
		sortStr = d.SortNameDesc
	case SortCountDesc:
		sortStr = d.SortCountDesc
	}

	hiddenStr := d.HiddenShown
	if !m.showHidden {
		hiddenStr = d.HiddenHidden
	}

	headerLine1 := titleStyle.Render(d.AppTitle) + "  " +
		headerPathStyle.Render(TruncateString(m.currDir, m.width-40)) + "  " +
		statusText

	headerLine2 := fmt.Sprintf("%s: %s  |  %s  |  Tab: %s",
		d.SortLabel,
		sortBadgeStyle.Render(sortStr),
		hiddenBadgeStyle.Render(hiddenStr),
		strings.ToUpper(string(m.lang)),
	)

	contentHeight := m.height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}

	var leftW, centerW, rightW int
	var showLeft, showRight bool

	if m.width >= 110 {
		leftW = m.width * 18 / 100
		if leftW < 16 {
			leftW = 16
		}
		if leftW > 24 {
			leftW = 24
		}
		centerW = m.width * 45 / 100
		if centerW < 36 {
			centerW = 36
		}
		rightW = m.width - leftW - centerW - 2
		showLeft = true
		showRight = true
	} else if m.width >= 72 {
		leftW = 0
		centerW = m.width * 50 / 100
		if centerW < 34 {
			centerW = 34
		}
		rightW = m.width - centerW - 1
		showLeft = false
		showRight = true
	} else {
		leftW = 0
		centerW = m.width
		rightW = 0
		showLeft = false
		showRight = false
	}

	leftLines := make([]string, contentHeight)
	centerLines := make([]string, contentHeight)
	rightLines := make([]string, contentHeight)

	if showLeft {
		parentTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#666666")).
			Render(TruncateString("< "+filepath.Base(m.parentDir), leftW-1))
		leftLines[0] = FitWidth(parentTitle, leftW)

		startRow := 0
		if m.parentCursor >= contentHeight-1 {
			startRow = m.parentCursor - (contentHeight - 2)
		}

		for r := 1; r < contentHeight; r++ {
			pIdx := startRow + (r - 1)
			if pIdx < len(m.parentItems) {
				pItem := m.parentItems[pIdx]
				isCurrent := (pIdx == m.parentCursor)
				pName := pItem.Name
				if pItem.IsDir {
					pName += "/"
				}
				if isCurrent {
					line := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).Render("> " + pName)
					leftLines[r] = FitWidth(line, leftW)
				} else {
					line := lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render("  " + pName)
					leftLines[r] = FitWidth(line, leftW)
				}
			} else {
				leftLines[r] = FitWidth("", leftW)
			}
		}
	}

	selectedBgStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2a2a2a")).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	for r := 0; r < contentHeight; r++ {
		idx := m.scroll + r
		if idx >= len(m.items) {
			centerLines[r] = FitWidth("", centerW)
			continue
		}

		item := m.items[idx]
		isSelected := (idx == m.cursor)

		var prefix string
		if isSelected {
			prefix = "> "
		} else {
			prefix = "  "
		}

		var sizeFormatted string
		if item.IsDir && !item.SizeKnown {
			sizeFormatted = "      ..."
		} else {
			sizeFormatted = fmt.Sprintf("%9s", FormatBytes(item.Size))
		}
		sizeStyled := lipgloss.NewStyle().Foreground(GetSizeColor(item.Size)).Render(sizeFormatted)

		var barStyled string
		if centerW >= 42 {
			ratio := 0.0
			if maxSize > 0 && item.Size > 0 {
				ratio = float64(item.Size) / float64(maxSize)
			}
			barStr := MakeProgressBar(ratio, 8)
			barStyled = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(barStr)
		}

		displayName := item.Name
		if item.IsDir {
			displayName += "/"
		}

		var nameStyle lipgloss.Style
		if item.IsDir {
			if strings.HasPrefix(item.Name, ".") {
				nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
			} else {
				nameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
			}
		} else {
			nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
		}
		nameStyled := nameStyle.Render(displayName)

		var rowContent string
		if barStyled != "" {
			rowContent = prefix + sizeStyled + "  " + barStyled + "  " + nameStyled
		} else {
			rowContent = prefix + sizeStyled + "  " + nameStyled
		}

		fittedRow := FitWidth(rowContent, centerW)
		if isSelected {
			centerLines[r] = selectedBgStyle.Render(fittedRow)
		} else {
			centerLines[r] = fittedRow
		}
	}

	if showRight {
		var previewRaw []string
		if len(m.items) > 0 && m.cursor >= 0 && m.cursor < len(m.items) {
			sel := m.items[m.cursor]
			if sel.IsDir {
				previewRaw = GetDirPreviewLines(m.dirPreview, totalSize, rightW, contentHeight, m.lang)
			} else {
				previewRaw = GetFilePreviewLines(m.filePreview, rightW, contentHeight, m.lang)
			}
		} else {
			previewRaw = []string{lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(d.EmptyDir)}
		}

		for r := 0; r < contentHeight; r++ {
			if r < len(previewRaw) {
				rightLines[r] = FitWidth(previewRaw[r], rightW)
			} else {
				rightLines[r] = FitWidth("", rightW)
			}
		}
	}

	sep := borderStyle.Render("|")
	var finalLines []string

	finalLines = append(finalLines, FitWidth(headerLine1, m.width))
	finalLines = append(finalLines, FitWidth(headerLine2, m.width))
	finalLines = append(finalLines, borderStyle.Render(strings.Repeat("-", m.width)))

	for r := 0; r < contentHeight; r++ {
		var rowStr string
		if showLeft && showRight {
			rowStr = leftLines[r] + sep + centerLines[r] + sep + rightLines[r]
		} else if showRight {
			rowStr = centerLines[r] + sep + rightLines[r]
		} else {
			rowStr = centerLines[r]
		}
		finalLines = append(finalLines, rowStr)
	}

	finalLines = append(finalLines, borderStyle.Render(strings.Repeat("-", m.width)))
	finalLines = append(finalLines, FitWidth(lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(d.HelpFooter), m.width))

	return strings.Join(finalLines, "\n")
}
