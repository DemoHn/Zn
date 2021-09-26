package lex

import (
	"testing"
)

func TestIdInRange(t *testing.T) {
	cases := []struct {
		name   string
		ids    []rune
		expect bool
	}{
		{
			name:   "operators",
			ids:    []rune{'+', '-', '*'},
			expect: true,
		},
		{
			name:   "operators (not /)",
			ids:    []rune{'/'},
			expect: false,
		},
		{
			name:   "numbers",
			ids:    []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
			expect: true,
		},
		{
			name:   "dot and underscore",
			ids:    []rune{'.', '^', '_', '%', '$'},
			expect: true,
		},
		{
			name:   "fullwidth characters",
			ids:    []rune{'＋', '－', '１', 'Ｅ', 'Ｇ', '＊', '／', '．', '＾', '＿'},
			expect: true,
		},
		{
			name:   "basic latin",
			ids:    []rune{'f', 'Z', 'u', 'Ñ', 'â', 'æ', 'ö', 'µ', 'Š', 'Ʒ'},
			expect: true,
		},
		{
			name:   "units",
			ids:    []rune{'°', '²', '³'},
			expect: true,
		},
		{
			name:   "kg, cm, cm2",
			ids:    []rune{'½', '㎏', '㎠', '㎝'},
			expect: false,
		},
		{
			name:   "last chars",
			ids:    []rune{0xffdb, 0xffdc, 0xffda},
			expect: true,
		},
		{
			name:   "Chinese characters",
			ids:    []rune{'之', '乎', '者', '也', '歷', '灣'},
			expect: true,
		},
		{
			name:   "Japenese characters (kanas & kanjis)",
			ids:    []rune{'あ', 'ぬ', 'り', 'リ', 'ツ', 'じ', 'ッ', 'ぱ', '認', '読'},
			expect: true,
		},
		{
			name:   "Korean characters (Hangul)",
			ids:    []rune{'ㅏ', 'ㅠ', 'ㅎ', '바', '음'},
			expect: true,
		},
		{
			name:   "other punctuators",
			ids:    []rune{',', '@'},
			expect: false,
		},
		{
			name:   "other random in range characters",
			ids:    []rune{0x1f38, 0xa950, 0xff49, 0xa392},
			expect: true,
		},
		{
			name:   "other random not in range characters",
			ids:    []rune{0xd7, 0xc91, 0x1dfa},
			expect: false,
		},
		{
			name:   "markers",
			ids:    []rune{'【', '】', '‘', '’', '「', '」', '<', '>', '？', '！', '；', '：', '~'},
			expect: false,
		},
		{
			name:   "other characters larger than 0xffff or smaller than 0",
			ids:    []rune{0x1ff1d, 0x2ffff, -1, 10},
			expect: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			for _, id := range tt.ids {
				res := idInRange(id)
				if res != tt.expect {
					t.Errorf("identifier %c (0x%x) expects %v, got %v", id, id, tt.expect, res)
				}
			}
		})
	}
}
