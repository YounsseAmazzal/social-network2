package utils

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"time"

	"social/internal/models"
)

var (
	rxEmail    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rxNickname = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	rxName     = regexp.MustCompile(`^[a-zA-Z]+$`)
)

var (
	ErrRequiredFields = errors.New("all required fields must be filled")
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrPasswordMatch  = errors.New("passwords do not matchhhhh")
	ErrAgeLimit       = errors.New("youuur age must be between 18 and 100")
)

func ValidateRegistration(u *models.UserInfo) error {
	//check emmmty fildes
	if isAnyEmpty(u.Email, u.Password, u.Firstname, u.Lastname, u.DateOfBirth) {
		return ErrRequiredFields
	}

	if u.Password != u.ConfirmPassword {
		return ErrPasswordMatch
	}

	if !rxEmail.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	if !rxName.MatchString(u.Firstname) || !rxName.MatchString(u.Lastname) {
		return errors.New("names must contain only letters")
	}
	if u.Nickname != "" && !rxNickname.MatchString(u.Nickname) {
		return errors.New("nickname must be 3-20 characters (letters, numbers, _, -)")
	}

	userage, err := calculateUserAge(u.DateOfBirth)
	if err != nil {
		return errors.New("invalid date of birth format (use YYYY-MM-DD)")
	}
	if userage < 18 || userage > 100 {
		return ErrAgeLimit
	}
	u.Age = userage

	if len(u.AboutMe) > 1000 {
		return errors.New("about me description is too long (max 1000 chars)")
	}

	sanitizeUser(u)
	return nil
}

// caaalculate the agee of user
func calculateUserAge(birthDateStr string) (int, error) {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	// fmt.Println("the returned age",age)
	return age, nil
}

func isAnyEmpty(fields ...string) bool {
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return true
		}
	}
	return false
}

// maaeke sanitizaaation
func sanitizeUser(u *models.UserInfo) {
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Firstname))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.AboutMe = html.EscapeString(strings.TrimSpace(u.AboutMe))
}
