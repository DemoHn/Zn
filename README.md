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
<pre class='zn-ref-ySWxJfAs' style='display: none'>zn: 注1：定义变量
令数据源为【1.24*10^2， 125.75， 225， 3.25e+2， 425， 525， 625.8】

注2：定义方法（也就是函数），通过「已知」语句指定其参数
如何计算均值？
    已知数据源，总数
    令I为0，S为0

    每当I小于总数：
        S为（X+Y：S，数据源#{I}）
        I为（X+Y：I，1）

    返回（X/Y：S，总数）

注3：执行方法
（显示：「此数据源之均值为」，（计算均值：数据源，7））</pre>
<pre class='zn-source-ySWxJfAs' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'><span></span><span style='color: #6a737d'>注1：定义变量</span>
<span></span><span style='color: #d73a49'>令</span>数据源<span style='color: #d73a49'>为</span>【<span style='color: #005cc5'>1.24*10^2</span>，<span>&nbsp;</span><span style='color: #005cc5'>125.75</span>，<span>&nbsp;</span><span style='color: #005cc5'>225</span>，<span>&nbsp;</span><span style='color: #005cc5'>3.25e+2</span>，<span>&nbsp;</span><span style='color: #005cc5'>425</span>，<span>&nbsp;</span><span style='color: #005cc5'>525</span>，<span>&nbsp;</span><span style='color: #005cc5'>625.8</span>】

<span></span><span style='color: #6a737d'>注2：定义方法（也就是函数），通过「已知」语句指定其参数</span>
<span></span><span style='color: #d73a49'>如何</span>计算均值<span style='color: #6f42c1'>？</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>已知</span>数据源，总数
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>令</span>I<span>&nbsp;</span><span style='color: #d73a49'>为</span><span style='color: #005cc5'>0</span>，S<span>&nbsp;</span><span style='color: #d73a49'>为</span><span style='color: #005cc5'>0</span>

<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>每当</span>I<span>&nbsp;</span><span style='color: #d73a49'>小于</span>总数<span style='color: #6f42c1'>：</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>S<span>&nbsp;</span><span style='color: #d73a49'>为</span>（X+Y<span style='color: #6f42c1'>：</span>S<span>&nbsp;</span>，数据源#{I<span>&nbsp;</span>}）
<span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>I<span>&nbsp;</span><span style='color: #d73a49'>为</span>（X+Y<span style='color: #6f42c1'>：</span>I<span>&nbsp;</span>，<span style='color: #005cc5'>1</span>）

<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>返回</span>（X/Y<span style='color: #6f42c1'>：</span>S<span>&nbsp;</span>，总数）

<span></span><span style='color: #6a737d'>注3：执行方法</span>
<span></span>（显示<span style='color: #6f42c1'>：</span><span style='color: #032f62'>「此数据源之均值为」</span>，（计算均值<span style='color: #6f42c1'>：</span>数据源，<span style='color: #005cc5'>7</span>））</pre>

以文件名 `计算均值.zn` 做为 `Zn` 的第一个参数执行命令，并得到结果：

```sh
$ Zn 计算均值.zn
此数据源之均值为 339.36429
```

## 快速入门

1. 声明变量：`令 〔变量名〕 为 〔值〕`
<pre class='zn-ref-eow488VI' style='display: none'>zn: 令PI为3.1415926
令格言为「千里之行，始于足下。」
令等差数列为【1，5，9，13，17，21】</pre>
<pre class='zn-source-eow488VI' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'><span style='color: #d73a49'>令</span>PI<span style='color: #d73a49'>为</span><span style='color: #005cc5'>3.1415926</span>
<span style='color: #d73a49'>令</span>格言<span style='color: #d73a49'>为</span><span style='color: #032f62'>「千里之行，始于足下。」</span>
<span style='color: #d73a49'>令</span>等差数列<span style='color: #d73a49'>为</span>【<span style='color: #005cc5'>1</span>，<span style='color: #005cc5'>5</span>，<span style='color: #005cc5'>9</span>，<span style='color: #005cc5'>13</span>，<span style='color: #005cc5'>17</span>，<span style='color: #005cc5'>21</span>】</pre>
变量须先声明，方可继续使用。

2. 对变量赋值：`〔变量名〕 为 〔值〕`
<pre class='zn-ref-wNh4S2vq' style='display: none'>zn: 进站人次-累计为1284500</pre>
<pre class='zn-source-wNh4S2vq' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'>进站人次-累计<span style='color: #d73a49'>为</span><span style='color: #005cc5'>1284500</span></pre>

`进站人次-累计` 为变量名（`-`也是变量名的一部分）， `1284500` 为值。相当于 `进站人次-累计 := 1284500`

3. 调用方法：`（〔变量名〕：〔参数1〕，...）`
<pre class='zn-ref-mVFmm5dI' style='display: none'>zn: （显示：「千里之行，始于足下。」）
令总和为（X+Y：100，200，300，400）</pre>
<pre class='zn-source-mVFmm5dI' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'>（显示<span style='color: #6f42c1'>：</span><span style='color: #032f62'>「千里之行，始于足下。」</span>）
<span style='color: #d73a49'>令</span>总和<span style='color: #d73a49'>为</span>（X+Y<span style='color: #6f42c1'>：</span><span style='color: #005cc5'>100</span>，<span style='color: #005cc5'>200</span>，<span style='color: #005cc5'>300</span>，<span style='color: #005cc5'>400</span>）</pre>

Zn 内置了 `显示` 方法，类似于 `console.log` 或者 `print`。  
`X+Y` 表示对所有参数求和；类似地还有 `X-Y`， `X*Y`， `X/Y`。

