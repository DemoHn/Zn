package exec

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/error"
)

const (
	maxDigitCount      = 18 // XXXXXXXX.XXXXXXXXXX
	maxLeadDecimalZero = 6  // 0.XXXXXX1234
	maxSciDigitCount   = 8  // 2.XXXXXXXX *10^ N
)

// Decimal - decimal number 「数值」型
type Decimal struct {
	// decimal internal properties
	co  *big.Int
	exp int
}

// NewDecimal -
func NewDecimal(value string) (*Decimal, *error.Error) {
	var decimal = &Decimal{
		exp: 0,
		co:  big.NewInt(0),
	}

	err := decimal.setValue(value)
	return decimal, err
}

// NewDecimalFromInt -
func NewDecimalFromInt(value int, exp int) *Decimal {
	return &Decimal{
		exp: exp,
		co:  big.NewInt(int64(value)),
	}
}

// String - show decimal display string
func (zd *Decimal) String() string {
	var sflag = ""
	if zd.co.Sign() < 0 {
		sflag = "-"
	}
	var txt = new(big.Int).Abs(zd.co).String()

	digitCount := len(txt)
	pointPos := zd.exp + digitCount

	// CASE I: no decimal point
	if digitCount <= pointPos && pointPos <= maxDigitCount {
		// subcase: add tail zeros
		if zd.exp > 0 {
			var zeros = strings.Repeat("0", zd.exp)
			return fmt.Sprintf("%s%s%s", sflag, txt, zeros)
		}
		return fmt.Sprintf("%s%s", sflag, txt)
	}
	// CASE II: with decimal point
	if pointPos <= maxDigitCount && pointPos < digitCount && pointPos > 0 && digitCount <= maxDigitCount {
		return fmt.Sprintf("%s%s.%s", sflag, txt[:pointPos], txt[pointPos:])
	}
	// CASE III: lead 0. 0s
	if pointPos <= 0 && pointPos > -maxLeadDecimalZero {
		var zeros = strings.Repeat("0", -pointPos)
		return fmt.Sprintf("%s0.%s%s", sflag, zeros, txt)
	}
	// CASE IV: sci format (1.23*10^-5)
	if digitCount > maxSciDigitCount {
		return fmt.Sprintf("%s%s.%s*10^%d", sflag, txt[0:1], txt[1:maxSciDigitCount+1], pointPos-1)
	} else if digitCount > 1 {
		return fmt.Sprintf("%s%s.%s*10^%d", sflag, txt[0:1], txt[1:], pointPos-1)
	}
	return fmt.Sprintf("%s%s*10^%d", sflag, txt[0:1], pointPos-1)
}

// SetValue - set decimal value from raw string
// raw string MUST be a valid number string
func (zd *Decimal) setValue(raw string) *error.Error {
	var intValS = []rune{}
	var expValS = []rune{}
	var dotNum = 0

	var rawR = []rune(raw)
	// similar with the regex parser in lexer.go
	const (
		sBegin  = 1
		sIntNum = 3
		sDotNum = 6
		sExpNum = 7
	)

	// parse string
	var state = sBegin
	var idx = 0
	var rawRL = len(rawR)
	for idx < rawRL {
		x := rawR[idx]
		// skip _
		if x == '_' {
			idx++
			continue
		}
		switch state {
		case sBegin:
			switch x {
			case '+':
				state = sIntNum
			case '-':
				state = sIntNum
				intValS = append(intValS, '-')
			case '.':
				state = sDotNum
			default:
				// real num and push first item
				state = sIntNum
				intValS = append(intValS, x)
			}
		case sIntNum:
			switch x {
			case '.':
				state = sDotNum
			case '*': // *10^ or *^
				state = sExpNum
				if rawR[idx+1] == '^' {
					idx = idx + 1
				} else {
					idx = idx + 3
				}
			case 'E', 'e':
				state = sExpNum
			default:
				intValS = append(intValS, x)
			}
		case sDotNum:
			switch x {
			case '*': // *10^ or *^
				state = sExpNum
				if rawR[idx+1] == '^' {
					idx = idx + 1
				} else {
					idx = idx + 3
				}
			case 'E', 'e':
				state = sExpNum
			default:
				intValS = append(intValS, x)
				dotNum++
			}
		case sExpNum:
			expValS = append(expValS, x)
		}
		idx++
	}

	// construct values
	if _, ok := zd.co.SetString(string(intValS), 10); !ok {
		return error.ParseFromStringError(raw)
	}

	var expInt = 0
	if len(expValS) > 0 {
		data, err := strconv.Atoi(string(expValS))
		if err != nil {
			return error.ParseFromStringError(raw)
		}
		expInt = data
	}
	zd.exp = expInt - dotNum
	return nil
}

// asInteger - if a decimal number is an integer (i.e. zd.exp >= 0), then export its
// value in (int) type; else return error.
func (zd *Decimal) asInteger() (int, *error.Error) {

	if zd.exp < 0 {
		raw := zd.String()
		return 0, error.ToIntegerError(raw)
	}
	if !zd.co.IsInt64() {
		raw := zd.String()
		return 0, error.ToIntegerError(raw)
	}
	return int(zd.co.Int64()), nil
}

// GetProperty - a null Value does not have ANY propreties.
func (zd *Decimal) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "文本*":
		return NewString(zd.String()), nil
	case "+1":
		v := ctx.arith.Add(*zd, *NewDecimalFromInt(1, 0))
		return &v, nil
	case "-1":
		v := ctx.arith.Add(*zd, *NewDecimalFromInt(-1, 0))
		return &v, nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty - a null Value does not have ANY propreties.
func (zd *Decimal) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod - a null value does not have ANY methods.
func (zd *Decimal) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	switch name {
	case "+1":
		v := ctx.arith.Add(*zd, *NewDecimalFromInt(1, 0))
		*zd = v
		return &v, nil
	case "-1":
		v := ctx.arith.Add(*zd, *NewDecimalFromInt(-1, 0))
		*zd = v
		return &v, nil
	}
	return nil, error.MethodNotFound(name)
}
