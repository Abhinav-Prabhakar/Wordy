package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Modern Vibrant Color Palette
	ColorPink        = lipgloss.Color("#EC4899") // Hot Pink
	ColorLightPink   = lipgloss.Color("#F472B6") // Soft Boba Pink
	ColorPurple      = lipgloss.Color("#A855F7") // Vibrant Purple
	ColorLavender    = lipgloss.Color("#C084FC") // Soft Lavender
	ColorDeepIndigo  = lipgloss.Color("#1E1B4B") // Deep Boba Background
	ColorCardBg      = lipgloss.Color("#0F172A") // Slate Dark Card
	ColorCardBorder  = lipgloss.Color("#334155") // Subtle Border
	ColorMatcha      = lipgloss.Color("#10B981") // Matcha Green
	ColorCyan        = lipgloss.Color("#06B6D4") // Bright Cyan
	ColorAmber       = lipgloss.Color("#F59E0B") // Warm Amber
	ColorCoral       = lipgloss.Color("#EF4444") // Coral Red
	ColorMuted       = lipgloss.Color("#64748B") // Muted Slate
	ColorSubtle      = lipgloss.Color("#94A3B8") // Subtle Gray
	ColorWhite       = lipgloss.Color("#F8FAFC") // Pure Soft White

	// App Header
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorPurple).
			Padding(0, 1).
			MarginRight(1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPink).
			Background(lipgloss.Color("#312E81")).
			Padding(0, 1).
			MarginRight(1)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Padding(0, 1).
				MarginRight(1)

	// Status Pills
	PillOnline = lipgloss.NewStyle().
			Foreground(ColorMatcha).
			Background(lipgloss.Color("#064E3B")).
			Padding(0, 1).
			Bold(true)

	PillSeed = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Background(lipgloss.Color("#0C4A6E")).
			Padding(0, 1).
			Bold(true)

	PillRateLimit = lipgloss.NewStyle().
			Foreground(ColorAmber).
			Background(lipgloss.Color("#78350F")).
			Padding(0, 1).
			Bold(true)

	// Card Containers
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Background(ColorCardBg).
			Padding(1, 2)

	CardBoxActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPink).
				Background(ColorCardBg).
				Padding(1, 2)

	// Text & Badges
	WordTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorLightPink)

	PhoneticStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Italic(true)

	BadgePosStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(lipgloss.Color("#4338CA")).
			Padding(0, 1).
			Bold(true)

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
				Background(ColorLavender).
				Padding(0, 1).
				Bold(true)

	BadgeRarityUncommon = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(ColorMatcha).
				Padding(0, 1).
				Bold(true)

	SectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorLavender).
				MarginTop(1)

	DefinitionTextStyle = lipgloss.NewStyle().
				Foreground(ColorWhite).
				MarginLeft(1)

	ExampleTextStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Italic(true).
				MarginLeft(2)

	AttributionStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true).
				MarginTop(1)

	// Related Words Links
	WordPillSelected = lipgloss.NewStyle().
				Foreground(ColorDeepIndigo).
				Background(ColorCyan).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	WordPillUnselected = lipgloss.NewStyle().
				Foreground(ColorCyan).
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#164E63")).
				Padding(0, 1).
				MarginRight(1)

	// Rating Buttons
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

	// Footer Help
	HelpBar = lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Background(lipgloss.Color("#1E1B4B")).
		Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorPink).
			Bold(true)
)

// RenderFrequencyBar returns a visual progress meter for word usage frequency.
func RenderFrequencyBar(corpusCount int64, zipf float64, width int) string {
	if width < 8 {
		width = 8
	}
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
