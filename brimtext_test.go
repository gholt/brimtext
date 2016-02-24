package brimtext

import (
	"fmt"
	"sort"
	"testing"
)

func TestOrdinal(t *testing.T) {
	for i, v := range map[int]string{
		0: "th", 1: "st", 2: "nd", 3: "rd", 4: "th",
		10: "th", 11: "th", 12: "th", 13: "th", 14: "th",
		20: "th", 21: "st", 22: "nd", 23: "rd", 24: "th",
		100: "th", 101: "st", 102: "nd", 103: "rd", 104: "th",
		110: "th", 111: "th", 112: "th", 113: "th", 114: "th",
		120: "th", 121: "st", 122: "nd", 123: "rd", 124: "th",
	} {
		if OrdinalSuffix(i) != v {
			t.Errorf("%#v != %#v", i, v)
		}
	}
}

func TestThousandsSep(t *testing.T) {
	for i, x := range map[int64]string{
		-1000:               "-1,000",
		-1:                  "-1",
		0:                   "0",
		999:                 "999",
		1000:                "1,000",
		100000:              "100,000",
		1000000:             "1,000,000",
		1000000000000000000: "1,000,000,000,000,000,000",
	} {
		o := ThousandsSep(i, ",")
		if o != x {
			t.Errorf("ThousandsSep(%#v) %#v != %#v", i, o, x)
		}
	}
}

func TestThousandsSepU(t *testing.T) {
	for i, x := range map[uint64]string{
		0:                   "0",
		999:                 "999",
		1000:                "1,000",
		100000:              "100,000",
		1000000:             "1,000,000",
		1000000000000000000: "1,000,000,000,000,000,000",
	} {
		o := ThousandsSepU(i, ",")
		if o != x {
			t.Errorf("ThousandsSepU(%#v) %#v != %#v", i, o, x)
		}
	}
}

func TestHumanSize(t *testing.T) {
	for i, v := range map[int64]string{
		0:                   "0",
		1:                   "1",
		512:                 "512",
		1023:                "1023",
		1024:                "1K",
		1535:                "1K",
		1536:                "2K",
		1048576:             "1M",
		1073741824:          "1G",
		1099511627776:       "1T",
		1125899906842624:    "1P",
		1152921504606846976: "1E",
	} {
		o := HumanSize(i, "")
		if o != v {
			t.Errorf("HumanSize(%#v) %#v != %#v", i, o, v)
		}
	}
	out := HumanSize(123, "b")
	exp := "123b"
	if out != exp {
		t.Errorf("%#v != %#v", out, exp)
	}
}

func TestSentence(t *testing.T) {
	for in, exp := range map[string]string{
		"":          "",
		"testing":   "Testing.",
		"'testing'": "'testing'.",
		"Testing.":  "Testing.",
		"testing.":  "Testing.",
	} {
		out := Sentence(in)
		if out != exp {
			t.Errorf("Sentence(%#v) %#v != %#v", in, out, exp)
		}
	}
}

func TestStringSliceToLowerSort(t *testing.T) {
	out := []string{"DEF", "abc"}
	sort.Sort(StringSliceToLowerSort(out))
	exp := []string{"abc", "DEF"}
	for i := 0; i < len(out); i++ {
		if out[i] != exp[i] {
			t.Fatalf("StringSliceToLowerSort fail at index %d %#v != %#v", i, out[i], exp[i])
			return
		}
	}
	out = []string{"DEF", "abc"}
	sort.Strings(out)
	exp = []string{"DEF", "abc"}
	for i := 0; i < len(out); i++ {
		if out[i] != exp[i] {
			t.Fatalf("sort.Strings sort fail at index %d %#v != %#v", i, out[i], exp[i])
			return
		}
	}
}

func TestWrap(t *testing.T) {
	in := ""
	out := Wrap(in, 79, "", "")
	exp := ""
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = "Just a test sentence."
	out = Wrap(in, 10, "", "")
	exp = `Just a
test
sentence.`
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = "Just   a   test   sentence."
	out = Wrap(in, 10, "", "")
	exp = `Just a
test
sentence.`
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = fmt.Sprintf("Just a %stest%s sentence.", string(ANSIEscape.Bold), string(ANSIEscape.Reset))
	out = Wrap(in, 10, "", "")
	exp = fmt.Sprintf(`Just a
%stest%s
sentence.`, string(ANSIEscape.Bold), string(ANSIEscape.Reset))
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = "Just a test sentence."
	out = Wrap(in, 10, "1234", "5678")
	exp = `1234Just a
5678test
5678sentence.`
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = `Just a test sentence. With
a follow on sentence.

And a separate paragraph.`
	out = Wrap(in, 10, "", "")
	exp = `Just a
test
sentence.
With a
follow on
sentence.

And a
separate
paragraph.`
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
	in = `Just a test sentence.  With     
          a follow           on sentence.

                And a separate  paragraph.       `
	out = Wrap(in, 10, "", "")
	exp = `Just a
test
sentence.
With a
follow on
sentence.

And a
separate
paragraph.`
	if out != exp {
		t.Errorf("Wrap(%#v) %#v != %#v", in, out, exp)
	}
}

func TestAllEqual(t *testing.T) {
	if !AllEqual() {
		t.Fatal("")
	}
	if !AllEqual([]string{}...) {
		t.Fatal("")
	}
	if !AllEqual("bob") {
		t.Fatal("")
	}
	if !AllEqual("bob", "bob") {
		t.Fatal("")
	}
	if !AllEqual("bob", "bob", "bob") {
		t.Fatal("")
	}
	if !AllEqual([]string{"bob", "bob", "bob"}...) {
		t.Fatal("")
	}
	if AllEqual("bob", "sue") {
		t.Fatal("")
	}
	if AllEqual("bob", "bob", "sue") {
		t.Fatal("")
	}
	if AllEqual([]string{"bob", "bob", "sue"}...) {
		t.Fatal("")
	}
}
