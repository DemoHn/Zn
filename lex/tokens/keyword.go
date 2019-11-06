package tokens

import "fmt"

// declare keywords
var (
	VarDeclare       = []rune{0x4EE4}                 // 令
	LogicIsI         = []rune{0x4E3A}                 // 为
	LogicIsII        = []rune{0x662F}                 // 是
	ValueAssign      = []rune{0x8BBE, 0x4E3A}         // 设为
	MethodClaim      = []rune{0x5982, 0x4F55}         // 如何
	ArgumentClaim    = []rune{0x5DF2, 0x77E5}         // 已知
	ReturnClaim      = []rune{0x8FD4, 0x56DE}         // 返回
	LogicIsNotI      = []rune{0x4E0D, 0x4E3A}         // 不为
	LogicIsNotII     = []rune{0x4E0D, 0x662F}         // 不是
	LogicNotEq       = []rune{0x4E0D, 0x7B49, 0x4E8E} // 不等于
	LogicGt          = []rune{0x5927, 0x4E8E}         // 大于
	LogicLte         = []rune{0x4E0D, 0x5927, 0x4E8E} // 不大于
	LogicLt          = []rune{0x5C0F, 0x4E8E}         // 小于
	LogicGte         = []rune{0x4E0D, 0x5C0F, 0x4E8E} // 不小于
	FirstArgClaim    = []rune{0x4EE5}                 // 以
	FirstArgJoin     = []rune{0x800C}                 // 而
	MethodYield      = []rune{0x5F97}                 // 得
	ConditionDeclare = []rune{0x5982, 0x679C}         // 如果
	ConditionThen    = []rune{0x5219}                 // 则
	ConditionElse    = []rune{0x5426, 0x5219}         // 否则
	WhileDeclare     = []rune{0x6BCF, 0x5F53}         // 每当
	ObjectAssign     = []rune{0x6210, 0x4E3A}         // 成为
	ClassDeclare     = []rune{0x5B9A, 0x4E49}         // 定义
	ClassTrait       = []rune{0x7C7B, 0x6BD4}         // 类比
	ClassSelf        = []rune{0x5176}                 // 其
	Comment          = []rune{0x6CE8}                 // 注
	ProperyDeclare   = []rune{0x4F55, 0x4E3A}         // 何为
	MethodArgClaim   = []rune{0x5728}                 // 在
	MethodArgJoin    = []rune{0x4E2D}                 // 中
	LogicOr          = []rune{0x6216}                 // 或
	LogicAnd         = []rune{0x4E14}                 // 且
)

// KeywordTokenType -
type KeywordTokenType int

// declare types
const (
	VarDeclareType KeywordTokenType = iota
	LogicIsIType
	LogicIsIIType
	ValueAssignType
	MethodClaimType
	ArgumentClaimType
	ReturnClaimType
	LogicIsNotIType
	LogicIsNotIIType
	LogicNotEqType
	LogicGtType
	LogicLteType
	LogicLtType
	LogicGteType
	FirstArgClaimType
	FirstArgJoinType
	MethodYieldType
	ConditionDeclareType
	ConditionThenType
	ConditionElseType
	WhileDeclareType
	ObjectAssignType
	ClassDeclareType
	ClassTraitType
	ClassSelfType
	CommentType
	ProperyDeclareType
	MethodArgClaimType
	MethodArgJoinType
	LogicOrType
	LogicAndType
)

// KeywordToken -
type KeywordToken struct {
	Type  KeywordTokenType
	Start int
	End   int
}

func (k KeywordToken) String(detailed bool) string {
	aliasNameMap := map[KeywordTokenType]string{
		VarDeclareType:       "VarDeclare<令>",
		LogicIsIType:         "LogicIsI<为>",
		LogicIsIIType:        "LogicIsII<是>",
		ValueAssignType:      "ValueAssign<设为>",
		MethodClaimType:      "MethodClaim<如何>",
		ArgumentClaimType:    "ArgumentClaim<已知>",
		ReturnClaimType:      "ReturnClaim<返回>",
		LogicIsNotIType:      "LogicIsNotI<不为>",
		LogicIsNotIIType:     "LogicIsNotII<不是>",
		LogicNotEqType:       "LogicNotEq<不等于>",
		LogicGtType:          "LogicGt<大于>",
		LogicLteType:         "LogicLte<不大于>",
		LogicLtType:          "LogicLt<小于>",
		LogicGteType:         "LogicGte<不小于>",
		FirstArgClaimType:    "FirstArgClaim<以>",
		FirstArgJoinType:     "FirstArgJoin<而>",
		MethodYieldType:      "MethodYield<得>",
		ConditionDeclareType: "ConditionDeclare<如果>",
		ConditionThenType:    "ConditionThen<则>",
		ConditionElseType:    "ConditionElse<否则>",
		WhileDeclareType:     "WhileDeclare<每当>",
		ObjectAssignType:     "ObjectAssign<成为>",
		ClassDeclareType:     "ClassDeclare<定义>",
		ClassTraitType:       "ClassTrait<类比>",
		ClassSelfType:        "ClassSelf<其>",
		CommentType:          "Comment<注>",
		ProperyDeclareType:   "ProperyDeclare<何为>",
		MethodArgClaimType:   "MethodArgClaim<在>",
		MethodArgJoinType:    "MethodArgJoin<中>",
		LogicOrType:          "LogicOr<或>",
		LogicAndType:         "LogicAnd<且>",
	}

	raw := aliasNameMap[k.Type]
	if detailed {
		return fmt.Sprintf("%s[%d,%d]", raw, k.Start, k.End)
	}

	return raw
}

// Position - get position
func (k KeywordToken) Position() (int, int) {
	return k.Start, k.End
}
