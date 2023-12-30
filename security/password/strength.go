package password

import (
	"unicode"
)

type PasswordStrength int // 密码强度级别

const (
	VeryWeak   PasswordStrength = iota // 非常弱
	Weak                               // 弱
	Medium                             // 中等
	Strong                             // 强
	VeryStrong                         // 非常强
)

// EvaluatePasswordStrength 评估密码强度，并返回相应的密码强度级别
func EvaluatePasswordStrength(password string) PasswordStrength {
	// 规则包括：
	// 密码长度至少为8个字符。
	// 包含至少一个小写字母。
	// 包含至少一个大写字母。
	// 包含至少一个数字。
	// 包含至少一个特殊字符（标点符号或符号）。
	// 密码长度至少为12个字符。

	length := len(password)
	if length < 8 {
		return VeryWeak
	}

	var (
		hasLowercase bool
		hasUppercase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var strength PasswordStrength

	if hasLowercase {
		strength++
	}
	if hasUppercase {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength++
	}
	if length >= 12 {
		strength++
	}

	return strength
}
