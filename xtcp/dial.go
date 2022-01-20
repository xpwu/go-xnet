package xtcp

import (
  "context"
  "errors"
  "github.com/xpwu/go-pipe/pipe"
  "net"
  "time"
)

// addr.go

func Dial(ctx context.Context, network, addr string) (c net.Conn, err error) {
  // 可能在上层调用时，network都是按照tcp的方式进行，这里需要根据addr再次做一次细分
  if network == tcpNetwork {
    switch {
    case isPipe(addr):
      network = pipeNetwork
    case isUnix(addr):
      network = unixNetwork
      addr = unixAddr(addr)
    }
  }

  switch network {
  case pipeNetwork:
    return pipe.Dial(ctx, addr)
  case unixNetwork,tcpNetwork:
    // net/http/transport.go  DefaultTransport
    return (&net.Dialer{
      Timeout:   30 * time.Second,
      KeepAlive: 30 * time.Second,
    }).DialContext(ctx, network, addr)
  }

  return nil, errors.New("not support network: " + network)
}

func CanProxy(addr string) bool {
  switch {
  case isPipe(addr) || isUnix(addr):
    return false
  default:
    return true
  }
}
