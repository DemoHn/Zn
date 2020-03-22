package exec

import "math/big"

// ArithInstance - arithmetic calculation (including + - * /) instance
type ArithInstance struct {
	precision int
}

// NewArithInstance -
func NewArithInstance(precision int) *ArithInstance {
	return &ArithInstance{precision}
}

// Add - A + B + C + D + ... = ?
func (ai *ArithInstance) Add(decimal1 *ZnDecimal, others ...*ZnDecimal) *ZnDecimal {
	var result = copyZnDecimal(decimal1)
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
func (ai *ArithInstance) Sub(decimal1 *ZnDecimal, others ...*ZnDecimal) *ZnDecimal {
	var result = copyZnDecimal(decimal1)
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
func (ai *ArithInstance) Mul(decimal1 *ZnDecimal, others ...*ZnDecimal) *ZnDecimal {
	// init result from decimal1
	var result = copyZnDecimal(decimal1)
	if len(others) == 0 {
		return result
	}

	for _, item := range others {
		result.co.Mul(result.co, item.co)
		result.exp = result.exp + item.exp
	}

	return result
}

//// arith helper

// rescalePair - make exps to be same
func rescalePair(d1 *ZnDecimal, d2 *ZnDecimal) (*ZnDecimal, *ZnDecimal) {
	intTen := big.NewInt(10)

	if d1.exp == d2.exp {
		return d1, d2
	}
	if d1.exp > d2.exp {
		// return new d1
		diff := d1.exp - d2.exp

		expVal := new(big.Int).Exp(intTen, big.NewInt(int64(diff)), nil)
		nD1 := &ZnDecimal{
			co:  new(big.Int).Mul(d1.co, expVal),
			exp: d2.exp,
		}
		return nD1, d2
	}
	// d1.exp < d2.exp
	// return new d2
	diff := d2.exp - d1.exp

	expVal := new(big.Int).Exp(intTen, big.NewInt(int64(diff)), nil)
	nD2 := &ZnDecimal{
		co:  new(big.Int).Mul(d2.co, expVal),
		exp: d1.exp,
	}
	return d1, nD2
}

// copyDecimal - duplicate deicmal value to a new variable
func copyZnDecimal(old *ZnDecimal) *ZnDecimal {
	var result ZnDecimal
	var newco big.Int
	result.exp = old.exp
	newco = *(old.co)
	result.co = &newco

	return &result
}
