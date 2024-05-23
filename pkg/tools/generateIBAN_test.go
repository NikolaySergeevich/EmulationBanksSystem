package tools_test

import (
	"strings"
	"testbanc/pkg/tools"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIBAN(t *testing.T) {
	iban := tools.GenerateIBAN()

	// Проверяем, что IBAN начинается с "BY"
	assert.True(t, strings.HasPrefix(iban.Iban, "BY"), "IBAN must start with BY")

	// Проверяем, что длина IBAN равна 28 символам
	assert.Equal(t, 28, utf8.RuneCountInString(iban.Iban), "The IBAN must be 28 characters long")

	// Проверяем, что IBAN состоит из букв и цифр
	for _, char := range iban.Iban {
		assert.True(t, unicode.IsLetter(char) || unicode.IsDigit(char), "IBAN must consist of letters and numbers")
	}
}