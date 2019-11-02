package tokens

import (
	"testing"

	"github.com/DemoHn/Zn/lex"
)

func TestSource_ConstructMarkerToken(t *testing.T) {

	tests := []struct {
		name        string
		code        string
		expectError bool
		errorCode   int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := []rune(tt.code)
			l := lex.NewLexer(code)
			for idx, ch := range code {
				token, _ := ConstructMarkerToken(l, ch, idx)
				l.AppendToken(token)
			}

			got := l.DisplayTokens()
			if got != tt.token {
				t.Errorf("ConstructMarkerToken -> %s, want -> %s", got, tt.token)
			}
		})
	}
}
