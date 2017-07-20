package brimtext

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type Alignment int

const (
	Left Alignment = iota
	Right
	Center
)

type AlignOptions struct {
	// Widths indicate the desired widths of each column. If nil or if a value
	// is 0, no rewrapping will be done.
	Widths     []int
	Alignments []Alignment
	// FirstDR etc. control what is output for situations with a prepended
	// display row, First row output with Down and Right connections, etc.
	FirstDR       string
	FirstLR       string
	FirstFirstDLR string
	FirstDLR      string
	FirstDL       string
	// RowFirstUD etc. control situations for each data row output.
	RowFirstUD  string
	RowSecondUD string
	RowUD       string
	RowLastUD   string
	// LeaveTrailingWhitespace should be set true if the last cell of data row
	// needs spaces to fill to the end (usually needed when setting RowLastUD).
	LeaveTrailingWhitespace bool
	// FirstNilFirstUDR etc. control situations when the first nil data row is
	// encountered. Can be used to separate the header from the rest of the
	// rows.
	FirstNilFirstUDR  string
	FirstNilLR        string
	FirstNilFirstUDLR string
	FirstNilUDLR      string
	FirstNilLastUDL   string
	// NilFirstUDR etc. control situations when the second and subsequent nil
	// data rows are encountered. Can be used to separate rows from each other.
	NilFirstUDR  string
	NilLR        string
	NilFirstUDLR string
	NilUDLR      string
	NilLastUDL   string
	// LastUR etc. control what is output for situations with an appended
	// display row.
	LastUR       string
	LastLR       string
	LastFirstULR string
	LastULR      string
	LastUL       string
	// NilBetweenEveryRow will add a nil data row between all rows; use to emit
	// FirstNil* and Nil* row separators.
	NilBetweenEveryRow bool
}

// NewDefaultAlignOptions gives:
//
//  &AlignOptions{RowSecondUD: " ", RowUD: " "}
//
// Which will format tables like:
//
//           Bob         Sue    John
//  Hometown San Antonio Austin New York
//  Mother   Bessie      Mary   Sarah
//  Father   Rick        Dan    Mike
func NewDefaultAlignOptions() *AlignOptions {
	return &AlignOptions{RowSecondUD: " ", RowUD: " "}
}

// NewSimpleAlignOptions gives:
//
//  return &AlignOptions{
//      FirstDR:                 "+-",
//      FirstLR:                 "-",
//      FirstFirstDLR:           "-+-",
//      FirstDLR:                "-+-",
//      FirstDL:                 "-+",
//      RowFirstUD:              "| ",
//      RowSecondUD:             " | ",
//      RowUD:                   " | ",
//      RowLastUD:               " |",
//      LeaveTrailingWhitespace: true,
//      FirstNilFirstUDR:        "+-",
//      FirstNilLR:              "-",
//      FirstNilFirstUDLR:       "-+-",
//      FirstNilUDLR:            "-+-",
//      FirstNilLastUDL:         "-+",
//      LastUR:                  "+-",
//      LastLR:                  "-",
//      LastFirstULR:            "-+-",
//      LastULR:                 "-+-",
//      LastUL:                  "-+",
//  }
//
// Which will format tables like:
//
//  +----------+-------------+--------+----------+
//  |          | Bob         | Sue    | John     |
//  +----------+-------------+--------+----------+
//  | Hometown | San Antonio | Austin | New York |
//  | Mother   | Bessie      | Mary   | Sarah    |
//  | Father   | Rick        | Dan    | Mike     |
//  +----------+-------------+--------+----------+
func NewSimpleAlignOptions() *AlignOptions {
	return &AlignOptions{
		FirstDR:                 "+-",
		FirstLR:                 "-",
		FirstFirstDLR:           "-+-",
		FirstDLR:                "-+-",
		FirstDL:                 "-+",
		RowFirstUD:              "| ",
		RowSecondUD:             " | ",
		RowUD:                   " | ",
		RowLastUD:               " |",
		LeaveTrailingWhitespace: true,
		FirstNilFirstUDR:        "+-",
		FirstNilLR:              "-",
		FirstNilFirstUDLR:       "-+-",
		FirstNilUDLR:            "-+-",
		FirstNilLastUDL:         "-+",
		LastUR:                  "+-",
		LastLR:                  "-",
		LastFirstULR:            "-+-",
		LastULR:                 "-+-",
		LastUL:                  "-+",
	}
}

