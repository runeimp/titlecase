package titlecase

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type tcCallback func(string, ...interface{}) string

type Mutagenic interface {
	AddValue(string)
	String() string
}

func appendage(tcLine []Mutagenic, word string) []Mutagenic {
	return append(tcLine, NewMutable(word))
}

type Mutagen struct {
	Value string
}

func (m Mutagen) AddValue(v string) {
	m.Value = v
}

func (m Mutagen) String() string {
	return m.Value
}

type Mutable struct {
	Mutagen
}

type Immutable struct {
	Mutagen
}

func NewImmutable(v string) Immutable {
	return Immutable{Mutagen{Value: v}}
}

func NewMutable(v string) Mutable {
	return Mutable{Mutagen{Value: v}}
}

const (
	smallWords  = `a|an|and|as|at|but|by|en|for|if|in|of|on|or|the|to|v\.?|via|vs\.?|with`
	punctuation = `!"“#$%&'‘()*+,\-–‒—―./:;?@[\\\]_` + "`" + `{|}~`
	Version     = "1.0.0"
)

var (
	AposSecond   = regexp.MustCompile(`^[dol]{1}['‘]{1}[a-z]+(?:['s]{2})?$`) // case insensitive
	CapFirst     = regexp.MustCompile(fmt.Sprintf(`^[%s]*?([A-Za-z])`, punctuation))
	InlinePeriod = regexp.MustCompile(`[a-z][.][a-z]`) // case insensitive
	MacMc        = regexp.MustCompile(`^([Mm]c|MC)(\w.+)`)
	SmallFirst   = regexp.MustCompile(fmt.Sprintf(`^([%s]*)(%s)\b`, punctuation, smallWords)) // case insensitive
	SmallLast    = regexp.MustCompile(fmt.Sprintf(`\b(%s)[%s]?$`, smallWords, punctuation))   // case insensitive
	SmallWords   = regexp.MustCompile(fmt.Sprintf(`^(%s)$`, smallWords))                      // case insensitive
	SubPhrase    = regexp.MustCompile(fmt.Sprintf(`([:.;?!\-–‒—―][ ])(%s)`, smallWords))
	UCElsewhere  = regexp.MustCompile(fmt.Sprintf(`[%s]*?[a-zA-Z]+[A-Z]+?`, punctuation))
	UCInitials   = regexp.MustCompile(`^(?:[A-Z]{1}\.{1}|[A-Z]{1}\.{1}[A-Z]{1})+$`)
)

// capitalize produces results the same way that Python's str.capitalize() does versus the way Go's strings.Title() works
func capitalize(s string) string {
	s = strings.ToLower(s)
	s = strings.Title(s)
	return s
}

/**
 * Convert input text
 *
 * This filter changes all words to Title Caps, and attempts to be clever
 * about *un*capitalizing small words like a/an/the in the input.
 *
 * The list of "small words" which are not capped comes from
 * the New York Times Manual of Style, plus 'vs', 'v', and 'with'.
 */
