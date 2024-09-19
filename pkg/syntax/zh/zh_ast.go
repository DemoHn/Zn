package zh

import (
	"github.com/DemoHn/Zn/pkg/syntax"
)

type consumerFunc func()

//////// Parse Methods

//// NOTE: the following methods are all using panic -> recover for zerr.management.
//// This is to expect eliminating `err != nil` statements.

// ParseStatement - a program consists of statements
//
// CFG:
// syntax.Statement -> VarDeclareStmt
//           -> BranchStmt
//           -> WhileLoopStmt
//           -> IterateStmt
//           -> FunctionDeclareStmt
//           -> FunctionReturnStmt
//           -> VarInputStmt
//           -> VOStmt
//           -> ImportStmt
//           -> ClassStmt
//           -> Expr
//           -> ；
func ParseStatement(p *ParserZH) syntax.Statement {
	var validTypes = []uint8{
		TypeStmtSep,
		TypeDeclareW,
		TypeCondW,
		TypeFuncW,
		TypeReturnW,
		TypeInputW,
		TypeWhileLoopW,
		TypeVarOneW,
		TypeIteratorW,
		TypeObjDefineW,
		TypeImportW,
		TypeThrowErrorW,
		TypeBreakW,
		TypeContinueW,
	}
	// remove lineTerminationFlag at the beginning of executing
	// a statement
	p.unsetStmtCompleteFlag()
	match, tk := p.tryConsume(validTypes...)

	var s syntax.Statement
	if match {
		switch tk.Type {
		case TypeStmtSep:
			// skip them because it's meaningless for syntax parsing
			return new(syntax.EmptyStmt)
		case TypeDeclareW:
			s = ParseVarDeclareStmt(p)
		case TypeCondW:
			s = ParseBranchStmt(p)
		case TypeFuncW:
			// constructor function: 如何成为XX？
			matchC, _ := p.tryConsume(TypeObjNewW)
			if matchC {
				s = ParseConstructorDeclareStmt(p)
			} else {
				// normal function: 如何xxx？
				s = ParseFunctionDeclareStmt(p)
			}
		case TypeReturnW:
			s = ParseFunctionReturnStmt(p)
		case TypeInputW:
			s = ParseVarInputStmt(p)
		case TypeWhileLoopW:
			s = ParseWhileLoopStmt(p)
		case TypeVarOneW:
			s = ParseVarOneLeadStmt(p) // parse any statements leads with 「以」
		case TypeIteratorW:
			s = ParseIteratorStmt(p)
		case TypeObjDefineW:
			s = ParseClassDeclareStmt(p)
		case TypeImportW:
			s = ParseImportStmt(p)
		case TypeThrowErrorW:
			s = ParseThrowExceptionStmt(p)
		case TypeBreakW:
			s = ParseBreakStmt(p)
		case TypeContinueW:
			s = ParseContinueStmt(p)
		}
		p.setStmtCurrentLine(s, tk)
	} else {
		// other case, parse syntax.syntax.Expression
		s = ParseExpression(p)
	}

	// normally, a complete statement should occupy a whole line
	// or following a stmt
	if !(p.stmtCompleteFlag || p.meetStmtBreak()) {
		panic(p.getInvalidSyntaxPeek())
	}
	return s
}

// ParseExpression - parse an syntax.Expression, see the following CFG for details
//
// CFG:
// Expr  -> AndE Expr'
// Expr' -> 或 AndE Expr'
//       ->
//
// AndE  -> EqE AndE'
// AndE' -> 且 EqE AndE'
//       ->
//
// EqE   -> VaE  EqE'
// EqE'  -> 等于  VaE
//       -> [不等于  /=]  VaE
//       -> [小于  <]  VaE
//       -> [不小于  >=]  VaE
//       -> [大于  >]  VaE
//       -> [不大于  <=] VaE
//       ->
//
// VaE   -> ArE VaE'
// VaE'  -> 为  IdxE
//       -> 为  & IdxE
//       ->
//
// precedences:
//
// #, #{}  >  *，/  >  +，-  >  为  >  等于，大于，etc.  >  且  >  或
func ParseExpression(p *ParserZH) syntax.Expression {
	cfg := syntax.EqMarkConfig{
		AsVarAssign: true,
	}
	return parseExpressionLv1(p, cfg)
}

func ParseExpressionMAP(p *ParserZH) syntax.Expression {
	cfg := syntax.EqMarkConfig{
		AsMapSign: true,
	}

	return parseExpressionLv1(p, cfg)
}

// parseExpressionLv1 - X 或 Y
func parseExpressionLv1(p *ParserZH, cfg syntax.EqMarkConfig) syntax.Expression {
	var parseTail func(syntax.Expression) syntax.Expression

	parseTail = func(el syntax.Expression) syntax.Expression {
		if match, tk := p.tryConsume(TypeLogicOrW); match {
			exprR := parseExpressionLv2(p, cfg)
			finalExpr := &syntax.LogicExpr{
				Type:      syntax.LogicOR,
				LeftExpr:  el,
				RightExpr: exprR,
			}
			p.setStmtCurrentLine(finalExpr, tk)
			return parseTail(finalExpr)
		}
		return el
	}

	exprL := parseExpressionLv2(p, cfg)
	return parseTail(exprL)
}

// parseExpressionLv2 - X 且 Y
func parseExpressionLv2(p *ParserZH, cfg syntax.EqMarkConfig) syntax.Expression {
	var parseTail func(syntax.Expression) syntax.Expression

	parseTail = func(el syntax.Expression) syntax.Expression {
		if match, tk := p.tryConsume(TypeLogicAndW); match {
			exprR := parseExpressionLv3(p, cfg)
			finalExpr := &syntax.LogicExpr{
				Type:      syntax.LogicAND,
				LeftExpr:  el,
				RightExpr: exprR,
			}
			p.setStmtCurrentLine(finalExpr, tk)
			return parseTail(finalExpr)
		}
		return el
	}

	exprL := parseExpressionLv3(p, cfg)
	return parseTail(exprL)
}

