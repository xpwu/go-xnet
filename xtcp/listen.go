package xtcp

import (
  "strings"
)

// addr.go

type Listen string

func (l *Listen) Network() string {
  if isPipe(string(*l)) {
    return pipeNetwork
  }
  if isUnix(string(*l)) {
    return unixNetwork
  }

  return tcpNetwork
}

// 能直接使用在net.Listen()的字符串
func (l *Listen) String() string {
  if isPipe(string(*l)) {
    return string(*l)
  }

  if isUnix(string(*l)) {
    return unixAddr(string(*l))
  }

  // 如果只有一个端口，Listen时必须以:开头
  if !strings.Contains(string(*l), ":") {
    return ":" + string(*l)
  }

  return string(*l)
}

func (l *Listen) On() bool {
  return *l != ""
}

func (l *Listen) CanTLS () bool {
  return true
}

func (l *Listen) LogString() string {
  return string(*l)
}
