// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hexya-erp/hexya/src/tools/logging"
)

var log logging.Logger

func init() {
	log = logging.GetLogger("strutils")
}

// SnakeCase convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func SnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// Title convert the given camelCase string to a title string.
// eg. MyHTMLData => My HTML Data
func Title(in string) string {

	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, ' ')
		}
		out = append(out, runes[i])
	}

	return string(out)
}

// GetDefaultString returns str if it is not an empty string or def otherwise
func GetDefaultString(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

// StartsAndEndsWith returns true if the given string starts with prefix
// and ends with suffix.
func StartsAndEndsWith(str, prefix, suffix string) bool {
	return strings.HasPrefix(str, prefix) && strings.HasSuffix(str, suffix)
}

// MarshalToJSONString marshals the given data to its JSON representation and
// returns it as a string. It panics in case of error.
func MarshalToJSONString(data interface{}) string {
	if _, ok := data.(string); !ok {
		domBytes, err := json.Marshal(data)
		if err != nil {
			log.Panic("Unable to marshal given data", "error", err, "data", data)
		}
		return string(domBytes)
	}
	return data.(string)
}

// HumanSize returns the given size (in bytes) in a human readable format
func HumanSize(size int64) string {
	units := []string{"bytes", "KB", "MB", "GB"}
	s, i := float64(size), 0
	for s >= 1024 && i < len(units)-1 {
		s /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", s, units[i])
}

// Substitute substitutes each occurrence of each key of mapping in str by the
// corresponding mapping value and returns the substituted string.
func Substitute(str string, mapping map[string]string) string {
	for key, val := range mapping {
		str = strings.Replace(str, key, val, -1)
	}
	return str
}

// DictToJSON sanitizes a python dict string representation to valid JSON.
func DictToJSON(dict string) string {
	dict = strings.Replace(dict, "'", "\"", -1)
	dict = strings.Replace(dict, "False", "false", -1)
	dict = strings.Replace(dict, "True", "true", -1)
	dict = strings.Replace(dict, "(", "[", -1)
	dict = strings.Replace(dict, ")", "]", -1)
	return dict
}

// MakeUnique returns an unique string in reference of the given pool
// its made of the base string plus a number if it exists within the pool
func MakeUnique(str string, pool []string) string {
	var nb int
	tested := str
	for tested == "" || IsIn(tested, pool...) {
		nb++
		tested = str + strconv.Itoa(nb)
	}
	return tested
}

// Reverse returns s reversed
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

// SplitEveryN retuns a slice containing str in multiple parts
// each part has its size to n (or less if last part)
func SplitEveryN(str string, n int) []string {
	var out []string
	var bSl []byte
	bStr := []byte(str)
	for i, b := range bStr {
		if (i)%n == 0 && i != 0 {
			out = append(out, string(bSl))
			bSl = nil
		}
		bSl = append(bSl, b)
	}
	out = append(out, string(bSl))
	return out
}

// SplitAtN returns the given strings separated as two string splited after the N-th character
func SplitAtN(str string, n int) (out1 string, out2 string) {
	if n > len(str) {
		out1 = str
		out2 = ""
	} else {
		out1 = str[:n]
		out2 = str[n:]
	}
	return
}

//ContainsOnly returns true if the given string only contains the characters given
func ContainsOnly(str string, lst ...byte) bool {
	for _, b := range []byte(str) {
		for _, l := range lst {
			if l != b {
				return false
			}
		}
	}
	return true
}

// NumberGrouping represents grouping values of a number as follows:
//  - it splits a number into groups of N, N being a value in the slice
//  - the last value represent the last group
//  - all values should be positive
//  - 0 means repetition of next int
//  - if the first value is not a 0, the grouping will end
//    e.g. :
//       3       -> 123456,789
//       0,3     -> 123,456,789
//       2,3     -> 1234,56,789
//       0,2,3   -> 12,34,56,789
type NumberGrouping []int

// FormatNumberStrWithGrouping formats a string (supposedly representing an integer)
// with number grouping defined by the given grouping slice. each group is separated with thSeparator
func FormatNumberStrWithGrouping(number string, grouping NumberGrouping, thSeparator string) string {
	number = Reverse(number)
	thSeparator = Reverse(thSeparator) //reverse strings
	var str string

	var revGrouping NumberGrouping
	for i := range grouping { //reverse grouping
		revGrouping = append(revGrouping, grouping[len(grouping)-i-1])
	}

	var out string
	last := 9999
	for _, n := range revGrouping { // use all grouping numbers
		if n == 0 {
			n = last
		}
		str, number = SplitAtN(number, n)
		if str != "" {
			out = out + thSeparator + str
		}
		if number == "" {
			return Reverse(strings.TrimPrefix(out, thSeparator))
		}
		last = n
	}
	// all grouping values exhausted, continue with N until number is empty
	for number != "" {
		n := 9999
		if grouping[0] == 0 {
			n = last
		}
		str, number = SplitAtN(number, n)
		if str != "" {
			out = out + thSeparator + str
		}
	}
	return Reverse(strings.TrimPrefix(out, thSeparator))
}

// FormatMonetary formats a float into a monetary string
// eg. FormatMonetary(3.14159, 2, 0, ",", "$", true) => "$ 3,14"
// Params:
//	value: the float value to be formated
//	digits: the ammount of digits written after the decimal point
//	grouping: (See type NumberGrouping for more on this)
//	separator: the character used as the decimal separator
//  thSeparator: the character used as grouping separator
//	symbol: the currency symbol
//	symPosLeft: whether or not the symbol shall be put before the value
func FormatMonetary(value float64, digits, grouping NumberGrouping, separator, thSeparator, symbol string, symToLeft bool) string {
	fmtStr := fmt.Sprintf("%%.%df", digits)
	str := fmt.Sprintf(fmtStr, value)
	strSpl := strings.Split(str, ".")
	str = FormatNumberStrWithGrouping(strSpl[0], grouping, thSeparator)
	if len(strSpl) > 1 {
		str = strings.Join([]string{str, strSpl[1]}, separator)
	}
	if symbol != "" {
		if symToLeft {
			str = fmt.Sprintf("%s %s", symbol, str)
		} else {
			str = fmt.Sprintf("%s %s", str, symbol)
		}
	}
	return str
}

// IsIn returns true if the given str is the same as one of the strings given in lst
func IsIn(str string, lst ...string) bool {
	for _, l := range lst {
		if str == l {
			return true
		}
	}
	return false
}
