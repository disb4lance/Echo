package entity

import (
	"testing"
	"time"
)

func TestUserProfile_CalculateAge(t *testing.T) {
	// Используем реальную текущую дату
	now := time.Now()
	currentYear := now.Year()

	tests := []struct {
		name      string
		birthDate time.Time
		expected  int
	}{
		{
			name:      "Точный возраст - день рождения уже был",
			birthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:  currentYear - 1990, // динамически
		},
		{
			name:      "День рождения еще не наступил в этом году",
			birthDate: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC),
			expected:  currentYear - 1990 - 1,
		},
		{
			name:      "Сегодня день рождения",
			birthDate: time.Date(1990, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected:  currentYear - 1990,
		},
		{
			name:      "Новорожденный",
			birthDate: time.Date(currentYear-1, 12, 31, 0, 0, 0, 0, time.UTC),
			expected:  0,
		},
		{
			name:      "Ровно 1 год",
			birthDate: time.Date(currentYear-1, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &UserProfile{
				BirthDate: tt.birthDate,
			}

			age := profile.CalculateAge()
			if age != tt.expected {
				t.Errorf("CalculateAge() = %v, want %v", age, tt.expected)
			}
		})
	}
}