// parseExpressionLv3 - X 等于/不等于 Y
// NOTE: by default '=' means nothing!
func parseExpressionLv3(p *ParserZH, cfg syntax.EqMarkConfig) syntax.Expression {
	validTypes := []uint8{
		TypeLogicEqualW, TypeEqualMark,
		TypeLogicNotEqW, TypeNEMark,
		TypeLogicGtW, TypeGTMark,
		TypeLogicGteW, TypeGTEMark,
		TypeLogicLtW, TypeLTMark,
		TypeLogicLteW, TypeLTEMark,
		TypeLogicYesW, TypeLogicNoW,
	}

	logicTypeMap := map[uint8]uint8{
		TypeLogicEqualW: syntax.LogicEQ,
		TypeLogicNotEqW: syntax.LogicNEQ,
		TypeNEMark:      syntax.LogicNEQ,
		TypeLogicGtW:    syntax.LogicGT,
		TypeGTMark:      syntax.LogicGT,
		TypeLogicGteW:   syntax.LogicGTE,
		TypeGTEMark:     syntax.LogicGTE,
		TypeLogicLtW:    syntax.LogicLT,
		TypeLTMark:      syntax.LogicLT,
		TypeLogicLteW:   syntax.LogicLTE,
		TypeLTEMark:     syntax.LogicLTE,
		TypeEqualMark:   syntax.LogicEQ,
		TypeLogicYesW:   syntax.LogicXEQ,
		TypeLogicNoW:    syntax.LogicXNEQ,
	}

	exprL := parseExpressionLv4(p, cfg)
	if match, tk := p.tryConsume(validTypes...); match {
		exprR := parseExpressionLv4(p, cfg)
		finalExpr := &syntax.LogicExpr{
			Type:      logicTypeMap[tk.Type],
			LeftExpr:  exprL,
			RightExpr: exprR,
		}

		p.setStmtCurrentLine(finalExpr, tk)
		return finalExpr
	}
	return exprL
}

// parseExpressionLv4 - X 设为 Y
// NOTE: by default '=' means nothing!
func parseExpressionLv4(p *ParserZH, cfg syntax.EqMarkConfig) syntax.Expression {
	validTypes := []uint8{
		TypeAssignW,
	}

	if cfg.AsVarAssign {
		validTypes = append(validTypes, TypeAssignMark)
	}

	exprL := ParseArithExpr(p)
	if match, tk := p.tryConsume(validTypes...); match {
		// parse &
		refMarkForLogicYes := false
		if match2, _ := p.tryConsume(TypeObjRef); match2 {
			refMarkForLogicYes = true
		}

		vid, ok := exprL.(syntax.Assignable)
		if !ok {
			panic(p.getExprMustTypeIDPeek())
		}
		exprR := ParseArithExpr(p)
		finalExpr := &syntax.VarAssignExpr{
			TargetVar:  vid,
			RefMark:    refMarkForLogicYes,
			AssignExpr: exprR,
		}

		p.setStmtCurrentLine(finalExpr, tk)
		return finalExpr
	}
	return exprL
}

// ParseArithExpr - parse arithmetic expression for calculation
//
// CFG:
// ArE   -> MdE ArE'
//
// ArE'  -> + MdE ArE'
//       -> - MdE ArE'
//       ->
//
// MdE   -> MemE MdE'
//
// MdE'  -> * MemE MdE'
//       -> / MemE MdE'
//       ->
func ParseArithExpr(p *ParserZH) syntax.Expression {
	var parseTail func(syntax.Expression) syntax.Expression

	parseTail = func(el syntax.Expression) syntax.Expression {
		if match, tk := p.tryConsume(TypePlus, TypeMinus); match {
			exprR := parseArithMulDivExpr(p)

			t := syntax.ArithAdd
			if tk.Type == TypeMinus {
				t = syntax.ArithSub
			}
			finalExpr := &syntax.ArithExpr{
				Type:      t,
				LeftExpr:  el,
				RightExpr: exprR,
			}
			p.setStmtCurrentLine(finalExpr, tk)
			return parseTail(finalExpr)
		}
		return el
	}

	exprL := parseArithMulDivExpr(p)
	return parseTail(exprL)
}

func parseArithMulDivExpr(p *ParserZH) syntax.Expression {
	var parseTail func(syntax.Expression) syntax.Expression

	parseTail = func(el syntax.Expression) syntax.Expression {
		if match, tk := p.tryConsume(TypeMultiply, TypeDivision, TypeIntDivMark, TypeModuloMark); match {
			exprR := ParseMemberExpr(p)

			var t uint8
			switch tk.Type {
			case TypeMultiply:
				t = syntax.ArithMul
			case TypeDivision:
				t = syntax.ArithDiv
			case TypeIntDivMark:
				t = syntax.ArithIntDiv
			case TypeModuloMark:
				t = syntax.ArithModulo
			}

			finalExpr := &syntax.ArithExpr{
				Type:      t,
				LeftExpr:  el,
				RightExpr: exprR,
			}
			p.setStmtCurrentLine(finalExpr, tk)
			return parseTail(finalExpr)
		}
		return el
	}

	exprL := ParseMemberExpr(p)
	return parseTail(exprL)
}

