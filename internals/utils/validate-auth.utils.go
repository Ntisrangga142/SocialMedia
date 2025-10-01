package utils

import "regexp"

func ValidateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	var (
		lowerRegex   = regexp.MustCompile(`[a-z]`)
		upperRegex   = regexp.MustCompile(`[A-Z]`)
		digitRegex   = regexp.MustCompile(`[0-9]`)
		specialRegex = regexp.MustCompile(`[\W_]`)
	)

	if len(password) < 8 {
		return false
	}
	return lowerRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		digitRegex.MatchString(password) &&
		specialRegex.MatchString(password)
}
