package password

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
	"unicode"

	"github.com/nilorg/sdk/random"
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

// RandomPassword 随机生成密码
func RandomPassword(strength PasswordStrength) string {
	// 根据密码强度级别生成密码
	switch strength {
	case VeryWeak:
		return random.Number(6)
	case Weak:
		return random.Number(8)
	case Medium:
		return randomaz(4) + randomNumber(4)
	case Strong:
		return randomaz(4) + randomNumber(4) + randomAZ(4)
	case VeryStrong:
		return randomaz(4) + randomNumber(4) + randomAZ(4) + randomSpecificSymbol(4)
	default:
		return ""
	}
}

// randomNumber 随机数字
func randomNumber(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString(strconv.Itoa(r.Intn(10)))
	}
	return buffer.String()
}

// randomaz 随机a-z字符串
func randomaz(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	strs := []byte("abcdefghijklmnopqrstuvwxyz")
	var buffer bytes.Buffer
	for i := 0; i < l; i++ {
		buffer.WriteByte(strs[r.Intn(len(strs))])
	}
	return buffer.String()
}

// randomAZ 随机A-Z字符串
func randomAZ(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	strs := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var buffer bytes.Buffer
	for i := 0; i < l; i++ {
		buffer.WriteByte(strs[r.Intn(len(strs))])
	}
	return buffer.String()
}

// randomSpecificSymbol 随机特殊字符
func randomSpecificSymbol(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	strs := []byte("!@#%&*?")
	var buffer bytes.Buffer
	for i := 0; i < l; i++ {
		buffer.WriteByte(strs[r.Intn(len(strs))])
	}
	return buffer.String()
}
