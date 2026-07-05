package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wordy-tui/wordy/pkg/api"
	"github.com/wordy-tui/wordy/pkg/seed"
	"github.com/wordy-tui/wordy/pkg/srs"
	"github.com/wordy-tui/wordy/pkg/storage"
	"github.com/wordy-tui/wordy/pkg/tui/styles"
	"github.com/wordy-tui/wordy/pkg/tui/views"
)

type wordListItem struct {
	word       string
	pos        string
	rarityTier string
}

func (i wordListItem) Title() string       { return i.word }
func (i wordListItem) Description() string { return fmt.Sprintf("%s • %s", i.pos, i.rarityTier) }
func (i wordListItem) FilterValue() string { return i.word }

type wordFetchedMsg struct {
	details seed.WordDetails
	err     error
}

type Model struct {
	cfg                 storage.Config
	store               *storage.Store
	apiClient           *api.Client
	activeTab           int // 0: Explore, 1: Review, 2: Settings
	width               int
	height              int
	wordList            list.Model
	exploreViewport     viewport.Model
	apiKeyInput         textinput.Model
	selectedWord        string
	selectedWordDetails seed.WordDetails
	wordHistory         []string
	relatedItems        []views.RelatedWordItem
	focusedRelIndex     int
	// SRS State
	srsDueWords      []string
	srsCurrentIdx    int
	srsCurrentCard   srs.CardState
	isCardFlipped    bool
	srsTotalMastered int
	// UI Feedback
	statusMsg string
	toastMsg  string
}