// NewBoxedAlignOptions gives:
//  &AlignOptions{
//      FirstDR:                 "+=",
//      FirstLR:                 "=",
//      FirstFirstDLR:           "=+=",
//      FirstDLR:                "=+=",
//      FirstDL:                 "=+",
//      RowFirstUD:              "| ",
//      RowSecondUD:             " | ",
//      RowUD:                   " | ",
//      RowLastUD:               " |",
//      LeaveTrailingWhitespace: true,
//      FirstNilFirstUDR:        "+=",
//      FirstNilLR:              "=",
//      FirstNilFirstUDLR:       "=+=",
//      FirstNilUDLR:            "=+=",
//      FirstNilLastUDL:         "=+",
//      NilFirstUDR:             "+-",
//      NilLR:                   "-",
//      NilFirstUDLR:            "-+-",
//      NilUDLR:                 "-+-",
//      NilLastUDL:              "-+",
//      LastUR:                  "+=",
//      LastLR:                  "=",
//      LastFirstULR:            "=+=",
//      LastULR:                 "=+=",
//      LastUL:                  "=+",
//      NilBetweenEveryRow:      true,
//  }
//
// Which will format tables like:
//
//  +==========+=============+========+==========+
//  |          | Bob         | Sue    | John     |
//  +==========+=============+========+==========+
//  | Hometown | San Antonio | Austin | New York |
//  +----------+-------------+--------+----------+
//  | Mother   | Bessie      | Mary   | Sarah    |
//  +----------+-------------+--------+----------+
//  | Father   | Rick        | Dan    | Mike     |
//  +==========+=============+========+==========+
func NewBoxedAlignOptions() *AlignOptions {
	return &AlignOptions{
		FirstDR:                 "+=",
		FirstLR:                 "=",
		FirstFirstDLR:           "=+=",
		FirstDLR:                "=+=",
		FirstDL:                 "=+",
		RowFirstUD:              "| ",
		RowSecondUD:             " | ",
		RowUD:                   " | ",
		RowLastUD:               " |",
		LeaveTrailingWhitespace: true,
		FirstNilFirstUDR:        "+=",
		FirstNilLR:              "=",
		FirstNilFirstUDLR:       "=+=",
		FirstNilUDLR:            "=+=",
		FirstNilLastUDL:         "=+",
		NilFirstUDR:             "+-",
		NilLR:                   "-",
		NilFirstUDLR:            "-+-",
		NilUDLR:                 "-+-",
		NilLastUDL:              "-+",
		LastUR:                  "+=",
		LastLR:                  "=",
		LastFirstULR:            "=+=",
		LastULR:                 "=+=",
		LastUL:                  "=+",
		NilBetweenEveryRow:      true,
	}
}

