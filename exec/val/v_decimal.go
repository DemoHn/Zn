package val

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

const (
	maxDigitCount      = 18 // XXXXXXXX.XXXXXXXXXX
	maxLeadDecimalZero = 6  // 0.XXXXXX1234
	maxSciDigitCount   = 12 // 2.XXXXXXXX *10^ N
	arithPrecision     = 16
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

// NewDecimalFromFloat64 -
func NewDecimalFromFloat64(value float64) *Decimal {
	valStr := fmt.Sprintf("%v", value)
	d, _ := NewDecimal(valStr)
	return d
}

// GetExp -
func (zd *Decimal) GetExp() int {
	return zd.exp
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

// AsInteger - if a decimal number is an integer (i.e. zd.exp >= 0), then export its
// value in (int) type; else return error.
func (zd *Decimal) AsInteger() (int, *error.Error) {

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

//// arithmetic calculations
func (zd *Decimal) Add(others ...*Decimal) *Decimal {
	var result = copyDecimal(zd)
	if len(others) == 0 {
		return result
	}

	for _, item := range others {
		r1, r2 := rescalePair(result, item)
		result.co.Add(r1.co, r2.co)
		result.exp = r1.exp
	}
	return result
}

// Sub - A - B - C - D - ... = ?
func (zd *Decimal) Sub(others ...*Decimal) *Decimal {
	var result = copyDecimal(zd)
	if len(others) == 0 {
		return result
	}

	for _, item := range others {
		r1, r2 := rescalePair(result, item)
		result.co.Sub(r1.co, r2.co)
		result.exp = r1.exp
	}
	return result
}

// Mul - A * B * C * D * ... = ?, ZnDeicmal value will be copied
func (zd *Decimal) Mul(others ...*Decimal) *Decimal {
	// init result from decimal1
	var result = copyDecimal(zd)
	if len(others) == 0 {
		return result
	}

	for _, item := range others {
		result.co.Mul(result.co, item.co)
		result.exp = result.exp + item.exp
	}

	// normalize 0
	if result.co.Sign() == 0 {
		result.exp = 0
	}
	return normalizeTailZero(result)
}

// Div - A / B / C / D / ... = ?, ZnDecimal value will be copied
// notice , when one of the dividends are 0, an ArithDivZeroError will be yield
func (zd *Decimal) Div(others ...*Decimal) (*Decimal, *error.Error) {
	var result = copyDecimal(zd)
	var num10 = big.NewInt(10)
	if len(others) == 0 {
		return result, nil
	}

	// if divisor is zero, return 0 directly
	if result.co.Sign() == 0 {
		result.exp = 0
		return result, nil
	}
	// ensure divisor and all divients are postive
	neg := false
	if result.co.Sign() < 0 {
		result.co.Neg(zd.co) // co := -co
		neg = !neg
	}

	for _, item := range others {
		// check if divident is zero
		if item.co.Sign() == 0 {
			return &Decimal{}, error.ArithDivZeroError()
		}
		if item.co.Sign() < 0 {
			item.co.Neg(item.co)
			neg = !neg
		}
		adjust := 0
		// adjust bits
		// C1 < C2
		if result.co.Cmp(item.co) < 0 {
			var c2_10x = new(big.Int).Mul(item.co, num10)
			for {
				if result.co.Cmp(item.co) >= 0 && result.co.Cmp(c2_10x) < 0 {
					break
				}
				// else, C1 = C1 * 10
				result.co.Mul(result.co, num10)
				adjust = adjust + 1
			}
		} else {
			var c1_10x = new(big.Int).Mul(result.co, num10)
			for {
				if item.co.Cmp(result.co) >= 0 && item.co.Cmp(c1_10x) < 0 {
					break
				}

				item.co.Mul(item.co, num10)
				adjust = adjust - 1
			}
		}

		// exp10x = 10^(precision)
		var precFactor = arithPrecision - 1
		if adjust < 0 {
			precFactor = arithPrecision
		}
		var exp10x *big.Int
		if arithPrecision >= 18 {
			exp10x = new(big.Int).Exp(num10, num10, big.NewInt(int64(precFactor-1)))
		} else {
			exp10x = big.NewInt(int64(math.Pow10(precFactor)))
		}

		// do div
		var mul10x = exp10x.Mul(result.co, exp10x)
		var xr = new(big.Int)
		var xq, _ = result.co.DivMod(mul10x, item.co, xr) // don't use QuoRem here!

		// rounding
		if xr.Mul(xr, big.NewInt(2)).Cmp(item.co) > 0 {
			xq = xq.Add(xq, big.NewInt(1))
		}

		// get final result
		result.co = xq
		result.exp = result.exp - item.exp - adjust - precFactor
	}
	// if final result is negative, invert result.co's sign
	if neg {
		result.co.Neg(result.co)
	}
	result = normalizeTailZero(result)
	return result, nil
}

// GetProperty - a null ctx.Value does not have ANY propreties.
func (zd *Decimal) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "文本":
		return NewString(zd.String()), nil
	case "+1":
		v := zd.Add(NewDecimalFromInt(1, 0))
		return v, nil
	case "-1":
		v := zd.Add(NewDecimalFromInt(-1, 0))
		return v, nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty - a null ctx.Value does not have ANY propreties.
func (zd *Decimal) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod - a null value does not have ANY methods.
func (zd *Decimal) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	case "加":
		if err := ValidateExactParams(values, "decimal"); err != nil {
			return nil, err
		}
		item := values[0].(*Decimal)
		v := zd.Add(item)
		return v, nil
	case "减":
		if err := ValidateExactParams(values, "decimal"); err != nil {
			return nil, err
		}
		item := values[0].(*Decimal)
		v := zd.Sub(item)
		return v, nil
	case "乘":
		if err := ValidateExactParams(values, "decimal"); err != nil {
			return nil, err
		}
		item := values[0].(*Decimal)
		v := zd.Mul(item)
		return v, nil
	case "除":
		if err := ValidateExactParams(values, "decimal"); err != nil {
			return nil, err
		}
		item := values[0].(*Decimal)
		return zd.Div(item)
	case "+1":
		v := zd.Add(NewDecimalFromInt(1, 0))
		*zd = *v
		return v, nil
	case "-1":
		v := zd.Add(NewDecimalFromInt(-1, 0))
		*zd = *v
		return v, nil
	}
	return nil, error.MethodNotFound(name)
}

//// arith helper
// normalizeTailZero - remove tail zero of an integer.
// e.g: 12.500 (12500 * 10^-3)  -->  12.5 (125 * 10^-1)
func normalizeTailZero(d1 *Decimal) *Decimal {
	intTen := big.NewInt(10)
	modResult := big.NewInt(0)
	divResult := copyDecimal(d1)
	for modResult.Sign() == 0 && divResult.co.Sign() != 0 {
		d1.co.Set(divResult.co)
		d1.exp = divResult.exp

		divResult.co.DivMod(divResult.co, intTen, modResult)
		divResult.exp += 1
	}

	return d1
}

// rescalePair - make exps to be same
func rescalePair(d1 *Decimal, d2 *Decimal) (*Decimal, *Decimal) {
	intTen := big.NewInt(10)

	if d1.exp == d2.exp {
		return d1, d2
	}
	if d1.exp > d2.exp {
		// return new d1
		diff := d1.exp - d2.exp

		expVal := new(big.Int).Exp(intTen, big.NewInt(int64(diff)), nil)
		nD1 := &Decimal{
			co:  new(big.Int).Mul(d1.co, expVal),
			exp: d2.exp,
		}
		return nD1, d2
	}
	// d1.exp < d2.exp
	// return new d2
	diff := d2.exp - d1.exp

	expVal := new(big.Int).Exp(intTen, big.NewInt(int64(diff)), nil)
	nD2 := &Decimal{
		co:  new(big.Int).Mul(d2.co, expVal),
		exp: d1.exp,
	}
	return d1, nD2
}

func copyDecimal(in *Decimal) *Decimal {
	x := new(big.Int)
	return &Decimal{
		co:  x.Set(in.co),
		exp: in.exp,
	}
}

// setValue - set decimal value from raw string
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
