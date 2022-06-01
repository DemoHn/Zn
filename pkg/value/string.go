package value

import (
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type strGetterFunc func(*String, *r.Context) (r.Value, error)
type strMethodFunc func(*String, *r.Context, []r.Value) (r.Value, error)

type String struct {
	value string
}

// replace special chars from {/xx} placeholders
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

func NewString(value string) *String {
	v := replaceSpecialChars(value)
	return  &String{v}
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
func (s *String) GetProperty(c *r.Context, name string) (r.Value, error) {
	strGetterMap := map[string]strGetterFunc{
		"长度": strGetLength,
		"字数": strGetLength,
		"文本": strGetText,
		"字符组": strGetCharArray,
	}
	if fn, ok := strGetterMap[name]; ok {
		return fn(s, c)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (s *String) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (s *String) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	strMethodMap := map[string]strMethodFunc{
		"替换": strExecReplace,
		"分隔": strExecSplit,
		"匹配": strExecMatch,
		"匹配开头": strExecMatchStart,
		"匹配结尾": strExecMatchEnd,
		"拼接": strExecJoin,
		"格式化": strExecFormat,
	}
	if fn, ok := strMethodMap[name]; ok {
		return fn(s, c, values)
	}
	return nil, zerr.MethodNotFound(name)
}

/////// getters, setters and methods
// getters
func strGetLength(s *String, c *r.Context) (r.Value, error) {
	l := utf8.RuneCountInString(s.value)
	return NewNumber(float64(l)), nil
}

func strGetText(s *String, c *r.Context) (r.Value, error) {
	return NewString(s.value), nil
}

func strGetCharArray(s *String, c *r.Context) (r.Value, error) {
	charArr := NewArray([]r.Value{})
	v := s.value

	for len(v) > 0 {
		r, size := utf8.DecodeRuneInString(v)
		charArr.AppendValue(NewString(string(r)))

		v = v[size:]
	}
	return charArr, nil
}

// methods
func strExecReplace(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	oldItem := values[0].(*String).String()
	newItem := values[1].(*String).String()

	result := strings.ReplaceAll(s.value, oldItem, newItem)
	return NewString(result), nil
}

func strExecSplit(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	sep := values[0].(*String).String()
	resultStrs := strings.Split(s.value, sep)

	result := NewArray([]r.Value{})
	for _, v := range resultStrs {
		result.AppendValue(NewString(v))
	}

	return result, nil
}

func strExecMatch(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.Contains(s.value, substr)

	return NewBool(result), nil
}

func strExecMatchStart(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasPrefix(s.value, substr)

	return NewBool(result), nil
}

func strExecMatchEnd(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasSuffix(s.value, substr)

	return NewBool(result), nil
}

func strExecJoin(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateAllParams(values, "string"); err != nil {
		return nil, err
	}
	result := s.value
	for _, v := range values {
		result += v.(*String).String()
	}

	return NewString(result), nil
}

// format string, like python's "<format_string:%s>" % (str)
func strExecFormat(s *String, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateAllParams(values, "string"); err != nil {
		return nil, err
	}

	var replacerArgs []string
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

