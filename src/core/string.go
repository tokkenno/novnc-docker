package core

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	randPtr = rand.New(rand.NewSource(time.Now().UnixNano()))
	lowerCaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
)

func randString(n int, set []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = set[randPtr.Intn(len(set))]
	}
	return string(b)
}

func GenerateUrl(from string) string {
	// Make a Regex to say we only want letters and numbers
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	cleanString := reg.ReplaceAllString(strings.Replace(strings.ToLower(from), " ", "_", -1), "")
	return cleanString + "_" + randString(10, lowerCaseRunes)
}