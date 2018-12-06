package helputils

import (
	"fmt"
	"math"
	"math/rand"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-errors/errors"
)

// RepeatWithSeparator is much like strings.Repeat but it also supports separator string.
//
// Examples:
// RepeatWithSeparator("foo", 2, " ") --> result: "foo foo"
// RepeatWithSeparator("bar", 3, "-") --> result: "bar-bar-bar"
func RepeatWithSeparator(s string, count int, sep string) string {
	all := []string{}
	for ; count > 0; count-- {
		all = append(all, s)
	}
	return strings.Join(all, sep)
}

// StructToStrings converts struct values to slice of strings.
// Unassigned fields will be assigned with "N/A".
//
// Example:
// type foo struct{
//     Bar int
// }
//
// StructToStrings(foo{Bar: 10}) --> []string{"10"}
// StructToStrings(foo{})        --> []string{"N/A"}
func StructToStrings(data interface{}) (line []string) {
	return structToStrings(data, false, "N/A")
}

func StructToNamedStrings(data interface{}) (line []string) {
	return structToStrings(data, true, "")
}

func structToStrings(data interface{}, names bool, nilValue string) (line []string) {
	if !IsReflectingStruct(data) {
		panic("must specify struct or ptr to struct")
	}

	StructExportedFieldIter(data, func(field *structs.Field) {
		var nameStr, valStr string
		if names {
			nameStr = field.Name() + ":"
		}
		if field.Kind() == reflect.Ptr && field.IsZero() { // nil values converted to nilValue
			if nilValue != "" {
				line = append(line, nameStr+nilValue)
			}
		} else {
			fieldValue := reflect.ValueOf(field.Value())
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			switch fieldValue.Kind() {
			case reflect.Float64:
				valStr = fmt.Sprintf("%.2f", fieldValue.Float())
			case reflect.Int:
				valStr = fmt.Sprintf("%v", fieldValue.Int())
			case reflect.String:
				valStr = fmt.Sprintf("%v", fieldValue.String())
			case reflect.Struct:
				valStr = "{" + strings.Join(structToStrings(fieldValue.Interface(), names, nilValue), " ") + "}"
			default:
				valStr = "ERROR"
			}
			line = append(line, nameStr+valStr)
		}
	})
	return line
}

// ContainsStr returns 'true' if target is matching one of given items
// compared by their default string representation, else returns 'false'.
// Can be combined with 'ToStrings' util. see unit tests for usage examples.
func ContainsStr(items []string, target ...string) bool {
	return FindStr(items, target...) > -1
}

func FindStr(items []string, target ...string) int {
	for i, item := range items {
		for _, t := range target {
			if item == t {
				return i
			}
		}
	}
	return -1
}

func EqualsAnyStr(str string, target ...string) bool {
	return ContainsStr([]string{str}, target...)
}

// MapStr takes a function from string to string and applies it to every
// item in `items' slice.  Returns the results of this application.
func MapStr(fn func(string) string, items []string) []string {
	var result []string
	for _, item := range items {
		result = append(result, fn(item))
	}
	return result
}

