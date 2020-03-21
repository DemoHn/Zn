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
func (ai *ArithInstance) Add(decimals ...*ZnDecimal) *ZnDecimal {
	sum, _ := NewZnDecimal("0")
	for _, decimal := range decimals {
		r1, r2 := rescalePair(sum, decimal)
		newco := new(big.Int).Add(r1.co, r2.co)

		sum.co = newco
		sum.exp = r1.exp
	}

	return sum
}

// Sub - A - B - C - D - ... = ?
func (ai *ArithInstance) Sub(decimals ...*ZnDecimal) *ZnDecimal {
	sum, _ := NewZnDecimal("0")
	for _, decimal := range decimals {
		r1, r2 := rescalePair(sum, decimal)
		negco := new(big.Int).Neg(r2.co)
		newco := new(big.Int).Add(r1.co, negco)

		sum.co = newco
		sum.exp = r1.exp
	}

	return sum
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
