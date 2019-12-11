<h1 align="center">Zn</h1>

<p align="center">千里之行，始于足下。</p>

## 简介 - What's Zn?

`Zn` （锌, Zion）是一门旨在帮助中文用户编写程序的脚本语言。它的命名来源于第30号元素「锌」的化学符号；同时「锌」的汉语发音正好与「新」相同，期望着这门「新」的语言能够在未来发展出自己的一片天地。

设计这门语言的初衷，来源于我在平时开发当中观察到的几个方面：

  - 中文用户的英文水平并没有高到无缝切换的地步。  
     
    虽然编程语言的关键字（基本上是英文单词）只有寥寥数十个，学习记住它自然十分容易；不过具体到业务代码中，大部分变量都需要用英文命名，于是这样的`中-英`转换就给开发人员带来不少负担了。更糟糕的是，当业务涉及到一些非常专业的领域时，对一些专业术语进行命名就成为了很大的挑战，而这些开销其实是没有必要的。

  - 开发业务代码并不需要太多抽象的算法。

    事实上，大部分业务代码的开发都可以用四个字母来概括：`CRUD`，亦即`增删改查`。和开发系统程序不同，开发业务代码并不需要很多精妙而难懂的算法，它需要更多对业务行为的抽象，讲求的是开发效率及可阅读性。

有鉴于此，`Zn` 语言在设计之初就向着以下几个特点迈进：

  - 使用中文关键字及标点符号，文法上更加贴近汉语语法。
  - 不再支持算术表达式（但是算术是支持的！）以换取更高的可读性。
  - 写出来的代码不仅要方便阅读和理解，而且颜值要高。

## 示例 - Demo

1. 计算两数之和
```zn
如何求和？
    已知左数，右数@整数
    返回【左数，右数】之和

显示：（求和：2，4） 注：应返回6
```

2. 动物庄园
```zn
定义动物：
    其名，年龄为空
    
    是为：名，年龄
    
    如何比较年龄？
        已知小A@动物

        如果小A之年龄大于其年龄：
            令差值为【小A之年龄，其年龄】之差
            以「「#1」的年龄比「#2」的年龄大#3岁」而连缀：小A之年龄，其年龄，差值
        否则：
            令差值为【其年龄，小A之年龄】之差
            以「「#1」的年龄比「#2」的年龄小#3岁」而连缀：小A之年龄，其年龄，差值
    
    如何显示姓名？ 其名

    如何显示年龄？ 其年龄

令小马成为动物：「小马」，20
令小象成为动物：「小象」，8

注1：最后应该显示 「小象」的年龄比「小马」的年龄小12岁
显示：（在小马中比较年龄：小象）
```

## 技术规范

TODO