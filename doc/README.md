> 千里之行，始于足下。

## 简介

Zn 是一门 _面向业务_ 的编程语言：即使您不是专业的程序员，也可以通过简洁的语法，丰富的方法库来完成日常的业务需求。

和市面上流行的编程语言相比，Zn 拥有以下优势：

- 采用 **中文关键词及标点**，无需绞尽脑汁将业务术语「翻译」成英文命名。
- 贴近汉语自身语法，阅读程序可以像阅读文章一样自然。
- 丰富的方法库，可以做到「开箱即用」，最大程度上满足您的需求。
- 完整的官方中文文档以及程序案例，快速上手编程开发。

## 安装及使用

_注1：Zn 采用 Go 语言开发，安装前请确保 Go 语言编译器已安装在机器上。（建议版本 `>= 1.13`）_  

### 下载
```sh
# 下载 及安装 Zn 语言
$ go get github.com/DemoHn/Zn -o zinc

# 查看 Zn 语言版本
$ zinc -v
```

### 进入交互模式

Zn 支持在命令行中以交互的方式返回结果（亦即REPL）。在命令行中直接输入 `zinc` ，即可进行交互模式。

交互模式运行时，当前行的前面会显示 `Zn>` 做为标识符。在 `>` 号后即可直接输入完整的表达式或者语句。如果中途发现有地方需要修改，即可使用 `左方向键` 及 `右方向键` 移动光标到对应的位置以编辑；输入完成后，敲击回车键即可直接执行，运行结果将直接在后面显示。

以下即为使用交互模式的一个例子，您可以切换到中文输入法，试着体验下：  
  > 注：按 `Ctrl + C` 即可退出交互执行模式。

```sh
$ zinc
Zn> 15 + 25
40
```

### 运行代码文件

Zn 语言目前亦支持执行某个文件中的程序，其格式为 `zinc <文件名>` （如 `zinc 快速排序.zn`）。文件路径可以是相对于当前目录的路径，亦可以是绝对路径。Zn 对于文件后缀名并没有要求，但是这里仍然建议代码文件以 `.zn` 做为后缀名保存。

> ⚠️ 代码文件须以 `utf-8` 编码储存，若以其他编码（包括`gb2312`, `gbk`）执行文件将会报错。

运行结果如下所示：
```sh
$ zinc 快速排序.zn
【-3、-2.3、0、1.32、5、8、12、19】
```

## 示例代码

以下代码片段分别展示 Zn 是如何解决具体问题的；您可以直接复制代码并运行程序，进而观察运行结果；如果想深入学习了解，可点击以下链接：

- 点击 [快速开始](./doc/zh-cn/manual/快速开始.md) 以快速了解并学习 Zn 的基本语法。  
- 点击 [用户手册](./doc/zh-cn/manual/README.md) 以查找 Zn 的全部语法等技术细节。 _TODO: 需要补全及修订内容_  

**1. 温度换算**

> 需求描述：摄氏度和华氏度是两种不同的温度单位；两者换算公式为 `F = 32 + 1.8C`（其中 F 为华氏度，C为摄氏度，如 `24°C = (32+1.8*24)°F = 75.2°F`）

```zn
如何换算温度？
    已知温度、单位
    如果单位 == “摄氏度”：
        返回32 + 1.8 * 温度
    再如单位 == “华氏度”：
        返回{温度 - 32} / 1.8
    否则：
        抛出异常：“无效的温度单位”！

（显示：（换算温度：24、“摄氏度”））  // 显示：75.2
```

**2. 解决「鸡兔同笼」问题**

> 需求描述：「鸡兔同笼」是一道经典算术问题；假设在一个笼子里同时关有数只鸡和兔子，从上面数有35个头，从下面数有94只脚，求解笼子里各有几只兔子和鸡？  
>  
> 求解这道问题有很多种方法；这里采用暴力解法实现，亦即从「全是鸡」到「全是兔」逐一列举，直到脚数满足题设条件为止。如果当遍历完成后，发现脚数不满足条件，那就抛出异常，表明「此题无解」！

```zn
如何解决鸡兔同笼？
    已知总头数、总脚数

    令鸡 = 0，兔 = 0
    注：“这里采用暴力解法，从全是鸡的情况到全是兔的情况逐一列举，
        直到满足条件并返回结果”
    每当鸡 <= 总头数：
        兔 = 总头数 - 鸡        
        如果{2 * 鸡} + {4 * 兔} == 总脚数：
            返回【鸡、兔】
        鸡 = 鸡 + 1

    注：如果上面的循环结束了都还没有返回结果的话，那就只能抛出异常了
    抛出异常：“此题无解”！

令结果 =（解决鸡兔同笼：35、100）
（显示：以「鸡数为{#1}，兔数为{#2}」（格式化：结果#1之文本、结果#2之文本））
```

**3. 对数组进行排序**

> 需求描述：对数组元素进行排序是一个在业务开发当中的高频需求；关于排序的算法汗牛充栋，如插入排序、快速排序、归并排序、冒泡排序等；本次程序演示的即是「冒泡排序」，亦即将一个乱序数组从小到大进行排序；
>  
> 冒泡排序的原理非常简单：它重复地走访过要排序的数列，一次比较两个元素，如果前者大于后者，就把它们交换过来；这个工作重复进行，直到没有再需要交换为止，说明排序已完成。_之所以叫「冒泡排序」，因为越小的元素像是水里的泡泡一样，经过交换慢慢「浮」到数列的项端。_


```zn
如何冒泡排序？
    已知 &数组
    令排序完成 = 假

    每当排序完成 /= 真：
        注：预先假设当前数组已经排序完成，无需交换
        排序完成 = 真
        以I、K遍历数组：
            注：防止访问数组越界
            如果 I >= 数组之长度：
                （结束循环）

            如果 数组#{I} > 数组#{I + 1}：
                以数组（交换：I、I + 1）
                排序完成 = 假
    返回数组

令数据源 =【1.32、5、-3、8、12、0、-2.3、19】

（冒泡排序：数据源）
（显示：数据源） 注：显示【-3、-2.3、0、1.32、5、8、12、19】
```

## 开发清单

- [ ] 完成用户手册的编写 (rev04)
- [X] 补充 `exec` 模块的单元测试
- [X] 开发 `Zn for VSCode` 插件，支持语法高亮
- [X] 添加数据类型的常用方法
- [X] 添加 `对于` 关键字 (rev05)
- [X] 添加异常处理 (rev05)

> 版本说明：Zn 语言一开始从 rev01 一路开发至 rev05，每个 rev 之间语法会有较大的差异，请以
> 最新的代码实现以及文档为准.
>
> 预计在 rev08 完成之后，版本将会升级至 v0.1，此时语法将会基本定形，但是标准库和 API 不会向下兼容
>
> 预计从 v1.0 起会公开发布正式版，此时 语法调整、标准库及API 不会随意删减（但是会添加），调整时会向下兼容
> v1.0 起的版本将会是稳定版本.

## 开源许可

此软件采用 `BSD-3` 开源许可，敬请注意其适用范围：

```
Copyright (c) 2020, Zn Dev Group
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the Zn dev group nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL ZN DEV GROUP BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```