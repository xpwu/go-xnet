package xtcp

import (
  "strings"
)

/**

url  scheme://<hostname>/path/to/

Dial and Listen is hostname

hostname   domainName:port
根据URL.Parse 的解析规则，hostname中不能有 / 符号，可以有 : ，但是最后一个:的后面一定是纯数字的字符串表示port
也可以完全没有:,也就没有端口
为了Dial与Listen时使用相同的格式，方便配置，对hostname约定如下：

1、normal net  与正常使用域名或者ip地址的方式一样  ---> xxx.xxx.xxx:xxx 对于Listen的情况，可以只是写一个port或者 :port

2、pipe --->  pipe:xxx  xxx唯一标识符, 必须是数字

3、unix --->  unix:|xxx|xxx|xxx|xxx.socket:0  用 | 作为路径分隔符，最后的:0是必须的，否则解析端口时会错误，但实际中并不使用，
如果是相对路径，可以不要第一个 |

*/

func isPipe(addr string) bool {
  return strings.HasPrefix(addr, "pipe:")
}

func isUnix(addr string) bool {
  return strings.HasPrefix(addr, "unix:")
}

func unixAddr(addr string) string {
  colonF := strings.Index(addr, ":")
  colonL := strings.LastIndex(addr, ":")
  if colonF >= colonL {
    return ""
  }

  return strings.Replace(addr[colonF+1:colonL], "|", "/", -1)
}

const (
  pipeNetwork = "pipe"
  unixNetwork = "unix"
  tcpNetwork = "tcp"
)
