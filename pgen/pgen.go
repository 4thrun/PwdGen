package pgen

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	// lowercase set
	Lowercase = "abcdefghijklmnopqrstuvwxyz"
	// uppercase set
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// digits set
	Digits = "1234567890"
	// special symbols set
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

// error types
var (
	// total length greater than designed
	ErrExceedsTotalLength = errors.New("too long")
	// number of letters exceeds the number of available letters
	// and repeats are not allowed
	ErrLettersExceedsAvailable = errors.New("excessive letters, repeats not allowed")
	// number of digits exceeds the number of available digits
	// and repeats are not allowed
	ErrDigitsExceedsAvailable = errors.New("excessive digits, repeats not allowed")
	// number of symbols exceeds the number of available symbols
	// and repeats are not allowed
	ErrSymbolsExceedsAvailable = errors.New("excessive symbols, repeats not allowed")
)

// Generator: stateful generator which can be used to customize
// the list of lowercase, uppercase, digits and symbols
type Generator struct {
	lowercase string
	uppercase string
	digits    string
	symbols   string
}

type GenInput struct {
	lowercase_input string
	uppercase_input string
	digits_input    string
	symbols_input   string
}

func NewGenerator(i *GenInput) (*Generator, error) {
	if i == nil {
		i = new(GenInput)
	}
	var g = &Generator{
		lowercase: i.lowercase_input,
		uppercase: i.uppercase_input,
		digits:    i.digits_input,
		symbols:   i.symbols_input,
	}
	if g.lowercase == "" {
		g.lowercase = Lowercase
	}
	if g.uppercase == "" {
		g.uppercase = Uppercase
	}
	if g.digits == "" {
		g.digits = Digits
	}
	if g.symbols == "" {
		g.symbols = Symbols
	}
	return g, nil
}

func (g *Generator) Generate(
	length int,
	numDigits int,
	numSymbols int,
	noUppercase bool,
	allowRepeat bool,
) (string, error) {
	var letters = g.lowercase
	if !noUppercase {
		letters += g.uppercase
	}
	var chars = length - numDigits - numSymbols
	if chars < 0 {
		return "", ErrExceedsTotalLength
	}
	if !allowRepeat && chars > len(letters) {
		return "", ErrLettersExceedsAvailable
	}
	if !allowRepeat && numDigits > len(g.digits) {
		return "", ErrDigitsExceedsAvailable
	}
	if !allowRepeat && numSymbols > len(g.symbols) {
		return "", ErrSymbolsExceedsAvailable
	}

	var result string
	// characters
	for i := 0; i < chars; i++ {
		ch, err := randomElement(letters)
		if err != nil {
			return "", err
		}
		if !allowRepeat && strings.Contains(result, ch) {
			i--
			continue
		}
		result, err = randomInsert(result, ch)
		if err != nil {
			return "", err
		}
	}
	// digits
	for i := 0; i < numDigits; i++ {
		d, err := randomElement(g.digits)
		if err != nil {
			return "", nil
		}
		if !allowRepeat && strings.Contains(result, d) {
			i--
			continue
		}
		result, err = randomInsert(result, d)
		if err != nil {
			return "", err
		}
	}
	// symbols
	for i := 0; i < numSymbols; i++ {
		s, err := randomElement(g.symbols)
		if err != nil {
			return "", nil
		}
		if !allowRepeat && strings.Contains(result, s) {
			i = i - 1
			continue
		}
		result, err = randomInsert(result, s)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

// MustGenerate: the same as Generate, but panics on error
func (g *Generator) MustGenerate(
	length int,
	numDigits int,
	numSymbols int,
	noUppercase bool,
	allowRepeat bool,
) string {
	result, err := g.Generate(length, numDigits, numSymbols, noUppercase, allowRepeat)
	if err != nil {
		panic(err)
	}
	return result
}

// Generate: the package shortcut for Generator.Generate
func Generate(length, numDigits, numSymbols int, noUppercase, allowRepeat bool) (string, error) {
	gen, err := NewGenerator(nil)
	if err != nil {
		return "", err
	}
	return gen.Generate(length, numDigits, numSymbols, noUppercase, allowRepeat)
}

// MustGenerate: the package shortcut for Generator.MustGenerate
func MustGenerate(length, numDigits, numSymbols int, noUppercase, allowRepeat bool) string {
	res, err := Generate(length, numDigits, numSymbols, noUppercase, allowRepeat)
	if err != nil {
		panic(err)
	}
	return res
}

// randomElement: extracts a random element from the given string
func randomElement(s string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}
	return string(s[n.Int64()]), nil
}

// randomInsert: randomly inserts the given character into the given string
func randomInsert(s string, ch string) (string, error) {
	if s == "" {
		return ch, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(1+len(s))))
	if err != nil {
		return "", err
	}
	var i = n.Int64()
	return s[0:i] + ch + s[i:], nil
}
