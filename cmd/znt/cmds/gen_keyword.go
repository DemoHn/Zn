package cmds

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type kwItem struct {
	name    string
	literal []rune
}

// GenKeywordCmd - generate keyword token definition from config file
var GenKeywordCmd = &cobra.Command{
	Use:   "gen-keyword [file]",
	Short: "根据关键词配置以生成对应（keyword）代码 - lex/keyword.go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dat, err := ioutil.ReadFile(args[0])
		if err != nil {
			panic(err)
		}
		splitInputFile(dat)
	},
}

// parse and insert to charMap & keywordMap
func splitInputFile(dat []byte) (map[rune]string, map[int]kwItem) {
	phase := 1
	lines := strings.Split(string(dat), "\n")
	charMap := map[rune]string{}
	keywordMap := map[int]kwItem{}
	for _, line := range lines {
		if strings.HasPrefix(line, "========") {
			phase = 2
			continue
		}
		// regard it as comment, ignore it
		if len(line) == 0 || strings.HasPrefix(line, "#") || strings.HasPrefix(line, " ") {
			continue
		}

		items := strings.Fields(line)
		// if not, ignore this line and parse next line
		if phase == 1 {
			if len(items) == 2 {
				r := []rune(items[1])
				charMap[r[0]] = fmt.Sprintf("Glyph%s", items[0])
			}
		} else if phase == 2 {
			if len(items) == 3 {
				t, e := strconv.Atoi(items[1])
				if e != nil {
					panic(e)
				}

				keywordMap[t] = kwItem{
					name:    fmt.Sprintf("Type%s", items[0]),
					literal: []rune(items[2]),
				}
			}
		}
	}
	return charMap, keywordMap
}

func exportKeywordLeadMap() map[rune][]int {
	keywordLeadMap := map[rune][]int{}
}

func exportCharContainMap() map[rune][]string {

}
