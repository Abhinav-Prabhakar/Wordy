package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wordy-tui/wordy/pkg/srs"
)

// UserData holds saved word data and SRS progress across sessions.
type UserData struct {
	Cards            map[string]srs.CardState `json:"cards"`
	BookmarkedWords  []string                 `json:"bookmarked_words"`
	StreakDays       int                      `json:"streak_days"`
	LastStudyDate    string                   `json:"last_study_date"`
	TotalWordsLearned int                     `json:"total_words_learned"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	Data UserData
}

// GetDataPath returns the path to user data file.
func GetDataPath() string {
	dataDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		dataDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(dataDir, "wordy", "data.json")
}

// NewStore initializes disk storage for user state.
func NewStore() (*Store, error) {
	path := GetDataPath()
	store := &Store{
		path: path,
		Data: UserData{
			Cards:           make(map[string]srs.CardState),
			BookmarkedWords: make([]string, 0),
		},
	}
	err := store.Load()
	return store, err
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var loaded UserData
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	if loaded.Cards == nil {
		loaded.Cards = make(map[string]srs.CardState)
	}
	s.Data = loaded
	return nil
}

func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// GetCard returns SRS card state for a word. Initializes if not present.
func (s *Store) GetCard(word string) srs.CardState {
	s.mu.RLock()
	card, exists := s.Data.Cards[word]
	s.mu.RUnlock()

	if !exists {
		return srs.NewCardState(word)
	}
	return card
}

// UpdateCard updates the SRS card state and updates streak stats if reviewed today.
func (s *Store) UpdateCard(card srs.CardState) error {
	s.mu.Lock()
	s.Data.Cards[card.Word] = card

	today := time.Now().Format("2006-01-02")
	if s.Data.LastStudyDate != today {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if s.Data.LastStudyDate == yesterday {
			s.Data.StreakDays++
		} else if s.Data.LastStudyDate == "" {
			s.Data.StreakDays = 1
		} else {
			s.Data.StreakDays = 1
		}
		s.Data.LastStudyDate = today
	}

	// Update total words learned count (repetitions >= 3)
	learned := 0
	for _, c := range s.Data.Cards {
		if c.Repetitions >= 3 {
			learned++
		}
	}
	s.Data.TotalWordsLearned = learned
	s.mu.Unlock()

	return s.Save()
}

// ToggleBookmark toggles bookmark status for a word.
func (s *Store) ToggleBookmark(word string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, w := range s.Data.BookmarkedWords {
		if w == word {
			s.Data.BookmarkedWords = append(s.Data.BookmarkedWords[:i], s.Data.BookmarkedWords[i+1:]...)
			_ = s.Save()
			return false
		}
	}
	s.Data.BookmarkedWords = append(s.Data.BookmarkedWords, word)
	_ = s.Save()
	return true
}

func (s *Store) IsBookmarked(word string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, w := range s.Data.BookmarkedWords {
		if w == word {
			return true
		}
	}
	return false
}