func Convert(input string, args ...interface{}) string {
	var (
		callback       tcCallback = nil
		processed                 = []string{}
		smallFirstLast            = true
	)

	for _, arg := range args {
		switch t := arg.(type) {
		case bool:
			smallFirstLast = arg.(bool)
		case tcCallback:
			callback = arg.(tcCallback)
		default:
			fmt.Println("Unhandled arg type: %s", t)
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("titlecase.Convert() | line = %q\n", line)
		lineAllCaps := false
		if strings.ToUpper(line) == line {
			lineAllCaps = true
		}
		words := strings.Fields(line)
		tcLine := []Mutagenic{}
		for _, word := range words {
			// wordAllCaps := false
			// if strings.ToUpper(word) == word {
			// 	wordAllCaps = true
			// }

			if callback != nil {
				newWord := callback(word, lineAllCaps)
				if newWord != "" {
					tcLine = append(tcLine, NewImmutable(newWord))
					continue
				}
			}

			if lineAllCaps && UCInitials.MatchString(word) {
				// fmt.Printf("titlecase.Convert() | UCInitials   | word = %q\n", word)
				tcLine = appendage(tcLine, word)
				continue
			}

			if AposSecond.MatchString(word) {
				// fmt.Printf("titlecase.Convert() | AposSecond   | word = %q\n", word)
				if len(word) > 0 && !strings.Contains("aeiouAEIOU", word[0:1]) {
					word = strings.ToLower(word[0:1]) + word[1:2] + strings.ToUpper(word[2:3]) + word[3:]
				} else {
					word = strings.ToUpper(word[0:1]) + word[1:2] + strings.ToUpper(word[2:3]) + word[3:]
				}
				tcLine = appendage(tcLine, word)
				continue
			}

			if MacMc.MatchString(word) {
				// fmt.Printf("titlecase.Convert() | MacMc       | word = %q\n", word)
				matches := MacMc.FindStringSubmatch(word)
				tcLine = appendage(
					tcLine,
					fmt.Sprintf("%s%s",
						capitalize(matches[1]),
						Convert(matches[2], callback, smallFirstLast),
					),
				)
				continue
			}

			if InlinePeriod.MatchString(word) || (!lineAllCaps && UCElsewhere.MatchString(word)) {
				// fmt.Printf("titlecase.Convert() | InlinePeriod | lineAllCaps = %t | wordAllCaps = %t | word = %q\n", lineAllCaps, wordAllCaps, word)
				tcLine = appendage(tcLine, word)
				continue
			}

			if SmallWords.MatchString(word) {
				// fmt.Printf("titlecase.Convert() | SmallWords   | word = %q\n", word)
				tcLine = appendage(tcLine, strings.ToLower(word))
				continue
			}

			if strings.Contains(word, "/") && !strings.Contains(word, "//") {
				// fmt.Printf("titlecase.Convert() | Slashes      | word = %q\n", word)
				slashed := titleDelimited(word, "-", callback, smallFirstLast)
				tcLine = appendage(tcLine, slashed)
				continue
			}

			if strings.Contains("-", word) {
				// fmt.Printf("titlecase.Convert() | Hyphens      | word = %q\n", word)
				hyphenated := titleDelimited(word, "-", callback, smallFirstLast)
				tcLine = appendage(tcLine, hyphenated)
				continue
			}

			if lineAllCaps {
				// fmt.Printf("titlecase.Convert() | ALLCAPS      | word = %q\n", word)
				word = strings.ToLower(word)
			}

			// fmt.Printf("titlecase.Convert() | CapFirst     | word = %q\n", word)
			tcLine = appendage(tcLine, CapFirst.ReplaceAllStringFunc(word, strings.ToUpper))
		} // end of word loop

		if smallFirstLast && len(tcLine) > 0 {
			switch tcLine[0].(type) {
			case Immutable:
				// ignore
			default:
				tcLine[0] = NewMutable(SmallFirst.ReplaceAllStringFunc(tcLine[0].String(), capitalize))
			}

			last := len(tcLine) - 1
			switch tcLine[last].(type) {
			case Immutable:
				// ignore
			default:
				tcLine[last] = NewMutable(SmallLast.ReplaceAllStringFunc(tcLine[last].String(), capitalize))
			}
		}

		result := ""
		// result = strings.Join(tcLine, " ")
		for _, word := range tcLine {
			result = result + " " + word.String()
		}
		result = result[1:]

		// NOTE: Not sure about the following block
		result = SubPhrase.ReplaceAllStringFunc(
			result,
			func(phrase string) string {
				matches := SubPhrase.FindStringSubmatch(phrase)
				return fmt.Sprintf("%s%s",
					matches[1],
					capitalize(matches[2]),
				)
			},
		)

		processed = append(processed, result)

	} // end of line loop

	return strings.Join(processed, "\n")
}

func titleDelimited(s, d string, callback tcCallback, smallFirstLast bool) string {
	delimited := strings.Split(s, d)
	for i, t := range delimited {
		delimited[i] = Convert(t, callback, smallFirstLast)
	}
	return strings.Join(delimited, d)
}
