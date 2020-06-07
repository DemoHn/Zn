# Zn
_千里之行，始于足下。_


## 简介

`Zn` 是一门 _面向业务_ 的新时代编程语言。

所谓「面向业务」，即是为用户开发业务代码时提供便利，使得用户能够编写出快速、稳定、可维护的代码以实现需求。换言之，和其他「面向计算机」的编程语言不同，`Zn` 更加强调「以人为本」，希望程序能够适应真实世界的需求而不是反过来。

为此，`Zn` 拥有以下独一无二的特性：

- 采用 **中文关键词、标点符号**。这样再也不用为「对专业术语用英文命名变量」之事发愁了。
- 默认使用 **高精度小数** 作为数值并参与运算，杜绝因浮点数计算所带来的计算误差。
  > 这一点对开发金融应用尤为关键。显然，诸如 `0.1 + 0.2 = 0.30000000000000004` 这样的结果在金融应用中是无法忍受的。
- 贴近汉语本身语法，阅读代码可以像阅读文章一样自然。
- 关键词之间不必用空格间断。

## 新手上路

_注：`Zn` 采用 Go 语言开发，安装前请确保 Go 语言编译器已安装在机器上。（建议版本 `>= 1.11`）_

1. 下载并安装 `Zn` 语言：
```sh
# 下载 Zn 语言
$ go get -u github.com/DemoHn/Zn

# 查看 Zn 语言版本，若显示 「Zn语言版本：rv2」或类似字样即表示安装成功
$ Zn -v
```

2. 进入交互执行模式 (REPL)：
  > 注1：下文中出现的直角引号 `「 」` 也可以用普通的双引号 `“ ”` 代替，如：  
  > `令BAT为【“字节”，“阿里”，“腾讯”】`，下同 。
  >
  > 注2：按 `Ctrl + C` 即可退出交互执行模式。
```sh
$ Zn
Zn> 令BAT为【「字节」，「阿里」，「腾讯」】
Zn> 令鹅厂为BAT#2
Zn> （显示：鹅厂）
腾讯
```

3. 执行文件

- 将以下内容输入到「计算均值.zn」文件中：

![code#1.png](/doc/images/code#1.png)

以文件名 `计算均值.zn` 做为 `Zn` 的第一个参数执行命令，并得到结果：

```sh
$ Zn 计算均值.zn
此数据源之均值为 339.36429
```

## 快速入门

1. 声明变量：`令 〔变量名〕 为 〔值〕`

![code#2.png](/doc/images/code#2.png)

变量须先声明，方可继续使用。

2. 对变量赋值：`〔变量名〕 为 〔值〕`

![code#3.png](/doc/images/code#3.png)

`进站人次-累计` 为变量名（`-`也是变量名的一部分）， `1284500` 为值。相当于 `进站人次-累计 := 1284500`

3. 调用方法：`（〔变量名〕：〔参数1〕，...）`

![code#4.png](/doc/images/code#4.png)

Zn 内置了 `显示` 方法，类似于 `console.log` 或者 `print`。  
`X+Y` 表示对所有参数求和；类似地还有 `X-Y`， `X*Y`， `X/Y`。

4. 流程控制：`如果 〔表达式〕：` 及 `每当 〔表达式〕：`

![code#5.png](/doc/images/code#5.png)

Zn 采用类似 Python 的方式，即是以缩进表示语法块。一个单位缩进为 `4个空格` 或者 `1个TAB`。同一个文件里的缩进类型要么都是空格，要么都是TAB，不能混用。

对于 `如果` 语句而言，若是后面的表达式值为 `真`，则执行「如果」后面的语句块；若为 `假`，则执行「否则」后面的语句块。

对于 `每当` 语句而言，当后面的表达式值为 `真`，则下属语句块会循环往复执行。

5. 定义方法：`如何 〔方法名〕 ？`

![code#6.png](/doc/images/code#6.png)

定义一个方法（又称函数）以 `如何` 开始，后面跟着要定义的方法名。而后在子语句块的第一项中定义参数列表：`已知 〔参数1〕，〔参数2〕，...`。 `返回`语句定义了执行方法之后所需的返回值。


## 参与开发

TODO

## 了解更多

TODO