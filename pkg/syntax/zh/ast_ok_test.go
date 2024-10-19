package zh

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
)

var testSuccessSuites = []string{
	varDeclCasesOK,
	whileLoopCasesOK,
	logicExprCasesOK,
	arrayListCasesOK,
	funcCallCasesOK,
	branchStmtCasesOK,
	stmtLineBreakCasesOK,
	memberExprCasesOK,
	memberMethodStmtCasesOK,
	iterateCasesOK,
	classDeclareCasesOK,
	functionDeclareCasesOK,
	importStmtCasesOK,
}

const logicExprCasesOK = `
========
1. low -> high precedence
--------
{A且B或C且D等于E且F设为100}等于0
--------
$PG($X(I=() S=($BK(
	$EQ(
		L=($OR(
				L=($AND(L=($ID(A)) R=($ID(B))))
				R=($AND(					
					L=($AND(
						L=($ID(C))
						R=($EQ(
							L=($ID(D))
							R=($ID(E))
						))
					))
					R=($VA(
						target=($ID(F))
						assign=($ID(100))
					))
				))
		))
		R=($ID(0))
	))) C=()
))
========
2. A or B or C or D or E
--------
A或B或C或D或E
--------
$PG($X(I=() S=($BK(
	$OR(L=(
		$OR(L=(
			$OR(L=(
				$OR(L=($ID(A)) R=(
					$ID(B)
				))
			) R=($ID(C)))
		) R=($ID(D)))
	) R=($ID(E)))
	)) C=()
))
`

const whileLoopCasesOK = `
========
1. one line block
--------
每当1：
	令A设为B
--------
$PG($X(I=() S=($BK(
	$WL(
		expr=($ID(1))
		block=($BK($VD($VP(
				vars[]=($ID(A))
				expr[]=($ID(B))
		))))
	)))
	C=()
))

========
2. nested while loop statement
--------
每当1：
	A设为B
	每当2：
		C设为D
		E设为F
	每当3：
		100
	G设为H
	K设为L

M设为N
--------
$PG($X(I=() S=($BK(
	$WL(
		expr=($ID(1))
		block=($BK(
			$VA(target=($ID(A)) assign=($ID(B)))
			$WL(
				expr=($ID(2))
				block=($BK(
					$VA(target=($ID(C)) assign=($ID(D)))
					$VA(target=($ID(E)) assign=($ID(F)))
				))
			)
			$WL(
				expr=($ID(3))
				block=($BK($ID(100)))
			)
			$VA(target=($ID(G)) assign=($ID(H)))
			$VA(target=($ID(K)) assign=($ID(L)))
		))
	)

	$VA(target=($ID(M)) assign=($ID(N)))
	)) C=()
))
`

