package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wordy-tui/wordy/pkg/seed"
	"github.com/wordy-tui/wordy/pkg/tui/styles"
)

// RenderNetworkView renders a visual graph cluster of word relations.
func RenderNetworkView(details seed.WordDetails, selectedRelIndex int, width int, height int) string {
	if details.Word == "" {
		return styles.CardBoxStyle.Width(width - 4).Render(
			lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("Select a word to visualize its semantic relation network graph."),
		)
	}

	var sb strings.Builder

	header := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorCyan).Render(fmt.Sprintf("🕸️  Semantic Relation Graph: %s", strings.ToUpper(details.Word)))
	sb.WriteString(header + "\n\n")

	// Central Node
	centerNode := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorBgDark).
		Background(styles.ColorPink).
		Padding(0, 2).
		Render(fmt.Sprintf("◉ %s", details.Word))

	sb.WriteString("              " + centerNode + "\n")
	sb.WriteString("                   │\n")
	sb.WriteString("      ┌────────────┼────────────┐\n")

	// Branches: Synonyms, Antonyms, Hypernyms/Related
	syns := details.RelatedWords["synonym"]
	ants := details.RelatedWords["antonym"]

	synBranchHeader := lipgloss.NewStyle().Foreground(styles.ColorMatcha).Bold(true).Render("▼ SYNONYMS")
	antBranchHeader := lipgloss.NewStyle().Foreground(styles.ColorCoral).Bold(true).Render("▼ ANTONYMS")
	relBranchHeader := lipgloss.NewStyle().Foreground(styles.ColorCyan).Bold(true).Render("▼ CONTEXTUAL")

	sb.WriteString(fmt.Sprintf("  %-20s %-20s %-20s\n", synBranchHeader, antBranchHeader, relBranchHeader))

	maxRows := 5
	for i := 0; i < maxRows; i++ {
		synItem := " "
		if i < len(syns) {
			synItem = styles.WordPillUnselected.Render("• " + syns[i])
		}
		antItem := " "
		if i < len(ants) {
			antItem = lipgloss.NewStyle().Foreground(styles.ColorCoral).Render("• " + ants[i])
		}
		relItem := " "
		// Fallback to other relation types
		for k, v := range details.RelatedWords {
			if k != "synonym" && k != "antonym" && i < len(v) {
				relItem = styles.WordPillUnselected.Render("[" + k + "] " + v[i])
				break
			}
		}

		sb.WriteString(fmt.Sprintf("  %-22s %-22s %-22s\n", synItem, antItem, relItem))
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Tip: Press Tab to switch to Explore view and click into any related word."))

	cardWidth := width - 4
	if cardWidth < 40 {
		cardWidth = 40
	}

	return styles.CardBoxStyle.Width(cardWidth).Render(sb.String())
}
