package password

import (
	"testing"
)

func TestValidatePasswordComplexity(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErrs int
	}{
		{"empty password", "", 5},
		{"too short", "Aa1!", 1},
		{"no uppercase", "abcdefg1!", 1},
		{"no lowercase", "ABCDEFG1!", 1},
		{"no digit", "Abcdefgh!", 1},
		{"no special char", "Abcdefg1", 1},
		{"valid password", "Abcdefg1!", 0},
		{"all failures except length", "abcdefgh", 3},
		{"exactly 8 chars valid", "Abcdef1!", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateComplexity(tt.password)
			if len(errs) != tt.wantErrs {
				t.Errorf(
					"ValidateComplexity(%q) returned %d errors, want %d: %v",
					tt.password,
					len(errs),
					tt.wantErrs,
					errs,
				)
			}
		})
	}
}
