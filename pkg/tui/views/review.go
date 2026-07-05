package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wordy-tui/wordy/pkg/seed"
	"github.com/wordy-tui/wordy/pkg/srs"
	"github.com/wordy-tui/wordy/pkg/tui/styles"
)

// RenderReviewView renders the SRS flashcard review screen.
func RenderReviewView(
	details seed.WordDetails,
	card srs.CardState,
	isFlipped bool,
	dueCount int,
	totalMastered int,
	streakDays int,
	width int,
	height int,
) string {
	var mainContent strings.Builder

	// Top stats bar
	duePill := lipgloss.NewStyle().Foreground(styles.ColorWhite).Background(styles.ColorDeepPurple).Padding(0, 1).Bold(true).Render(fmt.Sprintf("Due Today: %d", dueCount))
	masteredPill := lipgloss.NewStyle().Foreground(styles.ColorWhite).Background(styles.ColorMatcha).Padding(0, 1).Bold(true).Render(fmt.Sprintf("Mastered: %d", totalMastered))
	streakPill := lipgloss.NewStyle().Foreground(styles.ColorWhite).Background(styles.ColorPink).Padding(0, 1).Bold(true).Render(fmt.Sprintf("Streak: 🔥 %d days", streakDays))

	statsBar := lipgloss.JoinHorizontal(lipgloss.Center, duePill, "   ", masteredPill, "   ", streakPill)
	mainContent.WriteString(statsBar + "\n\n")

	if dueCount == 0 && details.Word == "" {
		// All cards reviewed for today
		emptyCard := styles.CardBoxStyle.Width(width - 6).Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Bold(true).Foreground(styles.ColorMatcha).Render("🎉 All Reviews Completed!"),
				"\n",
				lipgloss.NewStyle().Foreground(styles.ColorWhite).Render("You have mastered all scheduled vocabulary for today."),
				lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Check back tomorrow or explore new words in the Explore tab!"),
			),
		)
		mainContent.WriteString(emptyCard)
		return mainContent.String()
	}

	var cardBody strings.Builder

	// Header line
	title := styles.WordTitleStyle.Render(details.Word)
	phonetic := styles.PhoneticStyle.Render(details.Phonetic)
	pos := ""
	if details.PartOfSpeech != "" {
		pos = styles.BadgePosStyle.Render(details.PartOfSpeech)
	}
	rarity := styles.GetRarityBadge(details.RarityTier)

	cardBody.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", phonetic, "  ", pos, "  ", rarity))
	cardBody.WriteString("\n\n")

	// Frequency Bar
	cardBody.WriteString(styles.RenderFrequencyBar(details.CorpusCount, details.ZipfScore, 18) + "\n\n")

	if !isFlipped {
		// FRONT OF FLASHCARD
		cardBody.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render(fmt.Sprintf("SRS Status: Reps %d • Interval %dd • Ease %.2f", card.Repetitions, card.IntervalDays, card.EaseFactor)))
		cardBody.WriteString("\n\n")

		flipPrompt := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorPink).
			Border(lipgloss.NormalBorder()).
			BorderForeground(styles.ColorPink).
			Padding(0, 2).
			Render("Press [ SPACE ] to Flip Card")

		cardBody.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Render(flipPrompt))
	} else {
		// BACK OF FLASHCARD
		cardBody.WriteString(styles.SectionHeaderStyle.Render("📖 Definition"))
		cardBody.WriteString("\n")
		if len(details.Definitions) > 0 {
			cardBody.WriteString(styles.DefinitionTextStyle.Render(details.Definitions[0].Text))
		} else {
			cardBody.WriteString(styles.DefinitionTextStyle.Render("No definition available."))
		}
		cardBody.WriteString("\n\n")

		if len(details.Examples) > 0 {
			cardBody.WriteString(styles.SectionHeaderStyle.Render("💬 Example"))
			cardBody.WriteString("\n")
			cardBody.WriteString(styles.ExampleTextStyle.Render(fmt.Sprintf("“%s”", details.Examples[0].Text)))
			cardBody.WriteString("\n\n")
		}

		// Recall Rating Buttons (1-4)
		cardBody.WriteString(styles.SectionHeaderStyle.Render("Rate Your Recall Quality:"))
		cardBody.WriteString("\n\n")

		btn1 := styles.BtnAgainStyle.Render("[1] Again (Reset)")
		btn2 := styles.BtnHardStyle.Render("[2] Hard (1d)")
		btn3 := styles.BtnGoodStyle.Render("[3] Good (6d)")
		btn4 := styles.BtnEasyStyle.Render("[4] Easy (14d)")

		ratingsBar := lipgloss.JoinHorizontal(lipgloss.Center, btn1, btn2, btn3, btn4)
		cardBody.WriteString(ratingsBar)
	}

	cardWidth := width - 6
	if cardWidth < 40 {
		cardWidth = 40
	}

	mainContent.WriteString(styles.CardBoxActiveStyle.Width(cardWidth).Render(cardBody.String()))
	return mainContent.String()
}