// ParseMemberExpr -
//
// CFG:
//
// MemE  -> 其 PropE' IdxE'
//       -> BsE IdxE'
//
// IdxE' -> #  ID   IdxE'
//       -> #  String   IdxE'
//       -> #  {  Expr  }  IdxE'
//       -> 之  CallE' IdxE'
//       ->
//
// CallE' -> FuncID
//
// PropE' -> ID
//
// FuncID -> ID
func ParseMemberExpr(p *ParserZH) syntax.Expression {
	// internal functions
	var calleeTailParser func(bool, uint8, syntax.Expression) *syntax.MemberExpr
	var memberTailParser func(syntax.Expression) syntax.Expression

	// specially parsing items after 之 or 其
	calleeTailParser = func(hasRoot bool, rootType uint8, expr syntax.Expression) *syntax.MemberExpr {
		memberExpr := &syntax.MemberExpr{
			Root:     nil,
			RootType: rootType,
		}
		if hasRoot {
			memberExpr.Root = expr
		}

		if match, tk := p.tryConsume(TypeIdentifier); match {
			id := newID(p, tk)
			p.setStmtCurrentLine(id, tk)
			memberExpr.MemberType = syntax.MemberID
			memberExpr.MemberID = id

			return memberExpr
		}
		panic(p.getInvalidSyntaxPeek())
	}

	memberTailParser = func(expr syntax.Expression) syntax.Expression {
		mExpr := &syntax.MemberExpr{}
		// default rootType is RootTypeExpr
		mExpr.RootType = syntax.RootTypeExpr

		match, tk := p.tryConsume(TypeMapHash, TypeObjDotW, TypeObjDotIIW)
		if !match {
			return expr
		}
		p.setStmtCurrentLine(mExpr, tk)
		mExpr.Root = expr
		switch tk.Type {
		case TypeMapHash:
			match2, tk2 := p.tryConsume(TypeIdentifier, TypeString, TypeStmtQuoteL)
			if match2 {
				// set memberType
				mExpr.MemberType = syntax.MemberIndex
				switch tk2.Type {
				case TypeIdentifier:
					mExpr.MemberIndex = newID(p, tk2)
				case TypeString:
					mExpr.MemberIndex = newString(p, tk2)
				case TypeStmtQuoteL:
					mExpr.MemberIndex = ParseExpression(p)

					// #2. parse tail brace
					p.consume(TypeStmtQuoteR)
				}
				return memberTailParser(mExpr)
			}
			panic(p.getInvalidSyntaxPeek())
		case TypeObjDotW, TypeObjDotIIW:
			newExpr := calleeTailParser(true, syntax.RootTypeExpr, expr)
			// replace current memberExpr as newExpr
			return memberTailParser(newExpr)
		}

		panic(p.getInvalidSyntaxPeek())
	}

	// #1. parse 其 expr
	match, _ := p.tryConsume(TypeObjThisW) // 其
	if match {
		rootType := syntax.RootTypeProp // 其
		newExpr := calleeTailParser(false, rootType, nil)
		return memberTailParser(newExpr)
	}
	// #1. parse basic expr
	rootExpr := ParseBasicExpr(p)
	return memberTailParser(rootExpr)
}

// ParseBasicExpr - parse general basic syntax.Expression
//
// CFG:
// BsE   -> { E }
//       -> （ FuncID ： E、E、...）
//       -> 以 E （ FuncID ： E、E、...）
//       -> ID
//       -> String
//       -> ArrayList
//
// FuncID -> ID
func ParseBasicExpr(p *ParserZH) syntax.Expression {
	var validTypes = []uint8{
		TypeIdentifier,
		TypeString,
		TypeArrayQuoteL,
		TypeStmtQuoteL,
		TypeFuncQuoteL,
		TypeVarOneW,
	}

	match, tk := p.tryConsume(validTypes...)
	if match {
		var e syntax.Expression
		switch tk.Type {
		case TypeIdentifier:
			e = newID(p, tk)
		case TypeString:
			e = newString(p, tk)
		case TypeArrayQuoteL:
			e = ParseArrayExpr(p)
		case TypeStmtQuoteL:
			e = ParseExpression(p)
			p.consume(TypeStmtQuoteR)
		case TypeFuncQuoteL:
			matchNewObj, _ := p.tryConsume(TypeObjNewW)
			if matchNewObj {
				e = ParseObjNewExpr(p)
			} else {
				e = ParseFuncCallExpr(p, true)
			}
		case TypeVarOneW:
			e = ParseMemberFuncCallExpr(p)
		}
		p.setStmtCurrentLine(e, tk)
		return e
	}
	panic(p.getInvalidSyntaxPeek())
}

// ParseArrayExpr - yield ArrayExpr node (support both hashMap and arrayList)
// CFG:
// ArrayExpr -> 【 ItemList 】
//           -> 【 HashMapList 】
//           -> 【 = 】
// ItemList  -> Expr ItemList
//           ->
//
// HashMapList -> KeyID = Expr HashMapTail
//
// HashMapTail -> KeyID = Expr HashMapTail
//             ->
//
// KeyID     -> ID
//           -> String
//           -> Number
func ParseArrayExpr(p *ParserZH) syntax.UnionMapList {
	// #0. try to match if empty
	if match, emptyExpr := tryParseEmptyMapList(p); match {
		return emptyExpr
	}

	// define ArrayExpr & HashMapExpr
	var ar = &syntax.ArrayExpr{
		Items: []syntax.Expression{},
	}
	var hm = &syntax.HashMapExpr{
		KVPair: []syntax.HashMapKeyValuePair{},
	}

	var isArrayType = true
	// #1. consume first syntax.Expression
	exprI := ParseExpressionMAP(p)
	if match, tk := p.tryConsume(TypeAssignMark, TypePauseCommaSep, TypeArrayQuoteR); match {
		switch tk.Type {
		case TypeArrayQuoteR:
			ar.Items = append(ar.Items, exprI)
			return ar
		case TypeAssignMark:
			isArrayType = false
			// parse right expr
			exprR := ParseExpressionMAP(p)

			hm.KVPair = append(hm.KVPair, syntax.HashMapKeyValuePair{
				Key:   exprI,
				Value: exprR,
			})
			p.unsetStmtCompleteFlag()
		case TypePauseCommaSep:
			isArrayType = true
			// append item on array
			ar.Items = append(ar.Items, exprI)
		}
	} else {
		panic(p.getInvalidSyntaxPeek())
	}

	if isArrayType {
		// parse array like 【1、2、3、4、5】
		for {
			// if not, parse next expr
			expr := ParseExpressionMAP(p)
			ar.Items = append(ar.Items, expr)

			// if parse to end
			if match, tk := p.tryConsume(TypeArrayQuoteR, TypePauseCommaSep); match {
				if tk.Type == TypeArrayQuoteR {
					return ar
				}
			} else {
				panic(p.getInvalidSyntaxPeek())
			}
		}
	} else {
		// parse hashmap like 【A = 1，B = 2】
		for {
			if match, _ := p.tryConsume(TypeArrayQuoteR); match {
				return hm
			}

			exprL := ParseExpressionMAP(p)
			p.consume(TypeAssignMark)
			exprR := ParseExpressionMAP(p)

			hm.KVPair = append(hm.KVPair, syntax.HashMapKeyValuePair{
				Key:   exprL,
				Value: exprR,
			})
			p.unsetStmtCompleteFlag()
		}
	}
}

