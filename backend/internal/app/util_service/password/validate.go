package password

import (
	"unicode"
)

func ValidateComplexity(password string) []string {
	var errs []string

	if len(password) < 8 {
		errs = append(errs, "Must be at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errs = append(errs, "Must contain at least 1 uppercase letter")
	}
	if !hasLower {
		errs = append(errs, "Must contain at least 1 lowercase letter")
	}
	if !hasDigit {
		errs = append(errs, "Must contain at least 1 number")
	}
	if !hasSpecial {
		errs = append(errs, "Must contain at least 1 special character")
	}

	return errs
}