const varDeclCasesOK = `
========
1. inline one var
--------
令某变量设为100
--------
$PG($X(I=() S=($BK(
	$VD($VP(
		vars[]=($ID(某变量))
		expr[]=($ID(100))
	))))
	C=()
))
========
1a. inline one var: =
--------
令某变量 = 100
--------
$PG($X(I=() S=($BK(
	$VD($VP(
		vars[]=($ID(某变量))
		expr[]=($ID(100))
	))))
	C=()
))

========
2. two variables
--------
令变量1、变量2设为100
--------
$PG($X(I=() S=(
	$BK($VD($VP(
		vars[]=($ID(变量1) $ID(变量2))
		expr[]=($ID(100))
	))))
	C=()
))

========
3. paired variables inline (one pair only)
--------
令小A、小B、小C设为100
--------
$PG($X(I=() S=(
	$BK($VD($VP(
		vars[]=($ID(小A) $ID(小B) $ID(小C))
		expr[]=($ID(100))
	))))
	C=()
))

========
4. with varquotes
--------
令小A、` + "`" + `先令` + "`" + `设为200
--------
$PG($X(I=() S=($BK(
	$VD($VP(
		vars[]=($ID(小A) $ID(先令))
		expr[]=($ID(200))
	))))
	C=()
))

========
5. A -> B -> C
--------
令A设为B=C
--------
$PG($X(I=() S=($BK(
	$VD($VP(
		vars[]=($ID(A))
		expr[]=(
			$VA(
				target=($ID(B))
				assign=($ID(C))
			)
		)
	))))
	C=()
))

========
6. block var declare
--------
令：
	A设为1
	B设为2
	C、D设为3
	E、F设为4

令G设为5
--------
$PG($X(I=() S=($BK(
	$VD(
		$VP(vars[]=($ID(A))		expr[]=($ID(1)))
		$VP(vars[]=($ID(B))		expr[]=($ID(2)))
		$VP(vars[]=($ID(C) $ID(D))		expr[]=($ID(3)))
		$VP(vars[]=($ID(E) $ID(F))		expr[]=($ID(4)))
	)
	$VD($VP(
		vars[]=($ID(G))
		expr[]=($ID(5))
	))))
	C=()
))

========
7. define const variables
--------
令圆周率恒为3.1415926
--------
$PG($X(I=() S=($BK(
	$VD(
		$VP(const vars[]=($ID(圆周率)) expr[]=($ID(3.1415926)))
	)))
	C=()
))

========
8. define new object
--------
令圆周率 =（新建数学：3.1415926）
--------
$PG($X(I=() S=($BK($VD($VP(vars[]=($ID(圆周率)) expr[]=($NEW(class=($ID(数学)) params=($ID(3.1415926)))))))) C=()))
========
9. block declaration - mixture of const,assign,newObj
--------
令：
	高脚杯、小盅 =（新建SKU：「玻璃制品」、10、20、30）
	A、B、C设为「Amazon」
	D、E、F恒为空
	G恒为空
--------
$PG($X(I=() S=($BK($VD(
	$VP(vars[]=($ID(高脚杯) $ID(小盅)) expr[]=(
		$NEW(class=($ID(SKU)) params=($STR(玻璃制品) $ID(10) $ID(20) $ID(30))))
	)
	$VP(vars[]=($ID(A) $ID(B) $ID(C)) expr[]=($STR(Amazon)))
	$VP(const vars[]=($ID(D) $ID(E) $ID(F)) expr[]=($ID(空)))
	$VP(const vars[]=($ID(G)) expr[]=($ID(空))))))
	C=()
))

========
10. block declaration - new object without params
--------
令A =（新建B）
--------
$PG($X(I=() S=($BK($VD($VP(vars[]=($ID(A)) expr[]=($NEW(class=($ID(B)) params=())))))) C=()))
`
const funcCallCasesOK = `
========
1. success func call with no param
--------
（显示当前时间）
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前时间)) params=())
)) C=()
))
========
2. success func call with no param (varquote)
--------
（` + "`" + `显示当前之时间` + "`" + `）
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前之时间)) params=())
)) C=()
))

========
3. success func call with 1 parameter
--------
（显示当前时间：「今天」）
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(今天)))
)) C=()
))

========
4. success func call with 2 parameters
--------
（显示当前时间：「今天」、「15:30」）
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(今天) $STR(15:30)))
)) C=()
))

========
5. success func call with mutliple parameters
--------
（显示当前时间：「今天」、「15:30」、200、3000）
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(今天) $STR(15:30) $ID(200) $ID(3000)))
)) C=()
))

========
6. nested functions
--------
（显示当前时间：「今天」、「15:30」、（显示时刻））
--------
$PG($X(I=() S=($BK(
	$FN(name=($ID(显示当前时间)) params=(
		$STR(今天)
		$STR(15:30) 
		$FN(name=($ID(显示时刻)) params=())
	))
)) C=()
))
`