func tryParseEmptyMapList(p *ParserZH) (bool, syntax.UnionMapList) {
	emptyTrialTypes := []uint8{
		TypeArrayQuoteR, // for empty array
		TypeAssignMark,  // for empty hashmap
	}

	if match, tk := p.tryConsume(emptyTrialTypes...); match {
		switch tk.Type {
		case TypeArrayQuoteR:
			e := &syntax.ArrayExpr{Items: []syntax.Expression{}}
			p.setStmtCurrentLine(e, tk)
			return true, e
		case TypeAssignMark:
			p.consume(TypeArrayQuoteR)
			e := &syntax.HashMapExpr{KVPair: []syntax.HashMapKeyValuePair{}}
			p.setStmtCurrentLine(e, tk)
			return true, e
		}
	}
	return false, nil
}

// ParseFuncCallExpr - yield FuncCallExpr node
//
// CFG:
// FuncCallExpr  -> （ FuncID ： pcommaList ）YieldResultTail
//               -> （ FuncID ） YieldResultTail
//               -> （ 新建 ClassID ： pcommaList ）
//               -> （ 新建 ClassID ：）
// pcommaList     -> E pcommaListTail
// pcommaListTail -> 、 E pcommaListTail
//               ->
//
// FuncID   -> ID
//          -> Number
//
// YieldResultTail  ->  得到 ID
//                  ->
func ParseFuncCallExpr(p *ParserZH, parseYieldResult bool) *syntax.FuncCallExpr {
	var callExpr = &syntax.FuncCallExpr{
		Params:      []syntax.Expression{},
		YieldResult: nil,
	}
	// #1. parse ID
	callExpr.FuncName = parseFuncID(p)
	// #2. parse colon (maybe there's no params)
	match, _ := p.tryConsume(TypeFuncCall)
	if match {
		// #2.1 parse comma list
		parsePauseCommaList(p, func() {
			expr := ParseExpression(p)
			callExpr.Params = append(callExpr.Params, expr)
		})
	}

	// #3. parse right quote
	p.consume(TypeFuncQuoteR)

	// #4. parse yield result call
	if parseYieldResult {
		match2, _ := p.tryConsume(TypeGetResultW)
		if match2 {
			id := parseID(p)
			callExpr.YieldResult = id
		}
	}
	return callExpr
}

// ParseMemberFuncCallExpr - 以 ... （‹方法名›）
// CFG:
//
// FuncExpr -> 以 Expr （ FuncID ： commaList ）
//
// FuncID  -> ID
//         -> Number
func ParseMemberFuncCallExpr(p *ParserZH) *syntax.MemberMethodExpr {
	result := &syntax.MemberMethodExpr{}
	result.Root = ParseExpression(p)

	p.consume(TypeFuncQuoteL)
	// parse first function in function chain
	funcExprI := ParseFuncCallExpr(p, false)
	result.MethodChain = append(result.MethodChain, funcExprI)

	// then parse 、 and （  (for chain function)
	for {
		match, _ := p.tryConsume(TypePauseCommaSep)
		if !match {
			break
		}
		// parse （
		p.consume(TypeFuncQuoteL)
		funcExprN := ParseFuncCallExpr(p, false)
		result.MethodChain = append(result.MethodChain, funcExprN)
	}

	if match, _ := p.tryConsume(TypeGetResultW); match {
		id := parseID(p)
		result.YieldResult = id
	}

	return result
}

// ParseVarDeclareStmt - yield VarDeclare node
// CFG:
// VarDeclare -> 令 VDPair
//
// VDPair     -> VDItem VDPairTail
//
// VDPairTail -> VDItem VDPairTail
//            ->
//
// VDItem     -> IdfList 设为 Expr
//            -> IdfList 恒为 Expr
//
//    IdfList -> I I'
//         I' -> 、I I'
//            ->
//
// or block declaration:
//
// VarDeclare -> 令 ：
//           ...
//           ...     I3 、 I4、 I5 ...
func ParseVarDeclareStmt(p *ParserZH) *syntax.VarDeclareStmt {
	vNode := &syntax.VarDeclareStmt{
		AssignPair: []syntax.VDAssignPair{},
	}

	// #01. try to read colon
	// if colon exists -> parse comma list by block
	// if colon not exists -> parse comma list inline
	if match, _ := p.tryConsume(TypeFuncCall); match {
		expected, blockIndent := p.expectBlockIndent()
		if !expected {
			panic(p.getInvalidSyntaxCurr())
		}

		parseItemListBlock(p, blockIndent, func() {
			p.unsetStmtCompleteFlag()

			match, _ := p.tryConsume(TypeStmtSep)
			if match {
				return
			}
			// there are at least ONE vdAssignPair on each line!
			assignPair := parseVDAssignPair(p)

			// normally, a complete statement should occupy a whole line
			// or following a stmt
			if !(p.stmtCompleteFlag || p.meetStmtBreak()) {
				panic(p.getInvalidSyntaxPeek())
			}
			vNode.AssignPair = append(vNode.AssignPair, assignPair)
		})
	} else {
		// #02. consume identifier declare list (comma list) inline
		// there is ONLY ONE vdAssignPair along the statement!
		vNode.AssignPair = []syntax.VDAssignPair{
			parseVDAssignPair(p),
		}
	}

	return vNode
}