func NewModel() (*Model, error) {
	cfg := storage.LoadConfig()
	store, err := storage.NewStore()
	if err != nil {
		return nil, err
	}
	apiClient := api.NewClient(cfg.WordnikAPIKey)

	// Prepare initial seed list items
	seedWords := seed.GetSeedWords()
	items := make([]list.Item, len(seedWords))
	for i, w := range seedWords {
		items[i] = wordListItem{
			word:       w.Word,
			pos:        w.PartOfSpeech,
			rarityTier: w.RarityTier,
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 28, 20)
	l.Title = "🧋 Vocabulary"
	l.SetShowHelp(false)

	vp := viewport.New(60, 20)

	ti := textinput.New()
	ti.Placeholder = "Enter Wordnik API Key..."
	ti.SetValue(cfg.WordnikAPIKey)
	ti.CharLimit = 128
	ti.Width = 36

	m := &Model{
		cfg:             cfg,
		store:           store,
		apiClient:       apiClient,
		activeTab:       0,
		wordList:        l,
		exploreViewport: vp,
		apiKeyInput:     ti,
		wordHistory:     make([]string, 0),
	}

	m.refreshSRSQueue()
	if len(seedWords) > 0 {
		m.loadWord(seedWords[0].Word)
	}

	return m, nil
}

func (m *Model) refreshSRSQueue() {
	now := time.Now()
	due := make([]string, 0)
	mastered := 0

	for _, w := range seed.GetSeedWords() {
		card := m.store.GetCard(w.Word)
		if card.Repetitions >= 3 {
			mastered++
		}
		if card.IsDue(now) {
			due = append(due, w.Word)
		}
	}

	m.srsDueWords = due
	m.srsTotalMastered = mastered
	m.srsCurrentIdx = 0
	m.isCardFlipped = false

	if len(due) > 0 {
		m.loadReviewCard(due[0])
	}
}

func (m *Model) loadReviewCard(word string) {
	card := m.store.GetCard(word)
	m.srsCurrentCard = card
	details, err := m.apiClient.FetchWordDetails(word)
	if err == nil {
		m.selectedWordDetails = details
	}
}

func (m *Model) loadWord(word string) tea.Cmd {
	wordLower := strings.ToLower(strings.TrimSpace(word))
	if wordLower == "" {
		return nil
	}

	m.selectedWord = wordLower

	if len(m.wordHistory) == 0 || m.wordHistory[len(m.wordHistory)-1] != wordLower {
		m.wordHistory = append(m.wordHistory, wordLower)
	}

	return func() tea.Msg {
		details, err := m.apiClient.FetchWordDetails(wordLower)
		return wordFetchedMsg{details: details, err: err}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updateViewportContent() {
	isBookmarked := m.store.IsBookmarked(m.selectedWord)
	content := views.RenderExploreContent(
		m.selectedWordDetails,
		m.relatedItems,
		m.focusedRelIndex,
		isBookmarked,
		m.exploreViewport.Width,
	)
	m.exploreViewport.SetContent(content)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentHeight := msg.Height - 6
		if contentHeight < 10 {
			contentHeight = 10
		}

		sidebarWidth := 30
		m.wordList.SetSize(sidebarWidth, contentHeight)

		vpWidth := msg.Width - sidebarWidth - 6
		if vpWidth < 30 {
			vpWidth = 30
		}
		m.exploreViewport.Width = vpWidth
		m.exploreViewport.Height = contentHeight - 2
		m.updateViewportContent()

	case wordFetchedMsg:
		if msg.err == nil {
			m.selectedWordDetails = msg.details
			items := make([]views.RelatedWordItem, 0)
			for relType, words := range msg.details.RelatedWords {
				for _, w := range words {
					items = append(items, views.RelatedWordItem{
						RelType: relType,
						Word:    w,
					})
				}
			}
			m.relatedItems = items
			m.focusedRelIndex = 0
			m.updateViewportContent()
		} else {
			m.toastMsg = "API Error: " + msg.err.Error()
		}

	case tea.KeyMsg:
		// Global Navigation Keys
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.apiKeyInput.Focused() && !m.wordList.SettingFilter() {
				return m, tea.Quit
			}

		case "tab":
			if !m.apiKeyInput.Focused() {
				m.activeTab = (m.activeTab + 1) % 3
				return m, nil
			}

		case "shift+tab":
			if !m.apiKeyInput.Focused() {
				m.activeTab = (m.activeTab + 2) % 3
				return m, nil
			}

		case "1", "2", "3":
			if m.activeTab == 1 && m.isCardFlipped && len(m.srsDueWords) > 0 {
				rating := srs.Rating(msg.String()[0] - '0')
				m.submitSRSRating(rating)
				return m, nil
			} else if !m.apiKeyInput.Focused() && !m.wordList.SettingFilter() {
				m.activeTab = int(msg.String()[0] - '1')
				return m, nil
			}

		case "4":
			if m.activeTab == 1 && m.isCardFlipped && len(m.srsDueWords) > 0 {
				m.submitSRSRating(srs.RatingEasy)
				return m, nil
			}

		case "backspace":
			if m.activeTab == 0 && len(m.wordHistory) > 1 && !m.wordList.SettingFilter() {
				m.wordHistory = m.wordHistory[:len(m.wordHistory)-1]
				prevWord := m.wordHistory[len(m.wordHistory)-1]
				return m, m.loadWord(prevWord)
			}
		}

		// Mode Specific Keys
		switch m.activeTab {
		case 0: // Explore
			if msg.String() == "b" && !m.wordList.SettingFilter() {
				m.store.ToggleBookmark(m.selectedWord)
				m.updateViewportContent()
				return m, nil
			}
			if msg.String() == "left" || msg.String() == "h" {
				if m.focusedRelIndex > 0 {
					m.focusedRelIndex--
					m.updateViewportContent()
				}
				return m, nil
			}
			if msg.String() == "right" || msg.String() == "l" {
				if m.focusedRelIndex < len(m.relatedItems)-1 {
					m.focusedRelIndex++
					m.updateViewportContent()
				}
				return m, nil
			}
			if msg.String() == "enter" && len(m.relatedItems) > 0 && !m.wordList.SettingFilter() {
				targetWord := m.relatedItems[m.focusedRelIndex].Word
				return m, m.loadWord(targetWord)
			}

		case 1: // Review
			if msg.String() == " " {
				m.isCardFlipped = !m.isCardFlipped
				return m, nil
			}

		case 2: // Settings
			if msg.String() == "enter" {
				m.cfg.WordnikAPIKey = strings.TrimSpace(m.apiKeyInput.Value())
				_ = storage.SaveConfig(m.cfg)
				m.apiClient.SetAPIKey(m.cfg.WordnikAPIKey)
				m.statusMsg = "Settings saved!"
				m.apiKeyInput.Blur()
				return m, nil
			}
			if msg.String() == "c" && !m.apiKeyInput.Focused() {
				cache := api.NewDiskCache()
				cache.Items = make(map[string]api.CachedItem)
				_ = cache.Save()
				m.statusMsg = "Cache cleared!"
				return m, nil
			}
		}
	}

	// Route events to active components
	if m.activeTab == 0 {
		var cmd tea.Cmd
		m.wordList, cmd = m.wordList.Update(msg)
		cmds = append(cmds, cmd)

		if selectedItem, ok := m.wordList.SelectedItem().(wordListItem); ok {
			if selectedItem.word != m.selectedWord {
				cmds = append(cmds, m.loadWord(selectedItem.word))
			}
		}

		m.exploreViewport, cmd = m.exploreViewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.activeTab == 2 {
		var cmd tea.Cmd
		m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) submitSRSRating(rating srs.Rating) {
	if m.srsCurrentIdx >= len(m.srsDueWords) {
		return
	}
	word := m.srsDueWords[m.srsCurrentIdx]
	card := m.store.GetCard(word)

	updatedCard := srs.CalculateNextReview(card, rating, time.Now())
	_ = m.store.UpdateCard(updatedCard)

	m.srsCurrentIdx++
	m.isCardFlipped = false

	if m.srsCurrentIdx < len(m.srsDueWords) {
		m.loadReviewCard(m.srsDueWords[m.srsCurrentIdx])
	} else {
		m.refreshSRSQueue()
	}
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing Wordy TUI..."
	}

	var sb strings.Builder

	// Top Bar
	title := styles.LogoStyle.Render("🧋 Wordy")

	tabs := []string{"[1] Explore", "[2] Review (SRS)", "[3] Settings"}
	tabView := ""
	for i, t := range tabs {
		if i == m.activeTab {
			tabView += styles.TabActiveStyle.Render(t)
		} else {
			tabView += styles.TabInactiveStyle.Render(t)
		}
	}

	rl := m.apiClient.GetRateLimitInfo()
	statusPill := ""
	if m.cfg.WordnikAPIKey == "" {
		statusPill = styles.PillSeed.Render("🍃 Seed Mode")
	} else if rl.IsRateLimited {
		statusPill = styles.PillRateLimit.Render("⚡ Rate Limited")
	} else {
		statusPill = styles.PillOnline.Render(fmt.Sprintf("⚡ API: %d rem", rl.RemainingMinute))
	}

	topHeader := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", tabView, "  ", statusPill)
	sb.WriteString(topHeader + "\n\n")

	// Content Area Bounded strictly by Viewport Height
	switch m.activeTab {
	case 0: // Explore
		sidebar := m.wordList.View()
		cardView := styles.MainBoxStyle.Width(m.exploreViewport.Width).Render(m.exploreViewport.View())
		content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, "  ", cardView)
		sb.WriteString(content)

	case 1: // Review SRS
		reviewView := views.RenderReviewView(
			m.selectedWordDetails,
			m.srsCurrentCard,
			m.isCardFlipped,
			len(m.srsDueWords)-m.srsCurrentIdx,
			m.srsTotalMastered,
			m.store.Data.StreakDays,
			m.width-4,
			m.height-8,
		)
		sb.WriteString(reviewView)

	case 2: // Settings
		cachedCount := len(api.NewDiskCache().Items)
		settView := views.RenderSettingsView(
			m.cfg,
			m.apiKeyInput,
			rl,
			cachedCount,
			m.statusMsg,
			m.width-4,
			m.height-8,
		)
		sb.WriteString(settView)
	}

	// Bottom Help Bar
	sb.WriteString("\n")
	helpText := "Tab view • ↑/↓ navigate list • ←/→ select related • Enter jump • Space flip card • 1-4 rate • b bookmark • q quit"
	sb.WriteString(styles.HelpBar.Width(m.width).Render(helpText))

	return sb.String()
}