const branchStmtCasesOK = `
========
1. if-block only
--------
如果真：
    （X+Y：20、30）
--------
$PG($X(I=() S=($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(
			$FN(
				name=($ID(X+Y))
				params=($ID(20) $ID(30))
			)
		))
	)
)) C=()
))

========
2. if-block and else-block
--------
如果真：
    （X+Y：20、30）
否则：
    （X-Y：20、30）
--------
$PG($X(I=() S=($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(			
			$FN(
				name=($ID(X+Y))
				params=($ID(20) $ID(30))
			)
		))
		elseBlock=($BK(			
			$FN(
				name=($ID(X-Y))
				params=($ID(20) $ID(30))
			)
		))
	)
)) C=()
))

========
3. if-block & elseif blocks
--------
如果真：
    （X+Y：20、30）
再如A等于200：
    （X-Y：20、30）
再如A等于300：
    B设为10；
注：「‘这是一个多行注释’」
如果C /= 真：
    （ASDF）
--------
$PG($X(I=() S=($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(
			$FN(
				name=($ID(X+Y))
				params=($ID(20) $ID(30))
			)
		))
		otherExpr[]=($EQ(
			L=($ID(A))
			R=($ID(200))
		))
		otherBlock[]=($BK(
			$FN(
				name=($ID(X-Y))
				params=($ID(20) $ID(30))
			)
		))
		otherExpr[]=($EQ(
			L=($ID(A))
			R=($ID(300))
		))
		otherBlock[]=($BK(
			$VA(
				target=($ID(B))
				assign=($ID(10))
			)
		$))
	)
	$IF(
		ifExpr=($NEQ(L=($ID(C)) R=($ID(真))))
		ifBlock=($BK(
			$FN(name=($ID(ASDF)) params=())
		))
	)
)) C=()
))

========
4. if-elseif-else block
--------
如果真：
    （X+Y：20、30）
再如A == 100：
    （显示）
否则：
    （X-Y：20、30）
--------
$PG($X(I=() S=($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(			
			$FN(
				name=($ID(X+Y))
				params=($ID(20) $ID(30))
			)
		))
		elseBlock=($BK(			
			$FN(
				name=($ID(X-Y))
				params=($ID(20) $ID(30))
			)
		))
		otherExpr[]=(
			$EQ(L=($ID(A)) R=($ID(100)))
		)
		otherBlock[]=($BK(
			$FN(
				name=($ID(显示))
				params=()
			)
		))		
	)
)) C=()
))

========
5. if block - 为 as compar instead of var assignment
--------
如果A == 100：
    （X+Y：20、30）
再如B == 200：
    （显示）
--------
$PG($X(I=() S=($BK(
	$IF(
		ifExpr=($EQ(L=($ID(A)) R=($ID(100))))
		ifBlock=($BK(			
			$FN(
				name=($ID(X+Y))
				params=($ID(20) $ID(30))
			)
		))
		otherExpr[]=(
			$EQ(L=($ID(B)) R=($ID(200)))
		)
		otherBlock[]=($BK(
			$FN(
				name=($ID(显示))
				params=()
			)
		))
	)
)) C=()
))
`

