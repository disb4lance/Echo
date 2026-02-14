package entity

import (
	"time"

	"github.com/google/uuid"
)

type Gender int
type City int

const (
	PreferNotToSay Gender = iota
	Male
	Female
	Other
)

const (
	Moscow City = iota
	SaintPetersburg
)

func (g Gender) String() string {
	switch g {
	case Male:
		return "male"
	case Female:
		return "female"
	case Other:
		return "other"
	default:
		return "prefer_not_to_say"
	}
}

func (g Gender) IsValid() bool {
	switch g {
	case PreferNotToSay, Male, Female, Other:
		return true
	}
	return false
}

type UserProfile struct {
	UserID      uuid.UUID
	Name        string
	BirthDate   time.Time
	Gender      Gender
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *UserProfile) CalculateAge() int {
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()

	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}
