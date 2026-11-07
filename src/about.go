package src

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderAboutModal() string {
	d := T(m.lang)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff"))

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#aaaaaa")).
		Width(15)

	valStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	cardWidth := 46
	if m.width-6 < cardWidth && m.width > 10 {
		cardWidth = m.width - 6
	}

	var sb strings.Builder

	header := titleStyle.Render(d.AboutTitle) + "  " + valStyle.Render("v"+Version)
	sb.WriteString(lipgloss.NewStyle().Width(cardWidth).Align(lipgloss.Center).Render(header) + "\n")
	sb.WriteString(lipgloss.NewStyle().Width(cardWidth).Align(lipgloss.Center).Render(subtitleStyle.Render(d.AboutSubtitle)) + "\n\n")

	divider := dividerStyle.Render(strings.Repeat("-", cardWidth))
	sb.WriteString(divider + "\n")

	sb.WriteString(labelStyle.Render(d.AboutVersion) + valStyle.Render(Version) + "\n")
	sb.WriteString(labelStyle.Render(d.AboutCreator) + valStyle.Render("Ars-byte") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutRepo) + valStyle.Render("github.com/Ars-byte/bigcrock") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutLicense) + valStyle.Render(d.AboutLicenseVal) + "\n")
	sb.WriteString(labelStyle.Render(d.AboutYear) + valStyle.Render("2026") + "\n")
	sb.WriteString(labelStyle.Render(d.AboutBuiltWith) + valStyle.Render("Go") + "\n")

	sb.WriteString(divider + "\n\n")
	sb.WriteString(lipgloss.NewStyle().Width(cardWidth).Align(lipgloss.Center).Render(footerStyle.Render(d.AboutReturnHint)))

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1, 2)

	card := boxStyle.Render(sb.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}
