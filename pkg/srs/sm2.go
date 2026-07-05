package srs

import (
	"math"
	"time"
)

// Rating represents user recall quality during review.
type Rating int

const (
	RatingAgain Rating = 1 // Failed recall, reset repetition
	RatingHard  Rating = 2 // Remembered with significant effort
	RatingGood  Rating = 3 // Standard successful recall
	RatingEasy  Rating = 4 // Perfect recall with minimal effort
)

// CardState maintains the SRS metrics for a specific word.
type CardState struct {
	Word           string    `json:"word"`
	IntervalDays   int       `json:"interval_days"`
	Repetitions    int       `json:"repetitions"`
	EaseFactor     float64   `json:"ease_factor"`
	NextReviewDate time.Time `json:"next_review_date"`
	LastReviewed   time.Time `json:"last_reviewed"`
	TotalReviews   int       `json:"total_reviews"`
}

// NewCardState initializes a new word card for SRS tracking.
func NewCardState(word string) CardState {
	return CardState{
		Word:           word,
		IntervalDays:   0,
		Repetitions:    0,
		EaseFactor:     2.5, // Default SM-2 ease factor
		NextReviewDate: time.Now(),
		LastReviewed:   time.Time{},
		TotalReviews:   0,
	}
}

// CalculateNextReview computes the next review date, interval, and ease factor using SM-2.
func CalculateNextReview(state CardState, rating Rating, now time.Time) CardState {
	if rating < RatingAgain || rating > RatingEasy {
		rating = RatingGood
	}

	state.LastReviewed = now
	state.TotalReviews++

	// Map 1-4 rating scale to SM-2 0-5 scale equivalence:
	// 1 (Again) -> 1
	// 2 (Hard)  -> 3
	// 3 (Good)  -> 4
	// 4 (Easy)  -> 5
	qMap := map[Rating]float64{
		RatingAgain: 1.0,
		RatingHard:  3.0,
		RatingGood:  4.0,
		RatingEasy:  5.0,
	}
	q := qMap[rating]

	// Update Ease Factor: EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
	ef := state.EaseFactor + (0.1 - (5.0-q)*(0.08+(5.0-q)*0.02))
	if ef < 1.3 {
		ef = 1.3
	}
	state.EaseFactor = math.Round(ef*100) / 100

	if rating == RatingAgain {
		state.Repetitions = 0
		state.IntervalDays = 1
	} else {
		if state.Repetitions == 0 {
			state.IntervalDays = 1
		} else if state.Repetitions == 1 {
			state.IntervalDays = 6
		} else {
			state.IntervalDays = int(math.Round(float64(state.IntervalDays) * state.EaseFactor))
			if state.IntervalDays < 1 {
				state.IntervalDays = 1
			}
		}
		state.Repetitions++
	}

	state.NextReviewDate = now.Add(time.Duration(state.IntervalDays) * 24 * time.Hour)
	return state
}

// IsDue returns true if the card is due for review today.
func (c CardState) IsDue(now time.Time) bool {
	if c.TotalReviews == 0 {
		return true // New cards are due immediately
	}
	return now.After(c.NextReviewDate) || now.Equal(c.NextReviewDate)
}
