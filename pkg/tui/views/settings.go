package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/wordy-tui/wordy/pkg/api"
	"github.com/wordy-tui/wordy/pkg/storage"
	"github.com/wordy-tui/wordy/pkg/tui/styles"
)

// RenderSettingsView renders the settings menu for API Key & Rarity preferences.
func RenderSettingsView(
	cfg storage.Config,
	apiKeyInput textinput.Model,
	rateLimit api.RateLimitInfo,
	cachedCount int,
	statusMsg string,
	width int,
	height int,
) string {
	var sb strings.Builder

	header := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPurple).Render("⚙️  Wordy Configuration & API Settings")
	sb.WriteString(header + "\n\n")

	// API Key Section
	sb.WriteString(styles.SectionHeaderStyle.Render("🔑 Wordnik API Key"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Leave empty to use built-in offline seed dictionary. Sign up at https://developer.wordnik.com/"))
	sb.WriteString("\n\n")

	sb.WriteString("API Key: " + apiKeyInput.View() + "\n\n")

	// Rate Limits Section
	sb.WriteString(styles.SectionHeaderStyle.Render("⚡ API Rate-Limit Status"))
	sb.WriteString("\n")
	rateStr := fmt.Sprintf(
		"Minute: %d / %d calls remaining  │  Hour: %d / %d calls remaining",
		rateLimit.RemainingMinute, rateLimit.LimitMinute,
		rateLimit.RemainingHour, rateLimit.LimitHour,
	)
	if rateLimit.IsRateLimited {
		rateStr += lipgloss.NewStyle().Foreground(styles.ColorCoral).Bold(true).Render(" [429 RATE LIMITED - BACKOFF ACTIVE]")
	}
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorCyan).Render(rateStr) + "\n\n")

	// Storage & Cache Section
	sb.WriteString(styles.SectionHeaderStyle.Render("💾 Local Storage & Cache"))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Cached Word Definitions: %d terms stored in ~/.cache/wordy/\n", cachedCount))
	sb.WriteString("User Progress DB Path: " + storage.GetDataPath() + "\n\n")

	// Controls legend
	sb.WriteString(styles.SectionHeaderStyle.Render("⌨️ Actions"))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorPink).Render("[ENTER] Save API Key & Settings   │   [c] Clear Local Cache"))
	sb.WriteString("\n\n")

	if statusMsg != "" {
		sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(styles.ColorMatcha).Render("✓ "+statusMsg) + "\n")
	}

	cardWidth := width - 4
	if cardWidth < 40 {
		cardWidth = 40
	}

	return styles.CardBoxStyle.Width(cardWidth).Render(sb.String())
}