// FilterStr returns string slice built from specified strings,
// containing only items matching with specified filter func.
func FilterStr(items []string, fn func(string) bool) []string {
	var result []string
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// NonEmptyStr returns string slice from specified strings,
// containing only the non-empty strings
// (empty: space-trimming results with emtpy string)
func NonEmptyStr(items []string) []string {
	return FilterStr(items, IsNonEmptyStr)
}

// IsNonEmptyStr return wether specified string, trimmed from all space chars,
// is with length greater than zero.
func IsNonEmptyStr(valStr string) bool {
	return len(strings.TrimSpace(valStr)) > 0
}

// PrefixStr returns a function prepending `pref' to its argument.
func PrefixStr(pref string) func(string) string {
	return func(suffix string) string {
		return path.Join(pref, suffix)
	}
}

// ToStrings convert slices of different types to slice of strings.
// Any other type (e.g. builtin, struct, map..etc) will be
// converted to its default string representation (e.g. "%s").
func ToStrings(arr interface{}) []string {
	res := []string{}
	kind := reflect.TypeOf(arr).Kind()
	switch {
	case kind == reflect.Slice:
		// Found slice, start converting to string
		raw := reflect.ValueOf(arr)      // slice raw data
		slice := raw.Slice(0, raw.Len()) // use as a slice
		for i := 0; i < slice.Len(); i++ {
			res = append(res, fmt.Sprintf("%v", slice.Index(i)))
		}
	default:
		return []string{fmt.Sprintf("%s", reflect.ValueOf(arr))}
	}
	return res
}

func NumericExtractFromString(valStr string) int {
	re := regexp.MustCompile("[^0-9]")
	vcIdStr := re.ReplaceAllString(valStr, "")
	vcId, _ := strconv.Atoi(vcIdStr)
	return vcId
}

func KeyValuesStringer(keyValues []string) string {
	if len(keyValues)%2 != 0 {
		panic("dev fault: key-value slice must be with even length")
	}
	result := make([]string, len(keyValues)/2)
	for i, item := range keyValues {
		if i%2 == 0 {
			result[i/2] = item + ":"
		} else {
			if _, err := strconv.ParseInt(item, 10, 64); err == nil {
				result[(i-1)/2] += item
			} else {
				result[(i-1)/2] += "\"" + item + "\""
			}
		}
	}
	return strings.Join(result, " ")
}

// FindIP will return the matching IPv4 pattern in given string.
// 'count' will determine the number of occurences to return.
func FindIP(raw string, count int) []string {
	segment := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := segment + "\\." + segment + "\\." + segment + "\\." + segment

	regEx := regexp.MustCompile(regexPattern)
	ips := regEx.FindAllString(raw, count)

	return ips
}

func StringSliceToInterfaces(strs ...string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, valStr := range strs {
		result[i] = valStr
	}
	return result
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(strLen int) string {
	ret := make([]byte, strLen)
	for i := range ret {
		ret[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(ret)
}

func YamlToOneLine(body string, identSize int) string {
	indent := 0

	lines := FilterStr(strings.Split(body, "\n"),
		func(line string) bool { return len(line) > 0 },
	)

	var list bool

	lines = MapStr(
		func(line string) string {
			if len(line) > 0 {
				t := 0
				for ; line[t] == ' '; t++ {
				}
				if t%identSize != 0 {
					panic(fmt.Sprintf("expecting indent of %d:\n%s", identSize, body))
				}

				line = line[t:]
				for i := indent; i < t; i += identSize {
					line = "{" + line
				}
				for i := indent; i > t; i -= identSize {
					line = "}, " + line
				}
				indent = t

				if strings.Contains(line, ":") {
					val := strings.TrimSpace(strings.Split(line, ":")[1])
					if len(val) > 0 {
						line = line + ","
					}
					if list {
						line = "], " + line
						list = false
					}
				} else if strings.Contains(line, "-") {
					val := strings.TrimSpace(strings.Split(line, "-")[1])
					if len(val) > 0 {
						if !list {
							line = "[" + val + ","
							list = true
						} else {
							line = val + ","
						}
					}
				} else {
					panic("malformatted YAML line: " + line)
				}
			}
			return line
		},
		lines,
	)
	result := "{" + strings.Join(lines, " ") + "}"

	for i := indent; i > 0; i -= identSize {
		result = result + ", }"
	}

	return result
}

func ZipStrings(left []string, right []string) []string {
	max := int(math.Max(float64(len(left)), float64(len(right))))
	result := make([]string, 2*max)
	for i := 0; i < max; i++ {
		if i < len(left) {
			result[i*2] = left[i]
		}
		if i < len(right) {
			result[i*2+1] = right[i]
		}
	}
	return result
}

// GetValueByFieldName splits the string into fields, and returns the field following the specified name
func GetValueByFieldName(fieldName string, string2Parse string) (value string, err error) {
	valueIdx := -1
	fields := strings.Fields(string2Parse)
	for i, field := range fields {
		if field == fieldName {
			valueIdx = i + 1
		}
	}

	if valueIdx == -1 {
		err = errors.Errorf("Couldn't find field '%v' in '%v' (fields: %+v)", fieldName, string2Parse, fields)
		return
	}

	if valueIdx >= len(fields) {
		err = errors.Errorf("Couldn't find the value of field '%v' in '%v' (not enough fields: %+v)", fieldName, string2Parse, fields)
		return
	}

	value = fields[valueIdx]
	return
}