const arrayListCasesOK = `
========
1. empty array
--------
【】
--------
$PG($X(I=() S=($BK($ARR())) C=()))

========
2. empty hashmap
--------
【=】
--------
$PG($X(I=() S=($BK($HM())) C=()))

========
3. mixed string and decimal array
--------
【「MacBook Air 12"」，2080，3000】
--------
$PG($X(I=() S=($BK($ARR($STR(MacBook Air 12") $ID(2080) $ID(3000)))) C=()))

========
4. array with newline
--------
【
    「MacBook Air 12"」，
    2080，
    3000，
】
--------
$PG($X(I=() S=($BK($ARR($STR(MacBook Air 12") $ID(2080) $ID(3000)))) C=()))

========
5. array nest with array
--------
【
    「MacBook Air 12"」，
    2080，
    【100，200，300】
】
--------
$PG($X(I=() S=($BK(
	$ARR(
	$STR(MacBook Air 12") 
	$ID(2080) 
	$ARR($ID(100) $ID(200) $ID(300))
	))) C=()
))

========
6. array nest with array nest with array
--------
【
    「MacBook Air 12"」，
    2080，
    【100，200，300，
        【
            10000
        】
    】
】
--------
$PG($X(I=() S=($BK(
	$ARR(
	$STR(MacBook Air 12") 
	$ID(2080) 
	$ARR($ID(100) $ID(200) $ID(300) $ARR($ID(10000)))
	))) C=()
))

========
7. a simple hashmap
--------
【
		「数学」 = 80，
		「语文」 = 90
】
--------
$PG($X(I=() S=($BK(
	$HM(
	key[]=($STR(数学)) value[]=($ID(80)) 
	key[]=($STR(语文)) value[]=($ID(90))
	))) C=()
))

========
8. a hashmap nest with hashmap
--------
【
		「数学」 = 80，
		「语文」 = 【
				「阅读」 = 20，
				「听力」 = 30.5，
				「比例」 = 0.12345
		】
】
--------
$PG($X(I=() S=($BK(
	$HM(
	key[]=($STR(数学)) value[]=($ID(80)) 
	key[]=($STR(语文)) value[]=($HM(
		key[]=($STR(阅读)) value[]=($ID(20))
		key[]=($STR(听力)) value[]=($ID(30.5))
		key[]=($STR(比例)) value[]=($ID(0.12345))
	))
	))) 
	C=()
))
`

const stmtLineBreakCasesOK = `
========
1. a statement in oneline
--------
令香港记者设为记者名设为「张宝华」
--------
$PG($X(I=() S=($BK($VD($VP(vars[]=($ID(香港记者)) expr[]=($VA(target=($ID(记者名)) assign=($STR(张宝华)))))))) C=()))

========
2. a complete statement with comma list - 3 lines
--------
令树叶、鲜花、
    雪花、
                墨水设为「黑」
--------
$PG($X(I=() S=($BK($VD($VP(vars[]=($ID(树叶) $ID(鲜花) $ID(雪花) $ID(墨水)) expr[]=($STR(黑)))))) C=()))

========
3. nested function calls with multiple lines
--------
（显示：
    「1」、（调用参数：200、300、
        4000、5000））
--------
$PG($X(I=() S=($BK($FN(name=($ID(显示)) params=($STR(1) $FN(name=($ID(调用参数)) params=($ID(200) $ID(300) $ID(4000) $ID(5000))))))) C=()))

========
4. multi-line hashmap
--------
令对象表设为【
		1 = 「象」，
		2 = 「士」，
		3 = 「车」
】
--------
$PG($X(I=() S=($BK($VD($VP(vars[]=($ID(对象表)) expr[]=($HM(key[]=($ID(1)) value[]=($STR(象)) key[]=($ID(2)) value[]=($STR(士)) key[]=($ID(3)) value[]=($STR(车)))))))) C=()))
`

