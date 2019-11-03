package tokens

import (
	"testing"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

func TestSource_ConstructMarkerToken(t *testing.T) {

	tests := []struct {
		name        string
		code        string
		expectError bool
		errorCode   uint16
		token       string // we display
	}{
		{
			name:        "test all single markers",
			code:        "，：；？！&@#",
			expectError: false,
			errorCode:   0,
			token: "Comma<,>[0,0] Colon<:>[1,1] SemiColon<;>[2,2]" +
				" Question<?>[3,3] Bang<!>[4,4] Ref<&>[5,5] Annotation<@>[6,6] Hash<#>[7,7]",
		},
		{
			name:        "test double ems",
			code:        "——",
			expectError: false,
			errorCode:   0,
			token:       "Em<—>[0,1]",
		},
		{
			name:        "test double ellipsis",
			code:        "……",
			expectError: false,
			errorCode:   0,
			token:       "Ellipsis<…>[0,1]",
		},
		{
			name:        "regard other non-marker strings as empty",
			code:        "x&",
			expectError: false,
			errorCode:   0,
			token:       "NIL[0,0] Ref<&>[1,1]",
		},
		{
			name:        "only single em should throw error",
			code:        "—x",
			expectError: true,
			errorCode:   0x1101,
			token:       "",
		},
		{
			name:        "only single ellipsis should throw error",
			code:        "…y",
			expectError: true,
			errorCode:   0x1102,
			token:       "",
		},
		{
			name:        "no code, no token",
			code:        "",
			expectError: false,
			errorCode:   0,
			token:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := []rune(tt.code)
			l := lex.NewLexer(code)

			var lastError *error.Error
			for !l.End() {
				idx := l.GetIndex()
				ch := l.Next()
				token, err := ConstructMarkerToken(l, ch, idx)
				if err != nil {
					lastError = err
					break
				}
				l.AppendToken(token)
			}

			if lastError != nil && tt.expectError == false {
				t.Errorf("ConstructMarkerToken yields error: %s, want -> nil", lastError.Error())
			}
			if lastError == nil && tt.expectError == true {
				t.Errorf("ConstructMarkerToken yields: nil, want -> %x", tt.errorCode)
			}
			if lastError != nil && tt.expectError == true && lastError.GetCode() != tt.errorCode {
				t.Errorf("ConstructMarkerToken yields: error(%x), want -> error(%x)", lastError.GetCode(), tt.errorCode)
			}
			got := l.DisplayTokens()
			if got != tt.token {
				t.Errorf("ConstructMarkerToken -> %s, want -> %s", got, tt.token)
			}
		})
	}
}
