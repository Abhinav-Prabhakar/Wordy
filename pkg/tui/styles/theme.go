package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	ColorPurple      = lipgloss.Color("#A78BFA") // Soft Lavender
	ColorDeepPurple  = lipgloss.Color("#8B5CF6") // Vibrant Purple
	ColorPink        = lipgloss.Color("#F472B6") // Hot Pink / Boba Strawberry
	ColorMatcha      = lipgloss.Color("#34D399") // Matcha Green
	ColorCyan        = lipgloss.Color("#38BDF8") // Bright Cyan
	ColorAmber       = lipgloss.Color("#FBBF24") // Warm Amber
	ColorCoral       = lipgloss.Color("#F87171") // Coral Red
	ColorBgDark      = lipgloss.Color("#0F172A") // Slate Dark
	ColorCardBg      = lipgloss.Color("#1E293B") // Slate Medium Card
	ColorMuted       = lipgloss.Color("#64748B") // Muted Slate Gray
	ColorSubtle      = lipgloss.Color("#94A3B8") // Subtle Gray
	ColorWhite       = lipgloss.Color("#F8FAFC") // Off White

	// Base Layout Styles
	AppTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorDeepPurple).
			Padding(0, 1).
			MarginRight(1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPink).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorPink).
			Padding(0, 1)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Padding(0, 1)

	StatusPillOnline = lipgloss.NewStyle().
				Foreground(ColorMatcha).
				Background(lipgloss.Color("#064E3B")).
				Padding(0, 1).
				Bold(true)

	StatusPillSeed = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Background(lipgloss.Color("#0C4A6E")).
				Padding(0, 1).
				Bold(true)

	StatusPillRateLimit = lipgloss.NewStyle().
				Foreground(ColorAmber).
				Background(lipgloss.Color("#78350F")).
				Padding(0, 1).
				Bold(true)

	// Card Styles
	CardBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Background(ColorCardBg).
			Padding(1, 2)

	CardBoxActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPink).
				Background(ColorCardBg).
				Padding(1, 2)

	WordTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPink)

	PhoneticStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Italic(true)

	BadgePosStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorDeepPurple).
			Padding(0, 1).
			MarginLeft(1)

	BadgeRarityRare = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorCoral).
			Padding(0, 1).
			Bold(true)

	BadgeRarityObscure = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(ColorAmber).
				Padding(0, 1).
				Bold(true)

	BadgeRarityElegant = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(ColorPurple).
				Padding(0, 1).
				Bold(true)

	BadgeRarityUncommon = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(ColorMatcha).
				Padding(0, 1).
				Bold(true)

	// Definition & Section Styles
	SectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPurple).
				MarginTop(1)

	DefinitionTextStyle = lipgloss.NewStyle().
				Foreground(ColorWhite).
				MarginLeft(2)

	ExampleTextStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Italic(true).
				MarginLeft(4)

	AttributionStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true).
				MarginTop(1)

	// Related Words Selector Styles
	RelationTypeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAmber).
				MarginRight(1)

	WordPillSelected = lipgloss.NewStyle().
				Foreground(ColorBgDark).
				Background(ColorCyan).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	WordPillUnselected = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorCyan).
				Padding(0, 1).
				MarginRight(1)

	// Flashcard Ratings Buttons
	BtnAgainStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorCoral).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	BtnHardStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorAmber).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	BtnGoodStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorPurple).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	BtnEasyStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorMatcha).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	HelpFooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorPurple).
			Bold(true)
)

// RenderFrequencyBar returns a visual progress meter for word usage/corpus frequency.
func RenderFrequencyBar(corpusCount int64, zipf float64, width int) string {
	if width < 10 {
		width = 10
	}
	// Scale Zipf (1.0 to 5.0) to bar length
	ratio := (zipf - 1.0) / 4.0
	if ratio < 0.05 {
		ratio = 0.05
	}
	if ratio > 1.0 {
		ratio = 1.0
	}

	filledLen := int(float64(width) * ratio)
	emptyLen := width - filledLen

	filled := strings.Repeat("█", filledLen)
	empty := strings.Repeat("░", emptyLen)

	barStr := lipgloss.NewStyle().Foreground(ColorPink).Render(filled) +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(empty)

	return fmt.Sprintf("[%s] Zipf %.1f", barStr, zipf)
}

// GetRarityBadge returns a styled badge for rarity tier.
func GetRarityBadge(tier string) string {
	switch strings.Title(strings.ToLower(tier)) {
	case "Rare":
		return BadgeRarityRare.Render("✦ RARE")
	case "Obscure":
		return BadgeRarityObscure.Render("✧ OBSCURE")
	case "Elegant":
		return BadgeRarityElegant.Render("◇ ELEGANT")
	default:
		return BadgeRarityUncommon.Render("◈ UNCOMMON")
	}
}