func parseVDAssignPair(p *ParserZH) syntax.VDAssignPair {
	var idfList []*syntax.ID

	// #1. parse identifier
	parsePauseCommaList(p, func() {
		id := parseID(p)
		idfList = append(idfList, id)
	})

	// parse keyword
	validKeywords := []uint8{
		TypeAssignW,
		TypeAssignMark,
		TypeAssignConstW,
	}
	match, tk := p.tryConsume(validKeywords...)
	if !match {
		panic(p.getInvalidSyntaxPeek())
	}

	switch tk.Type {
	case TypeAssignW, TypeAssignMark:
		refMark := false
		if match, _ := p.tryConsume(TypeObjRef); match {
			refMark = true
		}
		expr := ParseExpression(p)

		return syntax.VDAssignPair{
			Type:       syntax.VDTypeAssign,
			Variables:  idfList,
			RefMark:    refMark,
			AssignExpr: expr,
		}
	case TypeAssignConstW:
		refMark := false
		if match, _ := p.tryConsume(TypeObjRef); match {
			refMark = true
		}
		expr := ParseExpression(p)

		return syntax.VDAssignPair{
			Type:       syntax.VDTypeAssignConst,
			Variables:  idfList,
			RefMark:    refMark,
			AssignExpr: expr,
		}
	default: // ObjNewW
		panic(p.getInvalidSyntaxPeek())
	}
}

// ParseObjNewExpr - yield ObjNewExpr node
// CFG:
// ObjNewExpr -> （创建 ID ：Expr、 Expr、 ... ）
//            -> （创建 ID ）
func ParseObjNewExpr(p *ParserZH) *syntax.ObjNewExpr {
	objClass := parseID(p)
	if match, _ := p.tryConsume(TypeFuncCall); match {
		params := []syntax.Expression{}
		parsePauseCommaList(p, func() {
			e := ParseExpression(p)
			params = append(params, e)
		})
		p.consume(TypeFuncQuoteR)
		return &syntax.ObjNewExpr{
			ClassName: objClass,
			Params:    params,
		}
	} else {
		p.consume(TypeFuncQuoteR)

		return &syntax.ObjNewExpr{
			ClassName: objClass,
			Params:    []syntax.Expression{},
		}
	}
}

// ParseWhileLoopStmt - yield while loop node
// CFG:
// WhileLoopStmt -> 每当 Expr ：
//               ..     Block
func ParseWhileLoopStmt(p *ParserZH) *syntax.WhileLoopStmt {
	// #1. consume expr
	// 为  as logicYES here
	trueExpr := ParseExpression(p)

	// #2. parse colon
	p.consume(TypeFuncCall)
	// #3. parse block
	expected, blockIndent := p.expectBlockIndent()
	if !expected {
		panic(p.getInvalidSyntaxPeek())
	}
	block := ParseBlockStmt(p, blockIndent)
	return &syntax.WhileLoopStmt{
		TrueExpr:  trueExpr,
		LoopBlock: block,
	}
}

// ParseBlockStmt - parse all statements inside a block
func ParseBlockStmt(p *ParserZH, blockIndent int) *syntax.StmtBlock {
	bStmt := &syntax.StmtBlock{
		Children: []syntax.Statement{},
	}

	// 01. parse all statements
	parseItemListBlock(p, blockIndent, func() {
		stmt := ParseStatement(p)
		bStmt.Children = append(bStmt.Children, stmt)
	})

	return bStmt
}

// ParseBranchStmt - yield BranchStmt node
// CFG:
// CondStmt -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 否则 ：
//         ...     IfFalseBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 再如 OtherExpr1 ：
//         ...     OtherBlock1
//         ... 再如 OtherExpr2 ：
//         ...     OtherBlock2
//         ... ....
//             否则 ：
//         ...     IfFalseBlock
func ParseBranchStmt(p *ParserZH) *syntax.BranchStmt {
	var condExpr syntax.Expression
	var condBlock *syntax.StmtBlock

	var stmt = new(syntax.BranchStmt)

	var condKeywords = []uint8{
		TypeCondElseW,
		TypeCondOtherW,
	}
	// by definition, the first Branch (if-branch) is required,
	// and the 如果 (if) keyword has been consumed before this function call.
	//
	// thus for other branches (like else-branch and elseif-branch),
	// we should consume the corresponding keyword explicitly. (否则，再如)
	const (
		stateInit        = 0
		stateIfBranch    = 1
		stateElseBranch  = 2
		stateOtherBranch = 3
	)

	mainIndent := p.getCurrIndent()
	var hState = stateInit

	for p.peek().Type != TypeEOF {
		// parse header
		switch hState {
		case stateInit:
			hState = stateIfBranch
		case stateIfBranch, stateOtherBranch:
			if p.getPeekIndent() != mainIndent {
				return stmt
			}
			// suppose the if-statement has NOT been finished yet...
			p.unsetStmtCompleteFlag()
			// parse related keywords (如果 expr： , 再如 expr：, 否则：)
			if match, tk := p.tryConsume(condKeywords...); match {
				if tk.Type == TypeCondOtherW {
					hState = stateOtherBranch
				} else {
					hState = stateElseBranch
				}
			} else {
				// NO 再如 orr 否则 matches, the if-statement has REALLY
				// DONE
				p.setStmtCompleteFlag()
				return stmt
			}
		case stateElseBranch:
			if p.getPeekIndent() != mainIndent {
				return stmt
			}
			if match, _ := p.tryConsume(TypeCondElseW); !match {
				return stmt
			}
		}

		// #1. parse condition expr
		if hState != stateElseBranch {
			condExpr = ParseExpression(p)
		}

		// #2. parse colon
		p.consume(TypeFuncCall)

		// #3. parse block statements
		ok, blockIndent := p.expectBlockIndent()
		if !ok {
			panic(p.getUnexpectedIndentPeek())
		}
		condBlock = ParseBlockStmt(p, blockIndent)

		// #4. fill data
		switch hState {
		case stateIfBranch:
			stmt.IfTrueExpr = condExpr
			stmt.IfTrueBlock = condBlock
		case stateOtherBranch:
			stmt.OtherExprs = append(stmt.OtherExprs, condExpr)
			stmt.OtherBlocks = append(stmt.OtherBlocks, condBlock)
		case stateElseBranch:
			stmt.HasElse = true
			stmt.IfFalseBlock = condBlock
			// only one else-branch is accepted
			return stmt
		}
	}
	return stmt
}

