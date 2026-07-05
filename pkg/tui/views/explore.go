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

// RenderExploreContent generates the inner text content for the Explore viewport.
func RenderExploreContent(
	details seed.WordDetails,
	relatedItems []RelatedWordItem,
	focusedRelIndex int,
	isBookmarked bool,
	width int,
) string {
	if details.Word == "" {
		return lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("Select a word from the left list to inspect definitions, frequency stats, and related words.")
	}

	var sb strings.Builder

	// Header line: Word Title, Bookmark, Phonetic, POS, Rarity
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

	// Frequency Meter
	freqBar := styles.RenderFrequencyBar(details.CorpusCount, details.ZipfScore, 16)
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Corpus Frequency: ") + freqBar + "\n\n")

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

	// Usage Examples Section
	if len(details.Examples) > 0 {
		sb.WriteString("\n")
		sb.WriteString(styles.SectionHeaderStyle.Render("💬 Examples"))
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

	// Related Words Cluster Section
	sb.WriteString("\n")
	sb.WriteString(styles.SectionHeaderStyle.Render("🔗 Related Words (Press Enter to Jump)"))
	sb.WriteString("\n")

	if len(relatedItems) == 0 {
		sb.WriteString(styles.DefinitionTextStyle.Render(lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("No related terms linked for this word.")))
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

		line := ""
		for _, p := range pills {
			if len(line)+lipgloss.Width(p) > width-8 {
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
		sb.WriteString(styles.AttributionStyle.Render("Source: " + details.AttributionText))
	}

	return sb.String()
}
