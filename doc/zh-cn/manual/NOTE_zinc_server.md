> 我们假设 Zn 的应用场景主要是寄生在由其他编程语言（python, JS, etc.） 所提供的服务里面，而且我们不做 Service，只做短链场景，这样就避免了内存泄漏等运营搞不定的问题

### FastCGI 协议 NOTE

NGINX 之类的 WebServer 通过发送 Records 请求信息；而 FastCGI Server 收到请求后进行处理，处理完之后再用 Records 把信息拎回去

Records 的格式如下：

0 ------- 1 ------ 2 --------- 3 ---- 4
| version | type   |   reqID [HL]     |
| contentLen [HL]  |   paddingLen     |
| reserved |   contentData...         |
    ...    |   paddingData            |


这里注意到几件事情：

1. reqID 最多只能 65536 个，所以 FastCGI Server 不能搞一堆长链把请求会 hang 住 （就像 php-fpm 一般单机也就开个 200 个worker）

2. contenLength 最长只能是 65536，所以当 content 过长时，只能把内容切分成几片发过去


从协议上看，NGINX 在发送请求时应该会按照顺序发送以下几个 Record:

1. FCGI_BEGIN_REQUEST: 开始请求
2. FCGI_PARAMS：把一些和请求相关的环境变量扔进来（历史遗留问题：CGI就是依赖这玩意读取数据的）
3. FCGI_STDIN：把请求数据放进去（比如 POST data 之类的，这个需要试验）


NGINX 在接收到 HTTP 请求后，先简单解析这协议，然后把它转成 fastCGI 协议发送给后端服务（php-fpm之类的），服务再把HTTP响应（content-type: text/html 开头...） 通过fastCGI 协议直接还回去

### 如何多模态处理请求？

> 我们在调查了几个编程语言(Go, JS, Python) 的 fastCGI 库，发现大家都对 fastCGI 库做了一层封装，将 fastCGI 输入请求自动转化成 HTTP 请求（比如将 `REQUEST_METHOD` 参数写进 http.Request 对象里）

> 我设想的其他语言(比如 nodejs) 调用 zinc 脚本模式如下：

```js
// 准备灌给 zinc 的变量
const pet0 = {
    "name": "旺财",
    "type": "dog",
    "age": 12
}

const store0 = "北大医院"

// 指定执行文件（或者 shared memory）
const result = zincReaction("./宠物医院.zn", 医院=store0, 宠物=pet0)
// result = "没得治"
```

```zn
// 宠物医院.zn

如果医院 == 「北大医院」：
    输出「没得治」
再如宠物#“age” > 20：
    输出「没得治」

输出「有得治」
```

> 对于 znServer 而言，如果 ZINC_ADAPTER 选择 http，那么输入值就按照 application 的格式转成 object；而后将输出对象转成文本以 text/html 形式输出

> 相对于直接架一个 HTTP Server，使用 fastCGI 的好处非常明显：可以在 FCGI_PARAMS 里面加上自定义的参数，比如 ZINC_ADAPTER, SCRIPT_FILENAME -- 这些「加上私货」的参数可以直接控制 ZnServer 的处理模式以及执行代码的文件位置。

> 在使用 `zinc` (不是 `zinc-server`) 执行程序时，最终向 stdout 输出的分为两个部分：

1. `显示` 语句所输出的文本
2. `输出` 语句所调用的值 （或者最后一个表达式）

    - 如果最后一个表达式是`空`，那就不显示

简单来说，把 `zinc` (以及 `zinc-server`) 当作一个大型的 jupyter notebook，把条件和表达式怼进去，出来最终的结果以及之前需要输出的值

### ZnServer 处理请求的流程

**总纲**
1. 输入和输出都使用 HTTP 格式
2. 输入的 Body 取决于 HTTP Header 中的 Content-Type 字段 (MIME Type)


**ZINC_ADAPTER 支持的模式**

1. playground

正常情况下，从 fastcgi body 中读取待执行的程序，再输出具体的值作为 http body 返回来。这玩意特别适合直接怼进 HTTP 页面时作为 Playground 运行时

2. http_handler

如同 PHP-FPM 的处理方式一样，从 WebServer 中接收HTTP请求，将请求参数灌到输入里面，再输出具体的值作为 http body 返回来。

如果 ZINC_ADAPTER 不是上述两个值或者根本就没给，直接返回 403


