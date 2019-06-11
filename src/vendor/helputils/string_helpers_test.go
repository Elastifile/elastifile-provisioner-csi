package helputils

import (
	"fmt"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"optional"
)

type some struct {
	A optional.String
	B optional.Int
	C optional.Float64
	D thing
}

type thing struct {
	E int
}

type bar struct{}

func (bar) String() string { return "Its my bar" }

func TestRepeatWithSeparator(t *testing.T) {
	//RepeatWithSeparator("bar", 3, "-") --> result: "bar-bar-bar"
	result := RepeatWithSeparator("bar", 3, "-")
	if result != "bar-bar-bar" {
		fmt.Println("Failed")
	}
}

func TestStructToStringsDefaults(t *testing.T) {
	test := some{}
	result := strings.Join(StructToStrings(test), " ")
	if result != "N/A N/A N/A" {
		t.Fatalf("unexpected result: '%v', for test: '%s'", result, test)
	}
	t.Log("[v] " + result)
}

func TestStructToStrings(t *testing.T) {
	a := "a"
	b := 10
	c := 1.1
	data := some{
		A: &a,
		B: &b,
		C: &c,
		D: thing{
			E: 2,
		},
	}

	result := strings.Join(StructToStrings(data), " ")
	if result != "a 10 1.10 {2}" {
		t.Fatalf("unexpected result: '%s'", result)
	}
	t.Logf("[v] %v", result)

	result = strings.Join(StructToNamedStrings(data), " ")
	if result != "A:a B:10 C:1.10 D:{E:2}" {
		t.Fatalf("unexpected result: '%s'", result)
	}
	t.Logf("[v] %v", result)
}

func TestContainsStrings(t *testing.T) {
	strs := []string{"a", "bb", "ccc"}
	ints := ToStrings([]int{1, 2, 3})
	complex := ToStrings(&bar{})

	for _, items := range [][]string{strs, ints, complex} {
		for _, item := range items {
			t.Log(fmt.Sprintf("Checking if %v contains %v.", items, item))
			if !ContainsStr(items, item) {
				t.Fatal("Failed")
			}
		}
	}
}

func TestFindIP(t *testing.T) {
	in := `
    Waiting for IP..
    Received IP address: 10.11.18.4...
    Done.
    `
	ips := FindIP(in, 100)
	t.Log(ips)
	if len(ips) == 0 {
		t.Fatalf("Failed to find ip pattern in:\n \"%v\"\n", in)
	}

	expected := "10.11.18.4"
	if ips[0] != expected {
		t.Fatalf("Expected ip: %v doesn't match ip: %v\n", expected, ips[0])
	}
}

func TestKeyValuesStringer(t *testing.T) {
	t.Log(KeyValuesStringer([]string{
		"a", "1",
		"b", "X",
		"c", "-1",
		"d", "--1",
	}))
}

func TestStringSliceToInterfaces(t *testing.T) {
	t.Log(StringSliceToInterfaces(
		"a", "1",
		"b", "X",
		"c", "-1",
		"d", "--1",
	))
}

func TestNumericExtractFromString(t *testing.T) {
	for test, expect := range map[string]int{
		"vc7c":    7,
		"vc8":     8,
		"32":      32,
		"x1y":     1,
		"ab":      0,
		"":        0,
		"1 2 3 4": 1234,
	} {
		result := NumericExtractFromString(test)
		if result != expect {
			t.Fatalf("Unexpected result:", result, "for", test, ", expected:", expect)
		}
		t.Log("[v] test:", test, "-> result:", result)
	}
}

func TestEqualsAnyStr(t *testing.T) {
	for test, expect := range map[string]bool{
		".=.|..":  true,
		"..=.|..": true,
		"=.|..":   false,
		"x=z|y|x": true,
		"xyz=":    false,
		"xyz=||":  false,
	} {
		testSplit := strings.Split(test, "=")
		str, opts := testSplit[0], testSplit[1]
		result := EqualsAnyStr(str, strings.Split(opts, "|")...)
		if result != expect {
			t.Fatalf("Unexpected result:", result, "for", str, opts, ", expected:", expect)
		}
		t.Log("[v] test:", str, "in", opts, "\t->", result)
	}
}

func TestYamlToOneLine(t *testing.T) {
	testStr := `
A: 1
B: 2
C: 
    D: 4
    E: 5
    F:
        G: 6
        H: 7
I: 8
J:
    K:
    	- 1
    	- 2
    	- 3
    L: 1
`
	var test struct {
		A int
		B int
		C struct {
			D int
			E int
			F struct {
				G int
				H int
			}
		}
		I int
		J struct {
			K []byte
			L int
		}
	}
	result := YamlToOneLine(testStr, 4)
	err := yaml.Unmarshal([]byte(result), &test)
	if err != nil {
		t.Fatal(err, "\n", result)
	}
	expect := "{A: 1, B: 2, C:  {D: 4, E: 5, F: {G: 6, H: 7, }, }, I: 8, J: {K: [1, 2, 3, ], L: 1,}, }"
	if result == expect {
		t.Logf(result)
	} else {
		t.Fatalf("unexpected result: '%s' for input: %s", result, testStr)
	}
}

func TestZipStrings(t *testing.T) {
	for test, expect := range map[string]string{
		"a,b,c:1,2,3": "a,1,b,2,c,3",
		"a,b:1,2,3":   "a,1,b,2,,3",
		"a,b,c:1,2":   "a,1,b,2,c,",
		":1,2,3":      ",1,,2,,3",
		"a,b,c:":      "a,,b,,c,",
	} {
		testSplit := strings.Split(test, ":")
		left := strings.Split(testSplit[0], ",")
		right := strings.Split(testSplit[1], ",")
		zip := ZipStrings(left, right)
		if strings.Join(zip, ",") != expect {
			t.Fatalf("unexpected result: %v, from %v + %v (expected: %v)",
				zip, left, right, strings.Split(expect, ","))
		} else {
			t.Logf("%v + %v => %v", left, right, zip)
		}
	}
}
