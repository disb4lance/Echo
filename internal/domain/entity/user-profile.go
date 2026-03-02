package entity

import (
	"database/sql/driver"
	"fmt"
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

func (g *Gender) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Gender", value)
	}

	switch str {
	case "male":
		*g = Male
	case "female":
		*g = Female
	case "other":
		*g = Other
	default:
		*g = PreferNotToSay
	}

	return nil
}

func (g Gender) Value() (driver.Value, error) {
	return g.String(), nil
}

func (g Gender) IsValid() bool {
	switch g {
	case PreferNotToSay, Male, Female, Other:
		return true
	}
	return false
}

func (u *UserProfile) CalculateAge() int {
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()

	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}
