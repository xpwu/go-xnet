package xtcp

import "net"

type VarObject interface {
  RemoteAddr() net.Addr
  LocalAddr() net.Addr
}

func GetVarValue(object VarObject, name string) (value string, ok bool) {
  varMap := map[string]func()string{
    "remote_addr":
    func() string {
      return object.RemoteAddr().String()
    },

    "local_addr":
    func() string {
      return object.LocalAddr().String()
    },
  }

  if f,ok := varMap[name]; ok {
    return f(), ok
  }

  return "", false
}
