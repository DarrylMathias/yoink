package tokenizer

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

var (
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

func Tokenize(text string) ([]string, error) {
	text, err := normalize(text)
	if err != nil{
		return nil, err
	}

	if text == "" {
		return nil, nil
	}
	filteredText, err := StemWords(RemoveStopWords(strings.Fields(text)))
	if err != nil{
		return nil, err
	}

	fmt.Println("list of generated tokens: ", filteredText)
	return filteredText, nil
}

func RemoveStopWords(text []string) ([]string){
	stopWords := []string{"a", "an", "also", "and", "are", "as", "at", "be", "because", "been", "but", "by", "for", "from", "have", "has", "however", "is", "it", "if", "not", "of", "on", "or", "so", "than", "that", "the", "their", "there", "these", "this", "was", "were", "whatever", "whether", "which", "with", "would", }

	filteredText := slices.DeleteFunc(text, func(word string) bool {
		return slices.Contains(stopWords, word)
	})
	return filteredText
}

// stemming is the process of converting partciples to nouns
func StemWords(text []string) ([]string, error){
	for i, word  := range text {
		stemmed, err := snowball.Stem(word, "english", true)
		if err != nil{
			return nil, err
		}
		text[i] = stemmed
	}
	return text, nil
}

func normalize(text string) (string, error) {
	text = strings.ToLower(text)

	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		var err error
		switch {
		case unicode.IsLetter(r):
			_, err = builder.WriteRune(r)

		case unicode.IsDigit(r):
			_, err = builder.WriteRune(r)

		case unicode.IsSpace(r):
			_, err = builder.WriteRune(' ')

		default:
			_, err = builder.WriteRune(' ')
		}
		if err != nil{
			return "", nil
		}
	}

	text = builder.String()

	text = whitespaceRegex.ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text), nil
}