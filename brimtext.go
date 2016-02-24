// Package brimtext contains tools for working with text. Probably the most
// complex of these tools is Align, which allows for formatting "pretty
// tables".
package brimtext

import (
	"bytes"
	"strconv"
	"strings"
)

// OrdinalSuffix returns "st", "nd", "rd", etc. for the number given (1st, 2nd,
// 3rd, etc.).
func OrdinalSuffix(number int) string {
	if (number/10)%10 == 1 || number%10 > 3 {
		return "th"
	} else if number%10 == 1 {
		return "st"
	} else if number%10 == 2 {
		return "nd"
	} else if number%10 == 3 {
		return "rd"
	}
	return "th"
}

// ThousandsSep returns the number formatted using the separator at each
// thousands position, such as ThousandsSep(1234567, ",") giving 1,234,567.
func ThousandsSep(v int64, sep string) string {
	s := strconv.FormatInt(v, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// ThousandsSepU returns the number formatted using the separator at each
// thousands position, such as ThousandsSepU(1234567, ",") giving 1,234,567.
func ThousandsSepU(v uint64, sep string) string {
	s := strconv.FormatUint(v, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

type humanSize struct {
	d int64
	s string
}

var humanSizes = []humanSize{
	humanSize{int64(1024), "K"},
	humanSize{int64(1024) << 10, "M"},
	humanSize{int64(1024) << 20, "G"},
	humanSize{int64(1024) << 30, "T"},
	humanSize{int64(1024) << 40, "P"},
	humanSize{int64(1024) << 50, "E"},
}

// Returns a more readable size format, such as HumanSize(1234567, "") giving
// "1M". For values less than 1K, it is common that no suffix letter should be
// added; but the appendBytes parameter is provided in case clarity is needed.
func HumanSize(b int64, appendBytes string) string {
	if b < 1024 {
		v := strconv.FormatInt(b, 10)
		if appendBytes != "" {
			return v + appendBytes
		}
		return v
	}
	c := b
	s := appendBytes
	for _, h := range humanSizes {
		c = b / h.d
		r := b % h.d
		if r >= h.d/2 {
			c++
		}
		if c < 1024 {
			s = h.s
			break
		}
	}
	return strconv.FormatInt(c, 10) + s
}

// Sentence converts the value into a sentence, uppercasing the first character
// and ensuring the string ends with a period. Useful to output better looking
// error.Error() messages, which are all lower case with no trailing period by
// convention.
func Sentence(value string) string {
	if value != "" {
		if value[len(value)-1] != '.' {
			value = strings.ToUpper(value[:1]) + value[1:] + "."
		} else {
			value = strings.ToUpper(value[:1]) + value[1:]
		}
	}
	return value
}

// StringSliceToLowerSort provides a sort.Interface that will sort a []string
// by their strings.ToLower values. This isn't exactly a case insensitive sort
// due to Unicode situations, but is usually good enough.
type StringSliceToLowerSort []string

func (s StringSliceToLowerSort) Len() int {
	return len(s)
}

func (s StringSliceToLowerSort) Swap(x int, y int) {
	s[x], s[y] = s[y], s[x]
}

func (s StringSliceToLowerSort) Less(x int, y int) bool {
	return strings.ToLower(s[x]) < strings.ToLower(s[y])
}

// Wrap wraps text for more readable output.
//
// The width can be a positive int for a specific width, 0 for the default
// width (attempted to get from terminal, 79 otherwise), or a negative number
// for a width relative to the default.
//
// The indent1 is the prefix for the first line.
//
// The indent2 is the prefix for any second or subsequent lines.
func Wrap(text string, width int, indent1 string, indent2 string) string {
	if width < 1 {
		width = GetTTYWidth() - 1 + width
	}
	bs := []byte(text)
	bs = wrap(bs, width, []byte(indent1), []byte(indent2))
	return string(bytes.Trim(bs, "\n"))
}

func wrap(text []byte, width int, indent1 []byte, indent2 []byte) []byte {
	if len(text) == 0 {
		return text
	}
	text = bytes.Replace(text, []byte{'\r', '\n'}, []byte{'\n'}, -1)
	var out bytes.Buffer
	for _, par := range bytes.Split([]byte(text), []byte{'\n', '\n'}) {
		par = bytes.Replace(par, []byte{'\n'}, []byte{' '}, -1)
		lineLen := 0
		start := true
		for _, word := range bytes.Split(par, []byte{' '}) {
			wordLen := len(word)
			if wordLen == 0 {
				continue
			}
			scan := word
			for len(scan) > 1 {
				i := bytes.IndexByte(scan, '\x1b')
				if i == -1 {
					break
				}
				j := bytes.IndexByte(scan[i+1:], 'm')
				if j == -1 {
					i++
				} else {
					j += 2
					wordLen -= j
					scan = scan[i+j:]
				}
			}
			if start {
				out.Write(indent1)
				lineLen += len(indent1)
				out.Write(word)
				lineLen += wordLen
				start = false
			} else if lineLen+1+wordLen > width {
				out.WriteByte('\n')
				out.Write(indent2)
				out.Write(word)
				lineLen = len(indent2) + wordLen
			} else {
				out.WriteByte(' ')
				out.Write(word)
				lineLen += 1 + wordLen
			}
		}
		out.WriteByte('\n')
		out.WriteByte('\n')
	}
	return out.Bytes()
}

// AllEqual returns true if all the values are equal strings; no strings,
// AllEqual() or AllEqual([]string{}...), are considered AllEqual.
func AllEqual(values ...string) bool {
	if len(values) < 2 {
		return true
	}
	compare := values[0]
	for _, v := range values[1:] {
		if v != compare {
			return false
		}
	}
	return true
}