const memberExprCasesOK = `
========
0. normal arith expr
--------
{5 + 8} * 3 - 7 / 4
--------
$PG($X(I=() S=($BK(
	$AR(type=(SUB) left=(
		$AR(type=(MUL) left=(
			$AR(type=(ADD) left=(
				$ID(5)) right=($ID(8))
			)) right=($ID(3)))
		) right=($AR(type=(DIV) left=($ID(7)) right=($ID(4)))))
)) C=()
))
========
1. normal dot member
--------
天之涯
--------
$PG($X(I=() S=($BK(
	$MB(root=($ID(天)) type=(mID) object=($ID(涯)))
)) C=()
))

========
2. normal dot member (nested)
--------
雪花之天涯之海角
--------
$PG($X(I=() S=($BK(
	$MB(
		root=(
			$MB(root=($ID(雪花)) type=(mID) object=($ID(天涯)))
		)
		type=(mID)
		object=($ID(海角))
	)
)) C=()
))

========
4. array index
--------
Array#123
--------
$PG($X(I=() S=($BK(
	$MB(root=($ID(Array)) type=(mIndex) object=($ID(123)))
)) C=()
))

========
5. array index (using {})
--------
Array#{天之涯}
--------
$PG($X(I=() S=($BK(
	$MB(root=($ID(Array)) type=(mIndex) object=(
		$MB(root=($ID(天)) type=(mID) object=($ID(涯)))
	))
)) C=()
))

========
6. array index (nested)
--------
Array#20#30#{QR}
--------
$PG($X(I=() S=($BK(
	$MB(
		root=(
			$MB(
				root=(
					$MB(
						root=($ID(Array))
						type=(mIndex)
						object=($ID(20))
					)
				)
				type=(mIndex)
				object=($ID(30))
			)
		)
		type=(mIndex)
		object=($ID(QR))
	)
)) C=()
))

========
7. mix methods & members & indexes
--------
Array#10之首
--------
$PG($X(I=() S=($BK(
	$MB(
		root=(
			$MB(
				root=($ID(Array))
				type=(mIndex)
				object=($ID(10))
			)
		)
		type=(mID)
		object=($ID(首))
	)
)) C=()
))

========
9. self root (rootProp)
--------
其年龄设为20
--------
$PG($X(I=() S=($BK(
	$VA(
		target=($MB(
			rootProp
			type=(mID)
			object=($ID(年龄))
		))		
		assign=($ID(20))
	)
)) C=()
))
========
10. mix rootProp and member
--------
其年龄之文本
--------
$PG($X(I=() S=($BK(
	$MB(
		root=(
			$MB(
				rootProp
				type=(mID)
				object=($ID(年龄))
			)
		)
		type=(mID)
		object=($ID(文本))
	)
)) C=()
))
`

const importStmtCasesOK = `
========
1. normal import
--------
导入《对象》
--------
$PG($IM(name=($STR(对象)) items=()))

========
2. normal import with items
--------
导入《对象》的名称、内容
--------
$PG($IM(name=($STR(对象)) items=($ID(名称) $ID(内容))))
`

const memberMethodStmtCasesOK = `
========
1. normal member method expr
--------
以A（运行方法）
--------
$PG($X(I=() S=($BK(
	$MMF(root=($ID(A)) chain=($FN(name=($ID(运行方法)) params=())))
)) C=()
))
========
2. normal member method expr with yield
--------
以A（运行方法），得到结果
--------
$PG($X(I=() S=($BK(
	$MMF(root=($ID(A)) chain=($FN(name=($ID(运行方法)) params=())) yield=($ID(结果)))
)) C=()
))
========
3. normal member method chain expr with yield
--------
以A（运行方法）、（方法2：A、B、C），得到结果
--------
$PG($X(I=() S=($BK(
	$MMF(root=($ID(A)) chain=(
		$FN(name=($ID(运行方法)) params=())
		$FN(name=($ID(方法2)) params=($ID(A) $ID(B) $ID(C)))
	) yield=($ID(结果)))
)) C=()
))
========
4. normal member method chain expr w/o yield
--------
以A（运行方法）、（方法2：A、B、C）
--------
$PG($X(I=() S=($BK(
	$MMF(root=($ID(A)) chain=(
		$FN(name=($ID(运行方法)) params=())
		$FN(name=($ID(方法2)) params=($ID(A) $ID(B) $ID(C)))
	))
)) C=()
))
========
5. normal member method chain more exprs
--------
以A（运行方法）、
	（方法2：A、B、C）、
	（QAQ：1、3、5、7）
--------
$PG($X(I=() S=($BK(
	$MMF(root=($ID(A)) chain=(
		$FN(name=($ID(运行方法)) params=())
		$FN(name=($ID(方法2)) params=($ID(A) $ID(B) $ID(C)))
		$FN(name=($ID(QAQ)) params=($ID(1) $ID(3) $ID(5) $ID(7)))
	))
)) C=()
))
`

