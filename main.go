package main

import (
	"fmt"
	"strings"
)

// T -
type T struct {
	Name string
	Code string
}

func main() {
	//gg := "等于"
	//for _, g := range gg {
	//	fmt.Printf("character %s is: %X \n", string(g), g)
	//}

	tt := []T{
		T{"VarDeclare", "	"},
		T{"LogicIsI", "为"},
		T{"LogicIsII", "是"},
		T{"ValueAssign", "设为"},
		T{"MethodClaim", "如何"},
		T{"ArgumentClaim", "已知"},
		T{"ReturnClaim", "返回"},
		T{"LogicIsNotI", "不为"},
		T{"LogicIsNotII", "不是"},
		T{"LogicNotEq", "不等于"},
		T{"LogicGt", "大于"},
		T{"LogicLte", "不大于"},
		T{"LogicLt", "小于"},
		T{"LogicGte", "不小于"},
		T{"FirstArgClaim", "以"},
		T{"FirstArgJoin", "而"},
		T{"MethodYield", "得"},
		T{"ConditionDeclare", "如果"},
		T{"ConditionThen", "则"},
		T{"ConditionElse", "否则"},
		T{"WhileDeclare", "每当"},
		T{"ObjectAssign", "成为"},
		T{"ClassDeclare", "定义"},
		T{"ClassTrait", "类比"},
		T{"ClassSelf", "其"},
		T{"Comment", "注"},
		T{"ProperyDeclare", "何为"},
		T{"MethodArgClaim", "在"},
		T{"MethodArgJoin", "中"},
		T{"LogicOr", "或"},
		T{"LogicAnd", "且"},
	}

	for _, t := range tt {
		codes := []string{}
		for _, ch := range []rune(t.Code) {
			codes = append(codes, fmt.Sprintf("0x%X", ch))
		}

		code := strings.Join(codes, ", ")
		fmt.Printf("%s = []rune{%s} // %s\n", t.Name, code, t.Code)

		//fmt.Printf("%sType: \"%s<%s>\",\n", t.Name, t.Name, t.Code)
	}
}
