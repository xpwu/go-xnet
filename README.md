# go-xtcp
扩展的 http/tcp 库，支持服务间的pipe通信方式，也提供了方便创建http/tcp服务的接口。

## 支持的 url 表示方式

```
url  scheme://<hostname>/path/to/

Dial and Listen is hostname

hostname   domainName:port
根据URL.Parse 的解析规则，hostname中不能有 / 符号，可以有 : ，但是最后一个:的后面一定是纯数字的字符串表示port
也可以完全没有:,也就没有端口
为了Dial与Listen时使用相同的格式，方便配置，对hostname约定如下：

1、normal net  与正常使用域名或者ip地址的方式一样  ---> xxx.xxx.xxx:xxx 对于Listen的情况，可以只是写一个port或者 :port

2、pipe --->  pipe:xxx  xxx唯一标识符, 必须是数字。该格式仅用于同一进程内goroutine之间的通信，基于chan而设计。

3、unix --->  unix:|xxx|xxx|xxx|xxx.socket:0  用 | 作为路径分隔符，最后的:0是必须的，
如果是相对路径，可以不要第一个 |。其中|xxx|xxx|xxx|xxx.socket 表示unix系统中的一个文件 '/xxx/xxx/xxx/xxx.socket'，
如果与其他不是基于xtcp写的网络库连接时，其他端根据/xxx/xxx/xxx/xxx.socket文件路径自行填写满足自己要求的unix连接格式。

```

## http 
创建服务   xhttp.SeverAndBlock(xxxx)   
客户端    xhttp.NewClient(), 也可以使用 xhttp.DefaultClient

## tcp
####创建服务   
s := &xtcp.Server{xxxx}    
s.ServeAndBlock()

####客户端    
conn := xtcp.Dial(xxxx)   
c := xtcp.NewConn(ctx, conn)  