const iterateCasesOK = `
========
1. normal iterate expr
--------
遍历【1，2，3】：
    令A设为值
    结束循环
--------
$PG($X(I=() S=($BK(
	$IT(
		target=($ARR($ID(1) $ID(2) $ID(3)))
		idxList=()
		block=($BK(
			$VD($VP(vars[]=($ID(A)) expr[]=(
			  	$ID(值)
			)))
			$BREAK
		))
	)
)) C=()
))

========
2. lead one var
--------
以K遍历代码：
    （显示：K）
--------
$PG($X(I=() S=($BK(
	$IT(
		target=($ID(代码))
		idxList=($ID(K))
		block=($BK(
			$FN(name=($ID(显示)) params=($ID(K)))
		))
	)
)) C=()
))
========
3. lead two vars
--------
以K、V遍历【
		「A」 = 1，
		「B」 = 2，
		「C」 = 3
】：
	（显示：K、V）
--------
$PG($X(I=() S=($BK(
	$IT(
		target=($HM(
			key[]=($STR(A)) value[]=($ID(1))
			key[]=($STR(B)) value[]=($ID(2))
			key[]=($STR(C)) value[]=($ID(3))
		))
		idxList=($ID(K) $ID(V))
		block=($BK($FN(
			name=($ID(显示))
			params=($ID(K) $ID(V))
		)))
	)
)) C=()
))
`

const classDeclareCasesOK = `
========
1. simplist class definition
--------
定义狗：
	其名设为“小黄”
	其品种设为“拉布拉多”
--------
$PG($X(I=() S=($BK(
	$CLS(
		name=($ID(狗))
		properties=(
			$PD(id=($ID(名)) expr=($STR(小黄)))
			$PD(id=($ID(品种)) expr=($STR(拉布拉多)))
		)		
		methods=()
		getters=()
	)
)) C=()
))

========
3. full class definition
--------
定义狗：
	其名设为“小黄”
	其年龄设为0

	如何狂吠？
		输出“汪汪汪”

	如何添加年龄？
		输出20

	何为总和？
		输出20
--------
$PG($X(I=() S=($BK(
	$CLS(
		name=($ID(狗))
		properties=(
			$PD(id=($ID(名)) expr=($STR(小黄)))
			$PD(id=($ID(年龄)) expr=($ID(0)))
		)
		methods=(
			$FN(
				type=FN
				name=($ID(狂吠))
				block=($X(I=() S=($BK(
					$RT($STR(汪汪汪))
				)) C=()))
			)
			$FN(
				type=FN
				name=($ID(添加年龄))
				block=($X(I=() S=($BK(
					$RT($ID(20))
				)) C=()))
			)
		)
		getters=(
			$FN(
				type=GET
				name=($ID(总和))
				block=($X(I=() S=($BK(
					$RT($ID(20))
				)) C=()))
			)
		)
	)
)) C=()
))
========
4. class definition with comment
--------
定义狗：
	注1：定义属性列表、并它们以默认值
	其名设为“小黄”
	其年龄设为0

	注2：方法列表
	如何狂吠？
		注：在方法里面添加注释
		输出“汪汪汪”

	如何添加年龄？
		输出20

	注3：getter列表
	何为总和？
		输出20
--------
$PG($X(I=() S=($BK(
	$CLS(
		name=($ID(狗))
		properties=(
			$PD(id=($ID(名)) expr=($STR(小黄)))
			$PD(id=($ID(年龄)) expr=($ID(0)))
		)
		methods=(
			$FN(
				type=FN
				name=($ID(狂吠))
				block=($X(I=() S=($BK(
					$RT($STR(汪汪汪))
				)) C=()))
			)
			$FN(
				type=FN
				name=($ID(添加年龄))
				block=($X(I=() S=($BK(
					$RT($ID(20))
				)) C=()))
			)
		)
		getters=(
			$FN(
				type=GET
				name=($ID(总和))
				block=($X(I=() S=($BK(
					$RT($ID(20))
				)) C=()))
			)
		)
	)
)) C=()
))
`

