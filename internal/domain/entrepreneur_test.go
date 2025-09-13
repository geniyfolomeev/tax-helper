package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestCalculateNextDeclarationDate(t *testing.T) {
	tests := []struct {
		name         string
		now          time.Time
		registeredAt time.Time
		lastSentAt   time.Time
		want         time.Time
	}{
		{
			name:         "No declarations yet, registration last month",
			now:          time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			registeredAt: time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
			lastSentAt:   time.Time{},
			want:         time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "No declarations yet, registration current month",
			now:          time.Date(2025, 9, 10, 10, 20, 0, 0, time.UTC),
			registeredAt: time.Date(2025, 9, 9, 0, 0, 0, 0, time.UTC),
			lastSentAt:   time.Time{},
			want:         time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "Old Entrepreneur, did not send for current month, now is sending day",
			now:          time.Date(2025, 9, 10, 20, 10, 30, 0, time.UTC),
			registeredAt: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			lastSentAt:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "Old Entrepreneur, sent for current month, now is sending day",
			now:          time.Date(2025, 9, 10, 20, 10, 30, 0, time.UTC),
			registeredAt: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			lastSentAt:   time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "Old Entrepreneur, sent for current month, now is sending day (New year)",
			now:          time.Date(2025, 12, 10, 20, 10, 30, 0, time.UTC),
			registeredAt: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			lastSentAt:   time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
			want:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentTimeFn = func() time.Time { return tt.now }
			defer func() { currentTimeFn = time.Now }()

			e := &Entrepreneur{
				TelegramID:      1,
				Status:          "active",
				RegisteredAt:    tt.registeredAt,
				LastSentAt:      tt.lastSentAt,
				YearTotalAmount: decimal.NewFromFloat(10),
			}
			actual := e.CalculateNextDeclarationDate()
			if !actual.Equal(tt.want) {
				t.Errorf("got %v, want %v", actual, tt.want)
			}
		})
	}
}
