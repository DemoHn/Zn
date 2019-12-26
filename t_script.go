package main

import (
	"fmt"
	"strings"
)

// T -
type T struct {
	Name   string
	Code   string
	Pinyin []string
}

func mainT() {
	//gg := "等于"
	//for _, g := range gg {
	//	fmt.Printf("character %s is: %X \n", string(g), g)
	//}

	keywords := []T{
		T{"VarDeclare", "令", []string{"LING"}},
		T{"LogicIsI", "为", []string{"WEI"}},
		T{"LogicIsII", "是", []string{"SHI"}},
		T{"ValueAssign", "设为", []string{"SHE", "WEI"}},
		T{"MethodClaim", "如何", []string{"RU", "HE"}},
		T{"ArgumentClaim", "已知", []string{"YI", "ZHIy"}},
		T{"ReturnClaim", "返回", []string{"FAN", "HUI"}},
		T{"LogicIsNotI", "不为", []string{"BU", "WEI"}},
		T{"LogicIsNotII", "不是", []string{"BU", "SHI"}},
		T{"LogicNotEq", "不等于", []string{"BU", "DENG", "YU"}},
		T{"LogicGt", "大于", []string{"DA", "YU"}},
		T{"LogicLte", "不大于", []string{"BU", "DA", "YU"}},
		T{"LogicLt", "小于", []string{"XIAO", "YU"}},
		T{"LogicGte", "不小于", []string{"BU", "XIAO", "YU"}},
		T{"FirstArgClaim", "以", []string{"YIi"}},
		T{"FirstArgJoin", "而", []string{"ER"}},
		T{"MethodYield", "得", []string{"DE"}},
		T{"ConditionDeclare", "如果", []string{"RU", "GUO"}},
		T{"ConditionThen", "则", []string{"ZE"}},
		T{"ConditionElse", "否则", []string{"FOU", "ZE"}},
		T{"WhileDeclare", "每当", []string{"MEI", "DANG"}},
		T{"ObjectAssign", "成为", []string{"CHENG", "WEI"}},
		T{"ObjectAlias", "作为", []string{"ZUO", "WEI"}},
		T{"ObjectConstructor", "是为", []string{"SHI", "WEI"}},
		T{"ClassDeclare", "定义", []string{"DING", "YIy"}},
		T{"ClassTrait", "类比", []string{"LEI", "BI"}},
		T{"ClassThis", "其", []string{"QI"}},
		T{"ClassSelf", "此", []string{"CI"}},
		T{"Comment", "注", []string{"ZHU"}},
		T{"ProperyDeclare", "何为", []string{"HE", "WEI"}},
		T{"MethodArgClaim", "在", []string{"ZAI"}},
		T{"MethodArgJoin", "中", []string{"ZHONG"}},
		T{"LogicOr", "或", []string{"HUO"}},
		T{"LogicAnd", "且", []string{"QIE"}},
		T{"AsI", "之", []string{"ZHI"}},
		T{"AsII", "的", []string{"DEo"}},
	}
	/**
	for _, t := range tt {
		codes := []string{}
		for _, ch := range []rune(t.Code) {
			codes = append(codes, fmt.Sprintf("0x%X", ch))
		}

		code := strings.Join(codes, ", ")
		fmt.Printf("%s = []rune{%s} // %s\n", t.Name, code, t.Code)

		//fmt.Printf("%sType: \"%s<%s>\",\n", t.Name, t.Name, t.Code)
	}*/
	printKeywordConsts(keywords)
	printLeadKeywords(keywords)
}

// KData -
type KData struct {
	Pinyin  string
	Index   int
	Phrases []string
}

func printKeywordConsts(keywords []T) {
	var leadMap = map[rune]*KData{}

	index := 0
	for _, keyword := range keywords {
		rr := []rune(keyword.Code)

		for i, r := range rr {
			_, ok := leadMap[r]
			if !ok {
				leadMap[r] = &KData{
					Pinyin:  keyword.Pinyin[i],
					Index:   index,
					Phrases: []string{},
				}
				index++
			}
			// add phrases
			v := leadMap[r]
			v.Phrases = append(v.Phrases, keyword.Code)

		}
	}

	var lines = make([]string, len(leadMap)*2)
	for k, m := range leadMap {
		varname := fmt.Sprintf("Glyph%s", m.Pinyin)
		ss := string([]rune{k})
		phrases := strings.Join(m.Phrases, "，")
		lines[m.Index*2] = fmt.Sprintf("// %s - %s - %s", varname, ss, phrases)
		// add vardef
		lines[m.Index*2+1] = fmt.Sprintf("%s rune = 0x%X", varname, k)
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func printLeadKeywords(keywords []T) {
	var leadMap = map[rune]*KData{}

	index := 0
	for _, keyword := range keywords {
		rr := []rune(keyword.Code)

		r := rr[0]
		_, ok := leadMap[r]
		if !ok {
			leadMap[r] = &KData{
				Pinyin:  keyword.Pinyin[0],
				Index:   index,
				Phrases: []string{},
			}
			index++
		}
		// add phrases
		v := leadMap[r]
		v.Phrases = append(v.Phrases, keyword.Code)
	}

	var lines = make([]string, len(leadMap))
	for _, m := range leadMap {
		varname := fmt.Sprintf("Glyph%s", m.Pinyin)
		// add vardef
		lines[m.Index] = fmt.Sprintf("%s,", varname)
	}
	fmt.Println(strings.Join(lines, "\n"))
}
