package lex

/**
// TestPushLine - test pushline only
func TestLineScanner_PushLine(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		result string
	}{
		{
			name:   "no push line",
			input:  "",
			result: "",
		},
		{
			name:   "push one line",
			input:  "p(2)",
			result: "Unknown<0>[0,2]",
		},
		{
			name:   "push two lines",
			input:  "p(2) p(4)",
			result: "Unknown<0>[0,2] Unknown<0>[0,4]",
		},
	}

	// run cases
	for _, tt := range cases {
		ls := NewLineScanner()
		execInput(tt.input, ls)
		// get result
		if StringifyLines(ls) != tt.result {
			t.Errorf("test pushline result expect:%s, actual:%s", tt.result, StringifyLines(ls))
		}
	}
}

func TestLineScanner_PushAndSetIndent(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		result    string
		withError bool
	}{
		{
			name:      "push and set line (space)",
			input:     "s(4,S,4) p(8) p(10) s(8,S,12) p(15)",
			result:    "Space<1>[4,8] Space<0>[0,10] Space<2>[12,15]",
			withError: false,
		},
		{
			name:      "push and set line (tag)",
			input:     "s(2,T,4) p(8) p(10) s(4,T,12) p(15)",
			result:    "Tab<2>[4,8] Tab<0>[0,10] Tab<4>[12,15]",
			withError: false,
		},
		{
			name:      "error: mix tab and space",
			input:     "s(2,T,4) p(8) p(10) s(4,S,12) p(15)",
			result:    "",
			withError: true,
		},
		{
			name:      "error: invalid space num!",
			input:     "s(7,S,4) p(15)",
			result:    "",
			withError: true,
		},
	}
	// run cases
	for _, tt := range cases {
		ls := NewLineScanner()
		err := execInput(tt.input, ls)
		// get result
		if err != nil && tt.withError == false {
			t.Errorf("test pushline should NOT throw error, but error: %v thrown", err)
		} else if err == nil && tt.withError == true {
			t.Error("test pushline result expected error, but no error thrown!")
		} else if tt.withError == false && StringifyLines(ls) != tt.result {
			t.Errorf("test pushline result expect:%s, actual:%s", tt.result, StringifyLines(ls))
		}

	}
}

// valid grammer examples:
// p(4) s(0,T,2)
// p(5) p(20)
func execInput(input string, ls *LineScanner) *error.Error {
	//re := regexp.MustCompile(`(((p)\((\d+)\))|(s)\((\d+),(T|S),(\d+)\)?!\s)`)
	//ops := re.FindAllStringSubmatch(input, -1)
	ops := strings.Split(input, " ")
	// do action
	rep := regexp.MustCompile(`^p\((\d+)\)$`)
	res := regexp.MustCompile(`^s\((\d+),(T|S),(\d+)\)$`)
	for _, op := range ops {
		rrep := rep.FindStringSubmatch(op)
		if len(rrep) > 0 {
			i, _ := strconv.Atoi(rrep[1])
			ls.PushLine(i)
		}

		rres := res.FindStringSubmatch(op)
		if len(rres) > 0 {
			c, _ := strconv.Atoi(rres[1])
			tp := IdetSpace
			if rres[2] == "T" {
				tp = IdetTab
			}
			si, _ := strconv.Atoi(rres[3])

			// exec
			if err := ls.SetIndent(c, tp, si); err != nil {
				return err
			}
		}
	}

	return nil
}
*/