4. 流程控制：`如果 〔表达式〕：` 及 `每当 〔表达式〕：`
<pre class='zn-ref-yAncjwXs' style='display: none'>zn: 令小明的成绩为78
如果小明的成绩大于60：
    （显示：「小明终于及格了。」）
否则：
    （显示：「小明还没有及格。」）

注1：以下流程用于计算 1 + 2 + 3 ... + 100
令I为1
令S为0

每当I不大于100：
    S为（X+Y：S，I）
    I为（X+Y：I，1）

（显示：「1 + 2 + 3 ... + 100 =」，S）</pre>
<pre class='zn-source-yAncjwXs' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'><span></span><span style='color: #d73a49'>令</span>小明的成绩<span style='color: #d73a49'>为</span><span style='color: #005cc5'>78</span>
<span></span><span style='color: #d73a49'>如果</span>小明的成绩<span style='color: #d73a49'>大于</span><span style='color: #005cc5'>60</span><span style='color: #6f42c1'>：</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span>（显示<span style='color: #6f42c1'>：</span><span style='color: #032f62'>「小明终于及格了。」</span>）
<span></span><span style='color: #d73a49'>否则</span><span style='color: #6f42c1'>：</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span>（显示<span style='color: #6f42c1'>：</span><span style='color: #032f62'>「小明还没有及格。」</span>）

<span></span><span style='color: #6a737d'>注1：以下流程用于计算 1 + 2 + 3 ... + 100</span>
<span></span><span style='color: #d73a49'>令</span>I<span>&nbsp;</span><span style='color: #d73a49'>为</span><span style='color: #005cc5'>1</span>
<span></span><span style='color: #d73a49'>令</span>S<span>&nbsp;</span><span style='color: #d73a49'>为</span><span style='color: #005cc5'>0</span>

<span></span><span style='color: #d73a49'>每当</span>I<span>&nbsp;</span><span style='color: #d73a49'>不大于</span><span style='color: #005cc5'>100</span><span style='color: #6f42c1'>：</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span>S<span>&nbsp;</span><span style='color: #d73a49'>为</span>（X+Y<span style='color: #6f42c1'>：</span>S<span>&nbsp;</span>，I<span>&nbsp;</span>）
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span>I<span>&nbsp;</span><span style='color: #d73a49'>为</span>（X+Y<span style='color: #6f42c1'>：</span>I<span>&nbsp;</span>，<span style='color: #005cc5'>1</span>）

<span></span>（显示<span style='color: #6f42c1'>：</span><span style='color: #032f62'>「1 + 2 + 3 ... + 100 =」</span>，S<span>&nbsp;</span>）</pre>

Zn 采用类似 Python 的方式，即是以缩进表示语法块。一个单位缩进为 `4个空格` 或者 `1个TAB`。同一个文件里的缩进类型要么都是空格，要么都是TAB，不能混用。

对于 `如果` 语句而言，若是后面的表达式值为 `真`，则执行「如果」后面的语句块；若为 `假`，则执行「否则」后面的语句块。

对于 `每当` 语句而言，当后面的表达式值为 `真`，则下属语句块会循环往复执行。

5. 定义方法：`如何 〔方法名〕 ？`

<pre class='zn-ref-JItoHFUf' style='display: none'>zn: 如何计算面积-三角形？
    已知高，底边长
    返回（X*Y：高，底边长，0.5）

如何计算面积-正方形？
    已知边长
    返回（X*Y：边长，边长）

如何计算面积-圆形？
    已知直径
    令PI为3.1415926
    令半径为（X*Y：直径，0.5）
    返回（X*Y：PI，半径，半径）

（计算面积-三角形：10，8）</pre>
<pre class='zn-source-JItoHFUf' style='font-family: Sarasa Mono SC, Microsoft YaHei, monospace'><span></span><span style='color: #d73a49'>如何</span>计算面积-三角形<span style='color: #6f42c1'>？</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>已知</span>高<span>&nbsp;</span>，底边长
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>返回</span>（X*Y<span style='color: #6f42c1'>：</span>高<span>&nbsp;</span>，底边长，<span style='color: #005cc5'>0.5</span>）

<span></span><span style='color: #d73a49'>如何</span>计算面积-正方形<span style='color: #6f42c1'>？</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>已知</span>边长
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>返回</span>（X*Y<span style='color: #6f42c1'>：</span>边长，边长）

<span></span><span style='color: #d73a49'>如何</span>计算面积-圆形<span style='color: #6f42c1'>？</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>已知</span>直径
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>令</span>PI<span style='color: #d73a49'>为</span><span style='color: #005cc5'>3.1415926</span>
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>令</span>半径<span style='color: #d73a49'>为</span>（X*Y<span style='color: #6f42c1'>：</span>直径，<span style='color: #005cc5'>0.5</span>）
<span>&nbsp;&nbsp;&nbsp;&nbsp;</span><span style='color: #d73a49'>返回</span>（X*Y<span style='color: #6f42c1'>：</span>PI，半径，半径）

<span></span>（计算面积-三角形<span style='color: #6f42c1'>：</span><span style='color: #005cc5'>10</span>，<span style='color: #005cc5'>8</span>）</pre>

定义一个方法（又称函数）以 `如何` 开始，后面跟着要定义的方法名。而后在子语句块的第一项中定义参数列表：`已知 〔参数1〕，〔参数2〕，...`。 `返回`语句定义了执行方法之后所需的返回值。


## 参与开发

TODO

## 了解更多

TODO