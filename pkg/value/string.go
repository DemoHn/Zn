package value

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	r "github.com/DemoHn/Zn/pkg/runtime"
)

type String struct {
	value string
	*r.ElementModel
}

// replace special chars from {/xx} placeholders
var specialCharMap = map[string]string{
	"`CR`":   "\r",
	"`LF`":   "\n",
	"`CRLF`": "\r\n",
	"`TAB`":  "\t",
	"`BK`":   "`",
}

func NewString(value string) *String {
	v := replaceSpecialChars(value)
	s := &String{v, r.NewElementModel()}

	// init getters & setters & methods
	s.RegisterGetter("长度", s.strGetLength)
	s.RegisterGetter("字数", s.strGetLength)
	s.RegisterGetter("文本", s.strGetText)
	s.RegisterGetter("字符组", s.strGetCharArray)

	s.RegisterMethod("替换", s.strExecReplace)
	s.RegisterMethod("分隔", s.strExecSplit)
	s.RegisterMethod("匹配", s.strExecMatch)
	s.RegisterMethod("匹配开头", s.strExecMatchStart)
	s.RegisterMethod("匹配结尾", s.strExecMatchEnd)
	s.RegisterMethod("拼接", s.strExecJoin)
	s.RegisterMethod("格式化", s.strExecFormat)

	return s
}

// String - display string value's string
func (s *String) String() string {
	return s.value
}

// replaceSpecialChars: A{/CR}B -> A\nB
func replaceSpecialChars(s string) string {
	re := regexp.MustCompile("`(CR|LF|CRLF|TAB|BK|U\\+[0-9A-Fa-f]{1,8})`")

	return re.ReplaceAllStringFunc(s, func(ss string) string {
		// #1. check if matched string is a special char (exclude U+xxxx)
		if res, ok := specialCharMap[ss]; ok {
			return res
		}

		// #2. match U+xxxx
		hexData, _ := strconv.ParseInt(ss[2:len(ss)-1], 16, 32)
		hexData32 := int32(hexData)
		return string([]rune{hexData32})
	})
}

/////// getters, setters and methods
// getters
func (s *String) strGetLength(c *r.Context) (r.Element, error) {
	l := utf8.RuneCountInString(s.value)
	return NewNumber(float64(l)), nil
}

func (s *String) strGetText(c *r.Context) (r.Element, error) {
	return NewString(s.value), nil
}

func (s *String) strGetCharArray(c *r.Context) (r.Element, error) {
	charArr := NewArray([]r.Element{})
	v := s.value

	for len(v) > 0 {
		r, size := utf8.DecodeRuneInString(v)
		charArr.AppendValue(NewString(string(r)))

		v = v[size:]
	}
	return charArr, nil
}

// methods
func (s *String) strExecReplace(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	oldItem := values[0].(*String).String()
	newItem := values[1].(*String).String()

	result := strings.ReplaceAll(s.value, oldItem, newItem)
	return NewString(result), nil
}

func (s *String) strExecSplit(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	sep := values[0].(*String).String()
	resultStrs := strings.Split(s.value, sep)

	result := NewArray([]r.Element{})
	for _, v := range resultStrs {
		result.AppendValue(NewString(v))
	}

	return result, nil
}

func (s *String) strExecMatch(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.Contains(s.value, substr)

	return NewBool(result), nil
}

func (s *String) strExecMatchStart(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasPrefix(s.value, substr)

	return NewBool(result), nil
}

func (s *String) strExecMatchEnd(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasSuffix(s.value, substr)

	return NewBool(result), nil
}

func (s *String) strExecJoin(c *r.Context, values []r.Element) (r.Element, error) {
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
func (s *String) strExecFormat(c *r.Context, values []r.Element) (r.Element, error) {
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
