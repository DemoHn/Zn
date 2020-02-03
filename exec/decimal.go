package exec

import (
	"math/big"
	"strconv"

	"github.com/DemoHn/Zn/error"
)

// ZnDecimal - decimal number
type ZnDecimal struct {
	ZnNullable
	// decimal internal properties
	sign bool // if true, this number is NEGATIVE
	co   big.Int
	exp  int
}

func (zd *ZnDecimal) String() string {
	// TODO -
	return "[NUMBER]"
}

// SetValue - set decimal value from raw string
// raw string MUST be a valid number string
func (zd *ZnDecimal) SetValue(raw string) *error.Error {
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
				zd.sign = true
				state = sIntNum
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
			case '*': // *10^
				state = sExpNum
				idx = idx + 3
			case 'E', 'e':
				state = sExpNum
			default:
				intValS = append(intValS, x)
			}
		case sDotNum:
			switch x {
			case '*':
				state = sExpNum
				idx = idx + 3
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
