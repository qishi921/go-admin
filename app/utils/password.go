package utils

import (
	"errors"
	"unicode"
)

// PasswordPolicy defines password requirements.
type PasswordPolicy struct {
	MinLength  int
	MaxLength  int
	RequireUpper bool
	RequireLower bool
	RequireDigit bool
	RequireSpecial bool
}

// DefaultPasswordPolicy returns the default password policy.
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false,
	}
}

// ValidatePassword validates a password against the policy.
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = DefaultPasswordPolicy()
	}

	if len(password) < policy.MinLength {
		return errors.New("密码长度不能少于 " + itoa(policy.MinLength) + " 个字符")
	}
	if len(password) > policy.MaxLength {
		return errors.New("密码长度不能超过 " + itoa(policy.MaxLength) + " 个字符")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return errors.New("密码必须包含至少一个大写字母")
	}
	if policy.RequireLower && !hasLower {
		return errors.New("密码必须包含至少一个小写字母")
	}
	if policy.RequireDigit && !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}
	if policy.RequireSpecial && !hasSpecial {
		return errors.New("密码必须包含至少一个特殊字符")
	}

	return nil
}

// ValidatePasswordSimple validates password with basic rules (min 8 chars, at least one digit and letter).
func ValidatePasswordSimple(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不能少于 8 个字符")
	}
	if len(password) > 128 {
		return errors.New("密码长度不能超过 128 个字符")
	}

	var hasLetter, hasDigit bool
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter {
		return errors.New("密码必须包含至少一个字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}

	return nil
}

// CommonWeakPasswords contains commonly used weak passwords.
var CommonWeakPasswords = map[string]bool{
	"123456":        true,
	"password":      true,
	"12345678":      true,
	"qwerty":        true,
	"123456789":     true,
	"12345":         true,
	"1234":          true,
	"111111":        true,
	"1234567":       true,
	"dragon":        true,
	"123123":        true,
	"baseball":      true,
	"abc123":        true,
	"football":      true,
	"monkey":        true,
	"letmein":       true,
	"696969":        true,
	"shadow":        true,
	"master":        true,
	"666666":        true,
	"qwertyuiop":    true,
	"123321":        true,
	"mustang":       true,
	"1234567890":    true,
	"michael":       true,
	"654321":        true,
	"pussy":         true,
	"superman":      true,
	"1qaz2wsx":      true,
	"7777777":       true,
	"admin":         true,
	"admin123":      true,
	"admin1234":     true,
	"password123":   true,
	"password1":     true,
	"qwerty123":     true,
}

// IsWeakPassword checks if the password is a commonly known weak password.
func IsWeakPassword(password string) bool {
	return CommonWeakPasswords[password]
}

// itoa converts int to string (simple implementation).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}