// ParseFunctionDeclareStmt - yield FunctionDeclareStmt node
// CFG:
// FunctionDeclareStmt -> 如何 FuncName ？
//       ...     已知 ID1、 & ID2、 ...
//       ...     ExecBlock
//       ... 拦截 ID3：
//       ...     CatchBlock
// 		 ... 拦截 ID4：
//       ...     CatchBlock
//
// FunctionDeclareStmt -> 如何 FuncName ？
//       ...     ExecBlock
//       ...     ....
//
func ParseFunctionDeclareStmt(p *ParserZH) *syntax.FunctionDeclareStmt {
	xID, xParamList, xExecBlock, xCatchBlocks := parseFunctionBlock(p)

	return &syntax.FunctionDeclareStmt{
		Name:        xID,
		DeclareType: syntax.DeclareTypeFunc,
		ParamList:   xParamList,
		ExecBlock:   xExecBlock,
		CatchBlocks: xCatchBlocks,
	}
}

// ParseFunctionDeclareStmt - similiar to ParseFunctionDeclareStmt, but it yields `如何新建XX？`
func ParseConstructorDeclareStmt(p *ParserZH) *syntax.FunctionDeclareStmt {
	xID, xParamList, xExecBlock, xCatchBlocks := parseFunctionBlock(p)

	return &syntax.FunctionDeclareStmt{
		Name:        xID,
		DeclareType: syntax.DeclareTypeConstructor,
		ParamList:   xParamList,
		ExecBlock:   xExecBlock,
		CatchBlocks: xCatchBlocks,
	}
}

func parseFunctionBlock(p *ParserZH) (*syntax.ID, []*syntax.ParamItem, *syntax.StmtBlock, []*syntax.CatchBlockPair) {
	// by definition, when 已知 syntax.Statement exists, it should be at first line
	// of function block
	const (
		stateParamList = 0
		stateFuncBlock = 2
	)
	var hState = stateParamList

	var xID *syntax.ID
	var xParamList []*syntax.ParamItem = make([]*syntax.ParamItem, 0)
	var xExecBlock *syntax.StmtBlock
	var xCatchBlocks []*syntax.CatchBlockPair = make([]*syntax.CatchBlockPair, 0)

	mainIndent := p.getCurrIndent()

	// #1. try to parse ID
	xID = parseFuncID(p)
	// #2. parse question mark
	p.consume(TypeFuncDeclare)

	// #3. parse block manually
	ok, blockIndent := p.expectBlockIndent()
	if !ok {
		panic(p.getUnexpectedIndentPeek())
	}
	// #3.1 parse param def list
	parseItemListBlock(p, blockIndent, func() {
		switch hState {
		case stateParamList:
			// parse 已知 expr
			if match, _ := p.tryConsume(TypeParamAssignW); match {
				parsePauseCommaList(p, func() {
					refMark, _ := p.tryConsume(TypeObjRef)

					idItem := parseID(p)
					xParamList = append(xParamList, &syntax.ParamItem{
						RefMark: refMark,
						ID:      idItem,
					})
				})

				// then change state
				hState = stateFuncBlock
			} else {
				hState = stateFuncBlock
			}
		case stateFuncBlock:
			xExecBlock = ParseBlockStmt(p, blockIndent)
		}
	})
	// #3.2 parse catch blocks
	for p.peek().Type != TypeEOF {
		if p.getPeekIndent() != mainIndent {
			break
		}
		p.unsetStmtCompleteFlag()
		if match, _ := p.tryConsume(TypeCatchErrorW); match {
			catchClass := parseFuncID(p)
			var catchBlock = &syntax.StmtBlock{
				Children: []syntax.Statement{},
			}
			// #2. parse question mark
			p.consume(TypeFuncCall)

			// #3. parse block manually
			ok, blockIndent := p.expectBlockIndent()
			if !ok {
				panic(p.getUnexpectedIndentPeek())
			}

			catchBlock = ParseBlockStmt(p, blockIndent)

			xCatchBlocks = append(xCatchBlocks, &syntax.CatchBlockPair{
				ExceptionClass: catchClass,
				StmtBlock:      catchBlock,
			})
		} else {
			p.setStmtCompleteFlag()
			break
		}
	}

	return xID, xParamList, xExecBlock, xCatchBlocks
}

