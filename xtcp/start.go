package xtcp

import (
  "context"
  "fmt"
  "github.com/xpwu/go-log/log"
  "net"
  "time"
)

type Handler interface {
  Handle(conn *Conn)
}

type HandlerFun func(conn *Conn)

func (h HandlerFun) Handle(conn *Conn) {
  h(conn)
}

type Server struct {
  Net     *Net
  Handler Handler
  Name    string
}

func (s *Server) ServeAndBlock() {
  StartAndBlock(s.Name, s.Net, s.Handler.Handle)
}

// Deprecated
func StartAndBlock(name string, config *Net, handler func(conn *Conn)) {

  if !config.Listen.On() {
    return
  }

  ctx, logger := log.WithCtx(context.Background())

  logger.Info("server(" + name + ") listen " + config.Listen.LogString())

  ln, err := NetListen(&config.Listen)
  if err != nil {
    panic(err)
  }
  defer func(ln net.Listener) {
    _ = ln.Close()
  }(ln)

  if config.TLS && config.Listen.CanTLS() {
    ln, err = NetListenTLS(ln, &config.TlsFile)
    if err != nil {
      panic(err)
    }
    defer func(ln net.Listener) {
      _ = ln.Close()
    }(ln)
  }

  ln, err = NetListenConcurrentAndName(ctx, ln, config.MaxConnections, name)
  if err != nil {
    panic(err)
  }
  defer func(ln net.Listener) {
    _ = ln.Close()
  }(ln)

  var tempDelay time.Duration

  for {
    co, err := ln.Accept()

    if err != nil {

      // copy from net/http/server.go and modify
      if ne, ok := err.(net.Error); ok && ne.Temporary() {
        if tempDelay == 0 {
          tempDelay = 5 * time.Millisecond
        } else {
          tempDelay *= 2
        }
        if max := 1 * time.Second; tempDelay > max {
          tempDelay = max
        }
        log.Error(fmt.Sprintf("tcp: Accept error: %v; retrying in %v", err, tempDelay))

        time.Sleep(tempDelay)
        continue
      }

      panic(err)
    }

    tempDelay = 0

    conn := co.(*Conn)

    go func() {
      defer func() {
        if r := recover(); r != nil {
          logger.Fatal(r)
        }
        _ = conn.Close()
      }()

      handler(conn)
    }()
  }
}
