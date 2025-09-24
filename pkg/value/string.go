package value

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type strGetterFunc func(*String) (r.Element, error)
type strMethodFunc func(*String, []r.Element) (r.Element, error)

type String struct {
	value string
}

func NewString(value string) *String {
	return &String{value}
}

// String - display string value's string
func (s *String) String() string {
	return s.value
}

// GetValues - alias of s.String()
func (s *String) GetValue() string {
	return s.value
}

// GetProperty -
func (s *String) GetProperty(c *r.Context, name string) (r.Element, error) {
	strGetterMap := map[string]strGetterFunc{
		"长度":  strGetLength,
		"字数":  strGetLength,
		"文本":  strGetText,
		"字符组": strGetCharArray,
	}
	if fn, ok := strGetterMap[name]; ok {
		return fn(s)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (s *String) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (s *String) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	strMethodMap := map[string]strMethodFunc{
		"替换":     strExecReplace,
		"分隔":     strExecSplit,
		"匹配":     strExecMatch,
		"匹配开头":   strExecMatchStart,
		"匹配结尾":   strExecMatchEnd,
		"取样":     strExecSlice,
		"去除空格":   strExecStripWhitespaces,
		"转小写-英文": strExecToLowerCase,
		"转大写-英文": strExecToUpperCase,
		"拼接":     strExecJoin,
		"格式化":    strExecFormat,
		"转换数值":   strExecAtoi,
	}
	if fn, ok := strMethodMap[name]; ok {
		return fn(s, values)
	}
	return nil, zerr.MethodNotFound(name)
}

// ///// getters, setters and methods
// getters
func strGetLength(s *String) (r.Element, error) {
	l := utf8.RuneCountInString(s.value)
	return NewNumber(float64(l)), nil
}

func strGetText(s *String) (r.Element, error) {
	return NewString(s.value), nil
}

func strGetCharArray(s *String) (r.Element, error) {
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
func strExecReplace(s *String, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	oldItem := values[0].(*String).String()
	newItem := values[1].(*String).String()

	result := strings.ReplaceAll(s.value, oldItem, newItem)
	return NewString(result), nil
}

func strExecSplit(s *String, values []r.Element) (r.Element, error) {
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

func strExecMatch(s *String, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.Contains(s.value, substr)

	return NewBool(result), nil
}

func strExecMatchStart(s *String, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasPrefix(s.value, substr)

	return NewBool(result), nil
}

func strExecMatchEnd(s *String, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	substr := values[0].(*String).String()
	result := strings.HasSuffix(s.value, substr)

	return NewBool(result), nil
}

func strExecSlice(s *String, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number", "number"); err != nil {
		return nil, err
	}
	ss := s.GetValue()
	startIdx := int(values[0].(*Number).GetValue())
	endIdx := int(values[1].(*Number).GetValue())
	if startIdx < 0 {
		startIdx = len(ss) + startIdx + 1
	}

	if startIdx < 1 {
		return nil, ThrowException("文本的起始索引需从1开始！")
	}

	if endIdx > len(ss) {
		return nil, ThrowException("文本的结束索引不能超过其长度！")
	}

	if endIdx < 0 {
		endIdx = len(ss) + endIdx + 1
	}
	if startIdx > endIdx {
		return NewString(""), nil
	}

	subString := ss[startIdx-1 : endIdx]
	return NewString(subString), nil
}

func strExecStripWhitespaces(s *String, value []r.Element) (r.Element, error) {
	ss := strings.TrimSpace(s.value)
	return NewString(ss), nil
}

func strExecToLowerCase(s *String, value []r.Element) (r.Element, error) {
	ss := strings.ToLower(s.value)
	return NewString(ss), nil
}

func strExecToUpperCase(s *String, value []r.Element) (r.Element, error) {
	ss := strings.ToUpper(s.value)
	return NewString(ss), nil
}

func strExecJoin(s *String, values []r.Element) (r.Element, error) {
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
func strExecFormat(s *String, values []r.Element) (r.Element, error) {
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

func strExecAtoi(s *String, values []r.Element) (r.Element, error) {
	v := strings.Replace(s.value, "*^", "e", 1)
	v = strings.Replace(v, "*10^", "e", 1)
	s.value = v

	num, err := strconv.ParseFloat(s.value, 64)
	if err != nil {
		return nil, ThrowException("转成数值失败，文本可能并不符合合适的数值格式")
	}

	return NewNumber(num), nil
}
