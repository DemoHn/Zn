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
		{
			name: "A -> B & B -> A",
			g: map[string][]string{
				"A1": {"B1", "B2", "B3"},
				"B1": {"C1", "A1"},
				"C1": {"B1"},
			},
			result: true,
		},
		{
			name: "valid dep tree#1",
			g: map[string][]string{
				"A1": {"B1", "B2", "B3"},
				"B1": {"C1", "B2"},
				"B2": {"C1"},
				"B3": {"B2"},
				"C1": {},
			},
			result: false,
		},
		{
			name: "a large circular loop (A->D->E->F->A)",
			g: map[string][]string{
				"A": {"B", "C", "D"},
				"B": {"C"},
				"C": {},
				"D": {"E", "C", "F"},
				"E": {"C", "B"},
				"F": {"A"},
			},
			result: true,
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