// ParseGetterDeclareStmt - yield GetterDeclareStmt node
// CFG:
// GetterDeclareStmt -> 何为 GetterName ？
//       ...     ExecBlock
//       ...     ....
//
func ParseGetterDeclareStmt(p *ParserZH) *syntax.FunctionDeclareStmt {
	xID, xParamList, xExecBlock, xCatchBlocks := parseFunctionBlock(p)

	return &syntax.FunctionDeclareStmt{
		DeclareType: syntax.DeclareTypeGetter,
		Name:        xID,
		ParamList:   xParamList,
		ExecBlock:   xExecBlock,
		CatchBlocks: xCatchBlocks,
	}
}

// ParseVarOneLeadStmt -
// There are 2 possible statements
//
// 1. 以 K、V 遍历...
// 2. 以 A （执行方法）
//
// CFG:
//
// VOStmt -> 以 ID  遍历 IStmtT'
// VOStmt -> 以 ID 、 ID  遍历 IStmtT'
//        -> 以 Expr  FuncExprT'
func ParseVarOneLeadStmt(p *ParserZH) syntax.Statement {
	// parse IterateStmt or FuncCallStmt
	exprI := ParseExpression(p)

	if match, tk := p.tryConsume(TypeIteratorW, TypeFuncQuoteL); match {
		switch tk.Type {
		case TypeIteratorW:
			if idX, ok := exprI.(*syntax.ID); ok {
				return parseIteratorStmtRest(p, []*syntax.ID{idX})
			}
			panic(p.getInvalidSyntaxPeek())
		case TypeFuncQuoteL:
			result := &syntax.MemberMethodExpr{
				Root:        exprI,
				MethodChain: []*syntax.FuncCallExpr{},
				YieldResult: nil,
			}

			// parse first function in function chain
			funcExprI := ParseFuncCallExpr(p, false)
			result.MethodChain = append(result.MethodChain, funcExprI)

			// then parse 、 and （  (for chain function)
			for {
				match, _ := p.tryConsume(TypePauseCommaSep)
				if !match {
					break
				}
				// parse （
				p.consume(TypeFuncQuoteL)
				funcExprN := ParseFuncCallExpr(p, false)
				result.MethodChain = append(result.MethodChain, funcExprN)
			}

			// then parse 得到
			if match, _ := p.tryConsume(TypeGetResultW); match {
				id := parseID(p)
				result.YieldResult = id
			}
			return result
		}
	}

	// Another case: 以 ID 、ID 遍历 IStmtT'
	p.consume(TypePauseCommaSep)
	exprII := ParseExpression(p)

	if match2, _ := p.tryConsume(TypeIteratorW); match2 {
		idX, okX := exprI.(*syntax.ID)
		idY, okY := exprII.(*syntax.ID)

		if okX && okY {
			return parseIteratorStmtRest(p, []*syntax.ID{idX, idY})
		}
		panic(p.getInvalidSyntaxPeek())
	}
	panic(p.getInvalidSyntaxPeek())
}

// ParseIteratorStmt - parse iterate stmt that starts with 遍历 keyword only
// CFG:
//
// IStmt -> 遍历 TargetExpr ：  StmtBlock
func ParseIteratorStmt(p *ParserZH) *syntax.IterateStmt {
	return parseIteratorStmtRest(p, []*syntax.ID{})
}

// parseIteratorStmtRest - parse after 以 ... and meet 遍历
// IStmtT'  -> [遍历] TargetExpr ：  StmtBlock
func parseIteratorStmtRest(p *ParserZH, idList []*syntax.ID) *syntax.IterateStmt {
	// 1. parse target expr
	targetExpr := ParseExpression(p)

	// 2. parse colon
	p.consume(TypeFuncCall)

	// 3. parse iterate block
	expected, blockIndent := p.expectBlockIndent()
	if !expected {
		panic(p.getInvalidSyntaxPeek())
	}
	block := ParseBlockStmt(p, blockIndent)

	return &syntax.IterateStmt{
		IterateExpr:  targetExpr,
		IndexNames:   idList,
		IterateBlock: block,
	}
}

// ParseFunctionReturnStmt - yield FuncParamList node (without head token: 输出)
//
// CFG:
// FRStmt -> 输出 syntax.Expression
func ParseFunctionReturnStmt(p *ParserZH) *syntax.FunctionReturnStmt {
	expr := ParseExpression(p)
	return &syntax.FunctionReturnStmt{
		ReturnExpr: expr,
	}
}

func ParseVarInputStmt(p *ParserZH) *syntax.VarInputStmt {
	var idList []*syntax.ID
	parsePauseCommaList(p, func() {
		id := parseID(p)
		idList = append(idList, id)
	})

	return &syntax.VarInputStmt{
		IDList: idList,
	}
}

// ParseThrowExceptionStmt - parse "throw" exception statement
//
// CFG:
//  TEStmt    ->  抛出  ID  ：  Expr   TExpr'  ！
//  TExpr'    ->  、  Expr  TExpr'
//            ->
func ParseThrowExceptionStmt(p *ParserZH) *syntax.ThrowExceptionStmt {
	var exprs []syntax.Expression
	// get id
	exceptionClass := parseID(p)
	p.consume(TypeFuncCall)

	exprI := ParseExpression(p)
	exprs = append(exprs, exprI)

	for {
		if match, _ := p.tryConsume(TypePauseCommaSep); match {
			exprX := ParseExpression(p)
			exprs = append(exprs, exprX)
		} else {
			break
		}
	}

	// parse "！" as last token
	p.consume(TypeExceptionT)

	return &syntax.ThrowExceptionStmt{
		ExceptionClass: exceptionClass,
		Params:         exprs,
	}
}