// NewUnicodeBoxedAlignOptions gives:
//  &AlignOptions{
//      FirstDR:                 "\u2554\u2550",
//      FirstLR:                 "\u2550",
//      FirstFirstDLR:           "\u2550\u2566\u2550",
//      FirstDLR:                "\u2550\u2564\u2550",
//      FirstDL:                 "\u2550\u2557",
//      RowFirstUD:              "\u2551 ",
//      RowSecondUD:             " \u2551 ",
//      RowUD:                   " \u2502 ",
//      RowLastUD:               " \u2551",
//      LeaveTrailingWhitespace: true,
//      FirstNilFirstUDR:        "\u2560\u2550",
//      FirstNilLR:              "\u2550",
//      FirstNilFirstUDLR:       "\u2550\u256c\u2550",
//      FirstNilUDLR:            "\u2550\u256a\u2550",
//      FirstNilLastUDL:         "\u2550\u2563",
//      NilFirstUDR:             "\u255f\u2500",
//      NilLR:                   "\u2500",
//      NilFirstUDLR:            "\u2500\u256b\u2500",
//      NilUDLR:                 "\u2500\u253c\u2500",
//      NilLastUDL:              "\u2500\u2562",
//      LastUR:                  "\u255a\u2550",
//      LastLR:                  "\u2550",
//      LastFirstULR:            "\u2550\u2569\u2550",
//      LastULR:                 "\u2550\u2567\u2550",
//      LastUL:                  "\u2550\u255d",
//      NilBetweenEveryRow:      true,
//  }
//
// Which will format tables like:
//
//  ╔══════════╦═════════════╤════════╤══════════╗
//  ║          ║ Bob         │ Sue    │ John     ║
//  ╠══════════╬═════════════╪════════╪══════════╣
//  ║ Hometown ║ San Antonio │ Austin │ New York ║
//  ╟──────────╫─────────────┼────────┼──────────╢
//  ║ Mother   ║ Bessie      │ Mary   │ Sarah    ║
//  ╟──────────╫─────────────┼────────┼──────────╢
//  ║ Father   ║ Rick        │ Dan    │ Mike     ║
//  ╚══════════╩═════════════╧════════╧══════════╝
func NewUnicodeBoxedAlignOptions() *AlignOptions {
	return &AlignOptions{
		FirstDR:                 "\u2554\u2550",
		FirstLR:                 "\u2550",
		FirstFirstDLR:           "\u2550\u2566\u2550",
		FirstDLR:                "\u2550\u2564\u2550",
		FirstDL:                 "\u2550\u2557",
		RowFirstUD:              "\u2551 ",
		RowSecondUD:             " \u2551 ",
		RowUD:                   " \u2502 ",
		RowLastUD:               " \u2551",
		LeaveTrailingWhitespace: true,
		FirstNilFirstUDR:        "\u2560\u2550",
		FirstNilLR:              "\u2550",
		FirstNilFirstUDLR:       "\u2550\u256c\u2550",
		FirstNilUDLR:            "\u2550\u256a\u2550",
		FirstNilLastUDL:         "\u2550\u2563",
		NilFirstUDR:             "\u255f\u2500",
		NilLR:                   "\u2500",
		NilFirstUDLR:            "\u2500\u256b\u2500",
		NilUDLR:                 "\u2500\u253c\u2500",
		NilLastUDL:              "\u2500\u2562",
		LastUR:                  "\u255a\u2550",
		LastLR:                  "\u2550",
		LastFirstULR:            "\u2550\u2569\u2550",
		LastULR:                 "\u2550\u2567\u2550",
		LastUL:                  "\u2550\u255d",
		NilBetweenEveryRow:      true,
	}
}

