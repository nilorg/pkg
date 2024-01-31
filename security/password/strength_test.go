package password

import (
	"testing"
)

func TestEvaluatePasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		expected PasswordStrength
	}{
		{"password", VeryWeak},
		{"Password", Weak},
		{"Password1", Medium},
		{"Password123", Strong},
		{"Password123!", VeryStrong},
		{"", VeryWeak},
	}

	for _, test := range tests {
		strength := EvaluatePasswordStrength(test.password)
		if strength != test.expected {
			t.Errorf("EvaluatePasswordStrength(%s) = %d, expected %d", test.password, strength, test.expected)
		}
	}
}

func TestRandomPassword(t *testing.T) {
	tests := []struct {
		strength PasswordStrength
	}{
		{VeryWeak},
		{Weak},
		{Medium},
		{Strong},
		{VeryStrong},
	}

	for _, test := range tests {
		password := RandomPassword(test.strength)
		strength := EvaluatePasswordStrength(password)
		t.Logf("RandomPassword(%d) = %s", test.strength, password)
		if strength < test.strength {
			t.Errorf("RandomPassword(%d) = %s, expected strength >= %d", test.strength, password, test.strength)
		}
	}
}
