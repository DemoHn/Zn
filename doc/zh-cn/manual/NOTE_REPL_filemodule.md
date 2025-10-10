### 不同的加载模块的方式

1. 读取文件夹
    - 需要 rootDir
    - 一开始就要指定一个文件作为初始模块
    - 根据 “AA-BB-CC” -> <rootDir>/AA/BB/CC.zn 的规则找到的对应的程序文件
        - 可以自定义一个 Finder 函数

2. 单模块模式（playground）
    - 没有 rootDir，只有一个主模块
    - 初始先读取程序，得到 lexer，再把lexer 放进 module 里面去初始化

3. REPL (单模块，lexer 初始值还是 nil)  -> 直接运行 `./zinc` 就进入REPL模式
    - 没有 rootDir，只有一个主模块
    - 第一次执行时读取程序，得到 lexer，再把 lexer 放进 Module 里面去初始化
    - 第二次执行时再得到一个新 lexer ，把当前 lexer 的


### use zinc as Go Library

```go
package main
import .../zinc 

func main() {
    interpreter := zinc.NewInterpreter()

    // exec path 1 -> for code snippets (main module ONLY) & REPL
    interpreter.ExecuteScript(script string, varInput map) -> (zinc.Element, error)

    // exec path 2 -> exeucte file in local file system
    interpreter.ExecuteFile(fileName string, varInput map) -> (zinc.Element, error)

    // exec path 3 -> custom module finder & execute
    interpreter.SetModuleFinder((moudleName string) -> string)
    interpreter.Execute(varInput map) -> (zinc.Element, error)

    //// if you want to run different modules many times, JUST
    //// create a new interpreter!

    interpreter2 := zinc.NewInterpreter()
    interpreter2.ExecuteFile(fileName string) -> (zinc.Element, error)

    ///// transform an element to a string for display usage 
    zinc.ElementToString(elem Element) -> string
}
```