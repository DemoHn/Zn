package val

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/exec/ctx"
)

func TestReduce_Array_RHS_OK(t *testing.T) {
	cases := []struct {
		name   string
		array  []ctx.Value
		index  int
		expect ctx.Value
	}{
		{
			"normal RHS value",
			[]ctx.Value{
				NewString("64"),
				NewString("128"),
			},
			0,
			NewString("64"),
		},
		{
			"normal RHS value #2",
			[]ctx.Value{
				NewBool(true),
				NewString("Hello World"),
				NewArray([]ctx.Value{
					NewString("A"),
					NewString("B"),
				}),
			},
			2,
			NewArray([]ctx.Value{
				NewString("A"),
				NewString("B"),
			}),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewContext(nil)

			iv := &IV{
				reduceType: IVTypeArray,
				root:       NewArray(tt.array),
				index:      tt.index,
			}

			v, err := iv.ReduceRHS(c)
			if err != nil {
				t.Errorf("reduce() should have no error - but error: %s occured", err.Error())
				return
			}

			if !reflect.DeepEqual(v, tt.expect) {
				t.Errorf("not same: expect=%v, reduced=%v", StringifyValue(v), StringifyValue(tt.expect))
			}
		})
	}
}

func TestReduce_Array_RHS_FAIL(t *testing.T) {
	cases := []struct {
		name    string
		array   []ctx.Value
		index   int
		errCode int
	}{
		{
			"negative index",
			[]ctx.Value{
				NewString("64"),
			},
			-1,
			0x2401,
		},
		{
			"exceed range",
			[]ctx.Value{
				NewArray([]ctx.Value{}),
			},
			2,
			0x2401,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewContext(nil)

			iv := &IV{
				reduceType: IVTypeArray,
				root:       NewArray(tt.array),
				index:      tt.index,
			}

			_, err := iv.ReduceRHS(c)
			if err == nil {
				t.Errorf("reduce() expect error - but no error")
				return
			}

			if int(err.GetCode()) != tt.errCode {
				t.Errorf("expect error code: 0x%x, got: 0x%x", tt.errCode, int(err.GetCode()))
			}
		})
	}
}

func TestReduace_Array_LHS_OK(t *testing.T) {
	cases := []struct {
		name  string
		array []ctx.Value
		index int
		input ctx.Value
	}{
		{
			"change first value",
			[]ctx.Value{
				NewString("First Item"),
			},
			0,
			NewString("Noah"),
		},
		{
			"change value #2",
			[]ctx.Value{
				NewString("First Item"),
				NewString("Second Item"),
				NewString("Third Item"),
				NewString("Last Item"),
			},
			2,
			NewString("Noah"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewContext(nil)

			root := NewArray(tt.array)
			// construct IV
			iv := &IV{
				reduceType: IVTypeArray,
				root:       root,
				index:      tt.index,
			}

			err := iv.ReduceLHS(c, tt.input)
			if err != nil {
				t.Errorf("reduceLHS() should have no error - but error: %s occured", err.Error())
				return
			}

			if !reflect.DeepEqual(root.value[tt.index], tt.input) {
				t.Errorf("unexpected reduced value, expect: %#v, got: %#v", tt.input, root.value[tt.index])
				return
			}
		})
	}
}

func TestReduce_Array_LHS_FAIL(t *testing.T) {
	cases := []struct {
		name    string
		array   []ctx.Value
		index   int
		errCode int
	}{
		{
			"negative index",
			[]ctx.Value{
				NewString("64"),
			},
			-1,
			0x2401,
		},
		{
			"exceed range",
			[]ctx.Value{
				NewArray([]ctx.Value{}),
			},
			2,
			0x2401,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := ctx.NewContext(nil)

			iv := &IV{
				reduceType: IVTypeArray,
				root:       NewArray(tt.array),
				index:      tt.index,
			}

			err := iv.ReduceLHS(c, NewString(""))
			if err == nil {
				t.Errorf("reduce() expect error - but no error")
				return
			}

			if int(err.GetCode()) != tt.errCode {
				t.Errorf("expect error code: 0x%x, got: 0x%x", tt.errCode, int(err.GetCode()))
			}
		})
	}
}