// Align will format a table according to options. If opts is nil,
// NewDefaultAlignOptions is used.
func Align(data [][]string, opts *AlignOptions) string {
	if data == nil || len(data) == 0 {
		return ""
	}
	if opts == nil {
		opts = NewDefaultAlignOptions()
	}
	newData := make([][]string, 0, len(data))
	for _, row := range data {
		if row == nil {
			if !opts.NilBetweenEveryRow {
				newData = append(newData, nil)
			}
			continue
		}
		if opts.Widths != nil {
			newRow := make([]string, 0, len(row))
			for col, cell := range row {
				if col >= len(opts.Widths) || opts.Widths[col] <= 0 {
					newRow = append(newRow, cell)
					continue
				}
				newRow = append(newRow, Wrap(cell, opts.Widths[col], "", ""))
			}
			row = newRow
		}
		work := make([][]string, 0, len(row))
		for _, cell := range row {
			cell = strings.Replace(cell, "\r\n", "\n", -1)
			work = append(work, strings.Split(cell, "\n"))
		}
		maxCells := 0
		for _, cells := range work {
			c := len(cells)
			if c > maxCells {
				maxCells = c
			}
		}
		newRows := make([][]string, 0)
		if opts.NilBetweenEveryRow && len(newData) != 0 {
			newData = append(newData, nil)
		}
		for c := 0; c < maxCells; c++ {
			newRow := make([]string, 0, len(work))
			for col := 0; col < len(work); col++ {
				if c < len(work[col]) {
					newRow = append(newRow, work[col][c])
				} else {
					newRow = append(newRow, "")
				}
			}
			newRows = append(newRows, newRow)
		}
		newData = append(newData, newRows...)
	}
	data = newData
	widths := make([]int, 0, len(data[0]))
	for _, row := range data {
		if row == nil {
			continue
		}
		for len(row) > len(widths) {
			widths = append(widths, len(row[len(widths)]))
		}
		for c, v := range row {
			if utf8.RuneCountInString(v) > widths[c] {
				widths[c] = utf8.RuneCountInString(v)
			}
		}
	}
	alignments := opts.Alignments
	if alignments == nil || len(alignments) < len(widths) {
		newal := append(make([]Alignment, 0, len(widths)), alignments...)
		for len(newal) < len(widths) {
			newal = append(newal, Left)
		}
		alignments = newal
	}
	est := utf8.RuneCountInString(opts.RowFirstUD)
	for _, w := range widths {
		est += w + utf8.RuneCountInString(opts.RowUD)
	}
	est += utf8.RuneCountInString(opts.RowLastUD) + 1
	est *= len(data)
	buf := bytes.NewBuffer(make([]byte, 0, est))
	if !AllEqual("", opts.FirstDR, opts.FirstFirstDLR, opts.FirstDLR, opts.FirstLR, opts.FirstDL) {
		buf.WriteString(opts.FirstDR)
		for col, width := range widths {
			if col == 1 {
				buf.WriteString(opts.FirstFirstDLR)
			} else if col != 0 {
				buf.WriteString(opts.FirstDLR)
			}
			for i := 0; i < width; i++ {
				buf.WriteString(opts.FirstLR)
			}
		}
		buf.WriteString(opts.FirstDL)
		buf.WriteByte('\n')
	}
	firstNil := true
	for _, row := range data {
		if row == nil {
			if firstNil {
				if !AllEqual("", opts.FirstNilFirstUDR, opts.FirstNilFirstUDLR, opts.FirstNilUDLR, opts.FirstNilLR, opts.FirstNilLastUDL) {
					buf.WriteString(opts.FirstNilFirstUDR)
					for col, width := range widths {
						if col == 1 {
							buf.WriteString(opts.FirstNilFirstUDLR)
						} else if col != 0 {
							buf.WriteString(opts.FirstNilUDLR)
						}
						for i := 0; i < width; i++ {
							buf.WriteString(opts.FirstNilLR)
						}
					}
					buf.WriteString(opts.FirstNilLastUDL)
				}
				firstNil = false
			} else {
				if !AllEqual("", opts.NilFirstUDR, opts.NilFirstUDLR, opts.NilUDLR, opts.NilLR, opts.NilLastUDL) {
					buf.WriteString(opts.NilFirstUDR)
					for col, width := range widths {
						if col == 1 {
							buf.WriteString(opts.NilFirstUDLR)
						} else if col != 0 {
							buf.WriteString(opts.NilUDLR)
						}
						for i := 0; i < width; i++ {
							buf.WriteString(opts.NilLR)
						}
					}
					buf.WriteString(opts.NilLastUDL)
				}
			}
			buf.WriteByte('\n')
			continue
		}
		buf.WriteString(opts.RowFirstUD)
		for c, v := range row {
			if c == 1 {
				buf.WriteString(opts.RowSecondUD)
			} else if c != 0 {
				buf.WriteString(opts.RowUD)
			}
			switch alignments[c] {
			case Right:
				for i := widths[c] - utf8.RuneCountInString(v); i > 0; i-- {
					buf.WriteRune(' ')
				}
				buf.WriteString(v)
			case Center:
				for i := (widths[c] - utf8.RuneCountInString(v)) / 2; i > 0; i-- {
					buf.WriteRune(' ')
				}
				buf.WriteString(v)
				if opts.LeaveTrailingWhitespace || c < len(row)-1 {
					for i := widths[c] - ((widths[c]-utf8.RuneCountInString(v))/2 + utf8.RuneCountInString(v)); i > 0; i-- {
						buf.WriteRune(' ')
					}
				}
			default:
				buf.WriteString(v)
				if opts.LeaveTrailingWhitespace || c < len(row)-1 {
					for i := widths[c] - utf8.RuneCountInString(v); i > 0; i-- {
						buf.WriteRune(' ')
					}
				}
			}
		}
		buf.WriteString(opts.RowLastUD)
		buf.WriteByte('\n')
	}
	if !AllEqual("", opts.LastUR, opts.LastFirstULR, opts.LastULR, opts.LastLR, opts.LastUL) {
		buf.WriteString(opts.LastUR)
		for col, width := range widths {
			if col == 1 {
				buf.WriteString(opts.LastFirstULR)
			} else if col != 0 {
				buf.WriteString(opts.LastULR)
			}
			for i := 0; i < width; i++ {
				buf.WriteString(opts.LastLR)
			}
		}
		buf.WriteString(opts.LastUL)
		buf.WriteByte('\n')
	}
	return buf.String()
}
