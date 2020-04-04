<div align="center">
  <div style="font-size: 32px;margin-bottom: 4px;">Zn</div>  
  <div style="color: #666">千里之行，始于足下。</div>
  <hr />
</div>

## 简介

`Zn` 是一门 _面向业务_ 的新时代编程语言。

所谓「面向业务」，即是为用户开发业务代码时提供便利，使得用户能够编写出快速、稳定、可维护的代码以实现需求。换言之，和其他「面向计算机」的编程语言不同，`Zn` 更加强调「以人为本」，希望程序能够适应真实世界的需求而不是反过来。

为此，`Zn` 拥有以下独一无二的特性：

- 采用 **中文关键词、标点符号**。这样再也不用为「对专业术语用英文命名变量」之事发愁了。
- 默认使用 **高精度小数** 作为数值并参与运算，杜绝因浮点数计算所带来的计算误差。
  > 这一点对开发金融应用尤为关键。显然，诸如 `0.1 + 0.2 = 0.30000000000000004` 这样的结果在金融应用中是无法忍受的。
- 贴近汉语本身语法，阅读代码可以像阅读文章一样自然。

## 新手上路

_注：`Zn` 采用 Go 语言开发，安装前请确保 Go 语言编译器已安装在机器上。（建议版本 `>= 1.11`）_

1. 下载并安装 `Zn` 语言：
```sh
# 下载 Zn 语言
go get -u github.com/DemoHn/Zn

# 查看 Zn 语言版本，若显示 「Zn语言版本：rv2」或类似字样即表示安装成功
Zn -v
```

2. 

