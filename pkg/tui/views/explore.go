package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wordy-tui/wordy/pkg/seed"
	"github.com/wordy-tui/wordy/pkg/tui/styles"
)

type RelatedWordItem struct {
	RelType string
	Word    string
}

// RenderExploreView builds the word detail and related word link explorer layout.
func RenderExploreView(
	details seed.WordDetails,
	relatedItems []RelatedWordItem,
	focusedRelIndex int,
	isBookmarked bool,
	width int,
	height int,
) string {
	if details.Word == "" {
		return styles.CardBoxStyle.Width(width - 4).Render(
			lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("Select a word from the list to explore definitions, frequency stats, and related words."),
		)
	}

	var sb strings.Builder

	// Header line: Word, Bookmark, Phonetic, Pos Badge, Rarity Badge
	bookmarkStar := "☆"
	if isBookmarked {
		bookmarkStar = lipgloss.NewStyle().Foreground(styles.ColorAmber).Render("★")
	}

	title := styles.WordTitleStyle.Render(details.Word) + " " + bookmarkStar
	phonetic := styles.PhoneticStyle.Render(details.Phonetic)
	posBadge := ""
	if details.PartOfSpeech != "" {
		posBadge = styles.BadgePosStyle.Render(details.PartOfSpeech)
	}
	rarityBadge := styles.GetRarityBadge(details.RarityTier)

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", phonetic, "  ", posBadge, "  ", rarityBadge))
	sb.WriteString("\n\n")

	// Frequency meter bar
	freqBar := styles.RenderFrequencyBar(details.CorpusCount, details.ZipfScore, 20)
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Corpus Usage: ") + freqBar + "\n\n")

	// Definitions Section
	sb.WriteString(styles.SectionHeaderStyle.Render("📖 Definitions"))
	sb.WriteString("\n")
	if len(details.Definitions) == 0 {
		sb.WriteString(styles.DefinitionTextStyle.Render("No definition available."))
		sb.WriteString("\n")
	} else {
		for i, def := range details.Definitions {
			defNum := lipgloss.NewStyle().Foreground(styles.ColorPurple).Bold(true).Render(fmt.Sprintf("%d. ", i+1))
			sb.WriteString(styles.DefinitionTextStyle.Render(defNum + def.Text))
			sb.WriteString("\n")
		}
	}

	// Examples Section
	if len(details.Examples) > 0 {
		sb.WriteString("\n")
		sb.WriteString(styles.SectionHeaderStyle.Render("💬 Usage Examples"))
		sb.WriteString("\n")
		for _, ex := range details.Examples {
			sb.WriteString(styles.ExampleTextStyle.Render(fmt.Sprintf("“%s”", ex.Text)))
			if ex.Author != "" || ex.Title != "" {
				source := strings.Trim(fmt.Sprintf("— %s, %s", ex.Author, ex.Title), ", ")
				sb.WriteString("\n" + styles.ExampleTextStyle.Render(lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(source)))
			}
			sb.WriteString("\n")
		}
	}

	// Related Words Interactive Explorer Section
	sb.WriteString("\n")
	sb.WriteString(styles.SectionHeaderStyle.Render("🔗 Related Words (Press Enter to Jump)"))
	sb.WriteString("\n")

	if len(relatedItems) == 0 {
		sb.WriteString(styles.DefinitionTextStyle.Render(lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("No related words linked for this term.")))
		sb.WriteString("\n")
	} else {
		var pills []string
		for i, item := range relatedItems {
			relPrefix := lipgloss.NewStyle().Foreground(styles.ColorAmber).Render("[" + item.RelType + "] ")
			pillText := item.Word

			if i == focusedRelIndex {
				pills = append(pills, styles.WordPillSelected.Render("➜ "+relPrefix+pillText))
			} else {
				pills = append(pills, styles.WordPillUnselected.Render(relPrefix+pillText))
			}
		}

		// Wrap pills nicely
		line := ""
		for _, p := range pills {
			if len(line)+lipgloss.Width(p) > width-10 {
				sb.WriteString("  " + line + "\n")
				line = p + " "
			} else {
				line += p + " "
			}
		}
		if line != "" {
			sb.WriteString("  " + line + "\n")
		}
	}

	// Source Attribution
	if details.AttributionText != "" {
		sb.WriteString("\n")
		sb.WriteString(styles.AttributionStyle.Render("Attribution: " + details.AttributionText))
	}

	cardWidth := width - 4
	if cardWidth < 40 {
		cardWidth = 40
	}

	return styles.CardBoxStyle.Width(cardWidth).Render(sb.String())
}