const functionDeclareCasesOK = `
========
1. simplist function
--------
如何搞个大新闻？
	1024
--------
$PG($X(I=() S=($BK(
	$FN(
		type=FN
		name=($ID(搞个大新闻))
		block=($X(
			I=()
			S=($BK($ID(1024)))
			C=()
		))
	)))
	C=()
))
========
2. with one param
--------
如何搞个大新闻？
	输入变量1
	1024
--------
$PG($X(I=() S=($BK(
	$FN(
		type=FN
		name=($ID(搞个大新闻))
		block=($X(
			I=($ID(变量1))
			S=($BK($ID(1024)))
			C=()
		))
	)))
	C=()
))
========
3. with multiple params
--------
如何搞个大新闻？
	输入A、B、` + "`" + `华为手机` + "`" + `
	1024
--------
$PG($X(I=() S=($BK(
	$FN(
		type=FN
		name=($ID(搞个大新闻))
		block=($X(
			I=($ID(A) $ID(B) $ID(华为手机))
			S=($BK($ID(1024)))
			C=()
		))
	)))
	C=()
))
========
4. with catch block
--------
如何搞个大新闻？
	如果C == 空：
		A = A * 2
		输出1024
	否则：
		输出1024	
	拦截A异常：
		输出233
	拦截B异常：
		输出566
--------
$PG($X(I=() S=($BK(
	$FN(
		type=FN
		name=($ID(搞个大新闻))
		block=($X(I=() S=(
		$BK(
			$IF(
				ifExpr=($EQ(L=($ID(C)) R=($ID(空))))
				ifBlock=($BK($VA(target=($ID(A)) assign=($AR(type=(MUL) left=($ID(A)) right=($ID(2))))) $RT($ID(1024))))
				elseBlock=($BK($RT($ID(1024))))
			)
		)
		) C=(
			cls[]=($ID(A异常)) stmt[]=($BK(
				$RT($ID(233))
			))
			cls[]=($ID(B异常)) stmt[]=($BK(
				$RT($ID(566))
			))
		)))
	)))
	C=()
))
`

// ////// BY FUNC ////////
// test ParseProgram() only
const testProgramOKCases = `
=========
1. normal case with both importBlock & execBlock
---------
导入《X》

输出666
---------
$PG($IM(name=($STR(X)) items=()) $X(I=() S=($BK($RT($ID(666)))) C=()))
=========
2. only statements
---------
输出668
---------
$PG($X(I=() S=($BK($RT($ID(668)))) C=()))
`

const testExecBlockOKCases = `
=========
1. normal case with both importBlock & execBlock
---------
输入A、B、C

令X = A + B + C

拦截异常1：
	输出123
拦截异常2：
	输出456
---------
$X(
	I=($ID(A) $ID(B) $ID(C)) 
	
	S=($BK($VD($VP(vars[]=($ID(X)) expr[]=($AR(type=(ADD) left=($AR(type=(ADD) left=($ID(A)) right=($ID(B)))) right=($ID(C)))))))) 
	
	C=( cls[]=($ID(异常1)) stmt[]=($BK($RT($ID(123)))) 
		cls[]=($ID(异常2)) stmt[]=($BK($RT($ID(456))))
	)
)
==========
2. normal case with importBlock & stmtBlock only
---------
输入A、B、C

233
---------
$X(
	I=($ID(A) $ID(B) $ID(C)) 
	
	S=($BK($ID(233))) 
	
	C=()
)
==========
3. normal case with stmtBlock only
---------
输出233
---------
$X(
	I=()
	S=($BK($RT($ID(233))))
	C=()
)
`

var testByFuncCaseList = []struct {
	cases   string
	astFunc func(*ParserZH) syntax.Node
}{
	{
		cases: testProgramOKCases,
		astFunc: func(pz *ParserZH) syntax.Node {
			return ParseProgram(pz)
		},
	},
	{
		cases: testExecBlockOKCases,
		astFunc: func(pz *ParserZH) syntax.Node {
			return ParseExecBlock(pz, 0)
		},
	},
}

type astSuccessCase struct {
	name    string
	input   string
	astTree string
}

