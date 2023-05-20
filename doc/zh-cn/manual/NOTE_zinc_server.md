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

Zinc Server 应处理两类请求
    1. 加载脚本
    2. 执行程序