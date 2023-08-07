package runtime

import "testing"

type graph = map[string][]string

type testCase struct {
	g      graph
	name   string
	result bool
}

func TestCircularDependency_CoreBFS(t *testing.T) {
	cases := []testCase{
		{
			name: "standard chain: A->B B->C C->D",
			g: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {"D"},
				"D": {},
			},
			result: false,
		},
		{
			name: "standard circle: A->B B->C C->A",
			g: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {"A"},
			},
			result: true,
		},
		{
			name: "C depends by A & B (triangle): A->B B->C A->C",
			g: map[string][]string{
				"A": {"B", "C"},
				"B": {"C"},
				"C": {},
			},
			result: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			moduleGraph := ModuleGraph{
				depGraph: tt.g,
			}
			res := moduleGraph.checkCircularDepedencyBFS()
			if tt.result != res {
				t.Errorf("expected %v, got %v", tt.result, res)
			}
		})
	}
}