func TestAST_OK(t *testing.T) {
	astCases := []astSuccessCase{}

	for _, suData := range testSuccessSuites {
		suites := splitTestSuites(suData)
		for _, suite := range suites {
			astCases = append(astCases, astSuccessCase{
				name:    suite[0],
				input:   suite[1],
				astTree: suite[2],
			})
		}
	}

	for _, tt := range astCases {
		t.Run(tt.name, func(t *testing.T) {
			l := syntax.NewLexer([]rune(tt.input))
			p := NewParserZH()

			pg, err := p.ParseAST(l)
			if err != nil {
				t.Errorf("expect no error, got error: %s", err)
			} else {
				// compare with ast
				expect := syntax.StringifyAST(pg)
				got := formatASTstr(tt.astTree)

				if expect != got {
					t.Errorf("AST compare:\nexpect ->\n%s\ngot ->\n%s", expect, got)
				}
			}
		})
	}
}

func TestAST_WithFunc_OK(t *testing.T) {
	var errX error
	defer func() {
		if r := recover(); r != nil {
			errX, _ = r.(error)
		}
	}()

	for _, suData := range testByFuncCaseList {
		astCases := []astSuccessCase{}

		suites := splitTestSuites(suData.cases)
		for _, suite := range suites {
			astCases = append(astCases, astSuccessCase{
				name:    suite[0],
				input:   suite[1],
				astTree: suite[2],
			})
		}

		for _, tt := range astCases {
			t.Run(tt.name, func(t *testing.T) {
				l := syntax.NewLexer([]rune(tt.input))
				p := NewParserZH()
				p.Lexer = l
				// advance token ONCE
				p.next()

				astFuncX := suData.astFunc
				node := astFuncX(p)
				if errX != nil {
					t.Errorf("expect no error, got error: %s", errX)
				} else {
					// compare with ast
					expect := syntax.StringifyAST(node)
					got := formatASTstr(tt.astTree)

					if expect != got {
						t.Errorf("AST compare:\nexpect ->\n%s\ngot ->\n%s", expect, got)
					}
				}
			})
		}
	}
}

func splitTestSuites(source string) [][3]string {
	result := [][3]string{}

	source = strings.Replace(source, "\r\n", "\n", -1)
	sourceArr := strings.Split(source, "\n")

	const (
		sInit    = 0
		sPartI   = 1
		sPartII  = 2
		sPartIII = 3
	)
	var state = sInit
	l1 := []string{}
	l2 := []string{}
	l3 := []string{}
	for _, line := range sourceArr {
		if strings.HasPrefix(line, "========") {
			// push old data
			if state == sPartIII {
				result = append(result, [3]string{
					strings.Join(l1, "\n"),
					strings.Join(l2, "\n"),
					strings.Join(l3, "\n"),
				})
			}
			state = sPartI
			// clear buffer
			l1 = []string{}
			l2 = []string{}
			l3 = []string{}
			continue
		}
		if strings.HasPrefix(line, "--------") {
			if state == sPartI {
				state = sPartII
			} else if state == sPartII {
				state = sPartIII
			}
			continue
		}

		switch state {
		case sPartI:
			l1 = append(l1, line)
		case sPartII:
			l2 = append(l2, line)
		case sPartIII:
			l3 = append(l3, line)
		}
	}

	// tail append
	if state == sPartIII {
		result = append(result, [3]string{
			strings.Join(l1, "\n"),
			strings.Join(l2, "\n"),
			strings.Join(l3, "\n"),
		})
	}

	return result
}

func formatASTstr(input string) string {
	reL := regexp.MustCompile(`\((\s)+`)
	reR := regexp.MustCompile(`(\s)+\)`)
	reS := regexp.MustCompile(`(\s)+`)

	input = reL.ReplaceAllString(input, "(")
	input = reR.ReplaceAllString(input, ")")
	input = reS.ReplaceAllString(input, " ")

	return strings.TrimSpace(input)
}
