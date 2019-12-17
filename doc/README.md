## `doc/` 文件夹说明 

- `spec`  
  此文档采用 ZnTex 格式编写。ZnTex 是专门为Zn语言所设计的一个基于 `tex` 及 `markdown` 语法的标记语言。它的主要目标就是生成相关语言规范及开发文档。

- `res`  
  存储相关资源，包括`ttf`, `html`, `css`

- `zntex`  
  ZnTex 源代码

### ZnTex

`[TODO]`

1. Command 命令
```tex
\commandname[opt1,opt2,...]{arg1}{arg2}{arg3}...
```

2. Switch 开关
```tex
\switchname ABC XYZ ABC XYZ
```

3. Environment 环境
```tex
\begin{tag}[opt1,opt2...]{arg1}{arg2}

XXX YYY ZZZ \switch XXXYYYZZZ
\commandname[pt1]{Hello}
...

\end{tag}
```

4. 换行符与 `\`, `{`, `}`, `$`
    - `\\`  表示换行  
    - `\bs{}` 表示 `\` 字符本身
    - `\{` 表示 `{` 字符本身
    - `\}` 表示 `}` 字符本身
    - `\$` 表示 `$` 字符本身

5. 全局变量
    - `$varname$` 前面及后面都有 `$` 表示全局变量（没有局部变量一说）

6. 单行注释（不支持多行，也没必要支持）
```tex
normal text % this is a comment 
```