// ParseImportStmt - parse import syntax.Statement
// CFG:
// ImportStmt  ->  导入 String ImportTail
//
// ImportTail  -> 之 ID IDTail
//             ->
//
// IDTail      -> 、 ID IDTail
//             ->
func ParseImportStmt(p *ParserZH) *syntax.ImportStmt {
	stmt := &syntax.ImportStmt{}
	match, tk := p.tryConsume(TypeLibString, TypeString)
	if !match {
		panic(p.getInvalidSyntaxPeek())
	}

	if tk.Type == TypeLibString {
		stmt.ImportLibType = syntax.LibTypeStd
	} else {
		stmt.ImportLibType = syntax.LibTypeCustom
	}

	stmt.ImportName = newString(p, tk)

	match2, _ := p.tryConsume(TypeObjDotW, TypeObjDotIIW)
	if !match2 {
		return stmt
	}
	// if match 导入 xxx 之 yyy、zzz
	parsePauseCommaList(p, func() {
		tk := parseFuncID(p)
		stmt.ImportItems = append(stmt.ImportItems, tk)
	})

	return stmt
}

// ParseClassDeclareStmt - define class structure
// A typical class may look like this:
//
// 定义 <NAME>：
//    其 <Prop1> 为 <Value1>     <-- PropertyDeclare (for listing all properties with initial value)
//    其 <Prop2> 为 <Value2>
//
//    如何 <Method1> ？    <-- MethodDeclare
//        <Blocks> ...
//        <Blocks> ...
//
//    何为 <Method1> ？    <-- GetterDeclare
//        <Blocks> ...
//        <Blocks> ...
//
// CFG:
// ClassStmt  ->  定义 ClassID ：
//                    ClassDeclareBlock
//
// ClassDeclareBlock  -> ClassDeclareBlockItem1  ClassDeclareBlockItem2 ...
//
// ClassDeclareBlockItem -> Constructor
//                       -> PropertyDeclareStmt
//                       -> FunctionDeclareStmt
//                       -> GetterDeclareStmt
//
func ParseClassDeclareStmt(p *ParserZH) *syntax.ClassDeclareStmt {
	var cdStmt = new(syntax.ClassDeclareStmt)
	// #1. consume ID
	cdStmt.ClassName = parseID(p)

	// #2. parse colon
	p.consume(TypeFuncCall)
	// #3. parse block
	expected, blockIndent := p.expectBlockIndent()
	if !expected {
		panic(p.getInvalidSyntaxPeek())
	}

	// parse block
	parseItemListBlock(p, blockIndent, func() {
		var validChildTypes = []uint8{
			TypeFuncW,
			TypeGetterW,
			TypeObjThisW,
		}
		p.unsetStmtCompleteFlag()

		match, tk := p.tryConsume(validChildTypes...)
		if !match {
			panic(p.getInvalidSyntaxPeek())
		}

		switch tk.Type {
		case TypeFuncW:
			stmt := ParseFunctionDeclareStmt(p)
			cdStmt.MethodList = append(cdStmt.MethodList, stmt)
		case TypeGetterW:
			stmt := ParseGetterDeclareStmt(p)
			cdStmt.GetterList = append(cdStmt.GetterList, stmt)
		case TypeObjThisW:
			stmt := parsePropertyDeclareStmt(p)
			cdStmt.PropertyList = append(cdStmt.PropertyList, stmt)
		}
	})

	return cdStmt
}

// parsePropertyDeclareStmt -
// CFG:
// PropertyDeclareStmt -> 其 ID 为 syntax.Expression
func parsePropertyDeclareStmt(p *ParserZH) *syntax.PropertyDeclareStmt {
	// #1. parse ID
	idItem := parseFuncID(p)
	// consume 为 or =
	p.consume(TypeAssignW, TypeAssignMark)

	// #2. parse expr
	initExpr := ParseExpression(p)

	return &syntax.PropertyDeclareStmt{
		PropertyID: idItem,
		InitValue:  initExpr,
	}
}

// ParseBreakStmt -
// CFG:
// ParseBreakStmt  ->  结束循环
func ParseBreakStmt(p *ParserZH) *syntax.BreakStmt {
	return &syntax.BreakStmt{}
}

// ParseContinueStmt -
// CFG:
// ParseContinueStmt  ->  继续循环
func ParseContinueStmt(p *ParserZH) *syntax.ContinueStmt {
	return &syntax.ContinueStmt{}
}

//// parse helpers
func parseID(p *ParserZH) *syntax.ID {
	match, tk := p.tryConsume(TypeIdentifier)
	if !match {
		panic(p.getInvalidSyntaxPeek())
	}
	return newID(p, tk)
}

// parseFuncID - allow parsing number (as string)
func parseFuncID(p *ParserZH) *syntax.ID {
	match, tk := p.tryConsume(TypeIdentifier)
	if !match {
		panic(p.getInvalidSyntaxPeek())
	}
	return newID(p, tk)
}

// parsePauseCommaList - 使用顿号来分隔
func parsePauseCommaList(p *ParserZH, consumer consumerFunc) {
	// first item MUST be consumed!
	consumer()

	// iterate to get value
	for {
		// consume comma
		if match, _ := p.tryConsume(TypePauseCommaSep); !match {
			// stop parsing immediately
			return
		}
		consumer()
	}
}

func parseItemListBlock(p *ParserZH, blockIndent int, consumer func()) {
	for (p.peek().Type != TypeEOF) && p.getPeekIndent() == blockIndent {
		consumer()
	}
}

func newID(p *ParserZH, tk *syntax.Token) *syntax.ID {
	id := new(syntax.ID)
	id.SetLiteral(tk.Literal)
	p.setStmtCurrentLine(id, tk)
	return id
}

func newString(p *ParserZH, tk *syntax.Token) *syntax.String {
	str := new(syntax.String)
	// remove first char and last char (that are left & right quotes)
	str.SetLiteral(tk.Literal)
	p.setStmtCurrentLine(str, tk)
	return str
}
