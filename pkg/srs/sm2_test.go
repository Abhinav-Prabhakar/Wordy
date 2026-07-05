package srs

import (
	"testing"
	"time"
)

func TestCalculateNextReview(t *testing.T) {
	now := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	card := NewCardState("perspicacious")

	if !card.IsDue(now) {
		t.Errorf("New card should be due immediately")
	}

	// First review: Good (rating 3)
	card = CalculateNextReview(card, RatingGood, now)
	if card.Repetitions != 1 {
		t.Errorf("Expected 1 repetition, got %d", card.Repetitions)
	}
	if card.IntervalDays != 1 {
		t.Errorf("Expected 1 day interval, got %d", card.IntervalDays)
	}

	// Second review: Easy (rating 4)
	now = now.Add(24 * time.Hour)
	card = CalculateNextReview(card, RatingEasy, now)
	if card.Repetitions != 2 {
		t.Errorf("Expected 2 repetitions, got %d", card.Repetitions)
	}
	if card.IntervalDays != 6 {
		t.Errorf("Expected 6 days interval, got %d", card.IntervalDays)
	}

	// Third review: Again (rating 1)
	card = CalculateNextReview(card, RatingAgain, now)
	if card.Repetitions != 0 {
		t.Errorf("Expected 0 repetitions after reset, got %d", card.Repetitions)
	}
	if card.IntervalDays != 1 {
		t.Errorf("Expected 1 day interval after reset, got %d", card.IntervalDays)
	}
}
