package exec

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/error"
)

// ZnDecimal - decimal number 「数值」型
type ZnDecimal struct {
	// decimal internal properties
	co  *big.Int
	exp int
}

// NewZnDecimal -
func NewZnDecimal(value string) (*ZnDecimal, *error.Error) {
	var decimal = &ZnDecimal{
		exp: 0,
		co:  big.NewInt(0),
	}

	err := decimal.setValue(value)
	return decimal, err
}

// String - show decimal display string
func (zd *ZnDecimal) String() (data string) {
	var sflag = ""
	if zd.co.Sign() < 0 {
		sflag = "-"
	}
	var txt = new(big.Int).Abs(zd.co).String()

	if zd.exp == 0 {
		data = fmt.Sprintf("%s%s", sflag, txt)
	} else if zd.exp > 0 {
		var zeros = strings.Repeat("0", zd.exp)
		data = fmt.Sprintf("%s%s%s", sflag, txt, zeros)
	} else {
		// case: zd.exp < 0
		if zd.exp+len(txt) <= 0 {
			var zeros = strings.Repeat("0", -(zd.exp + len(txt)))
			data = fmt.Sprintf("%s0.%s%s", sflag, zeros, txt)
		} else {
			pt := zd.exp + len(txt)
			data = fmt.Sprintf("%s%s.%s", sflag, txt[:pt], txt[pt:])
		}
	}
	return
}

// Compare - ZnDecimal
func (zd *ZnDecimal) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	var valR *ZnDecimal
	var targetRes = 0
	switch v := val.(type) {
	case *ZnDecimal:
		valR = v
	case *ZnNull:
		return NewZnBool(false), nil
	default:
		if cmpType == compareTypeEq || cmpType == compareTypeIs {
			return NewZnBool(false), nil
		}
		return nil, error.InvalidExprType("decimal")
	}

	switch cmpType {
	case compareTypeEq, compareTypeIs:
		targetRes = 0
	case compareTypeGt:
		targetRes = 1
	case compareTypeLt:
		targetRes = -1
	}
	r1, r2 := rescalePair(zd, valR)
	if res := r1.co.Cmp(r2.co); res == targetRes {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// SetValue - set decimal value from raw string
// raw string MUST be a valid number string
func (zd *ZnDecimal) setValue(raw string) *error.Error {
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
		return error.NewErrorSLOT("parse BigInt error")
	}

	var expInt = 0
	if len(expValS) > 0 {
		data, err := strconv.Atoi(string(expValS))
		if err != nil {
			return error.NewErrorSLOT("atoi error")
		}
		expInt = data
	}
	zd.exp = expInt - dotNum
	return nil
}

// asInteger - if a decimal number is an integer (i.e. zd.exp >= 0), then export its
// value in (int) type; else return error.
func (zd *ZnDecimal) asInteger() (int, *error.Error) {
	if zd.exp < 0 {
		return 0, error.NewErrorSLOT("this decimal not belongs to integer")
	}
	if !zd.co.IsInt64() {
		return 0, error.NewErrorSLOT("cast to int64 fail")
	}
	return int(zd.co.Int64()), nil
}