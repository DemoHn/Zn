package val

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// 将对应字符表示替换成对应的实际字符
var specialCharMap = map[string]string{
	"{/CR}":  "\r",
	"{/r}":   "\r",
	"{/LF}":  "\n",
	"{/n}":   "\n",
	"{/TAB}": "\t",
	"{/t}":   "\t",
	"{/SO}":  "/",
	"{/s}":   "/",
}

// String - represents for Zn's 文本型
type String struct {
	value string
}

// NewString - new string ctx.Value Object from raw string
func NewString(value string) *String {
	v := replaceSpecialChars(value)
	return &String{v}
}

// String - display string value's string
func (s *String) String() string {
	return s.value
}

// replaceSpecialChars: A{/CR}B -> A\nB
func replaceSpecialChars(s string) string {
	re := regexp.MustCompile(`\{\/(\+?[0-9a-zA-Z]+?)\}`)
	reUnicode := regexp.MustCompile(`\{\/\+([0-9a-fA-F]+)\}`)

	return re.ReplaceAllStringFunc(s, func(ss string) string {
		// #1. check if matched string is a special char
		if res, ok := specialCharMap[ss]; ok {
			return res
		}

		// #2. check if matched string is a unicode representation
		if matches := reUnicode.FindStringSubmatch(ss); len(matches) > 1 {
			hexData, _ := strconv.ParseInt(matches[1], 16, 32)
			hexData32 := int32(hexData)
			return string([]rune{hexData32})
		}

		// #3. otherwise, return directly
		return ss
	})
}

// GetProperty -
func (s *String) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "长度":
		l := utf8.RuneCountInString(s.value)
		return NewDecimalFromInt(l, 0), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (s *String) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (s *String) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	// S 之（替换：<旧项>, <新项>）
	case "替换":
		if err := ValidateExactParams(values, "string", "string"); err != nil {
			return nil, err
		}
		oldItem := values[0].(*String).String()
		newItem := values[1].(*String).String()

		result := strings.ReplaceAll(s.value, oldItem, newItem)
		return NewString(result), nil
	// S 之（分隔：<分隔符>） -> 【A，B，C，...】
	case "分隔":
		if err := ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}
		sep := values[0].(*String).String()
		resultStrs := strings.Split(s.value, sep)

		result := NewArray([]ctx.Value{})
		for _, v := range resultStrs {
			result.AppendValue(NewString(v))
		}

		return result, nil
	case "匹配":
		if err := ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}
		substr := values[0].(*String).String()
		result := strings.Contains(s.value, substr)

		return NewBool(result), nil
	case "匹配开头":
		if err := ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}
		substr := values[0].(*String).String()
		result := strings.HasPrefix(s.value, substr)

		return NewBool(result), nil
	case "匹配结尾":
		if err := ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}
		substr := values[0].(*String).String()
		result := strings.HasSuffix(s.value, substr)

		return NewBool(result), nil
	// 「xxx」 + 「yyy」 + 「zzz」
	case "拼接":
		if err := ValidateAllParams(values, "string"); err != nil {
			return nil, err
		}
		result := s.value
		for _, v := range values {
			result += v.(*String).String()
		}

		return NewString(result), nil
	// 「xxx {#1} yyy {#2}」.Format(「A」,「B」)
	case "格式拼接":
		if err := ValidateAllParams(values, "string"); err != nil {
			return nil, err
		}

		replacerArgs := []string{}
		for idx, v := range values {
			format := fmt.Sprintf("{#%d}", idx+1)
			value := v.(*String).String()

			replacerArgs = append(replacerArgs, format, value)
		}

		r := strings.NewReplacer(replacerArgs...)
		// replace {#1} with value1, {#2} with value2, ...
		result := r.Replace(s.value)

		return NewString(result), nil
	}
	return nil, error.MethodNotFound(name)
}
