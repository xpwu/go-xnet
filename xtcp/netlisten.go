package xtcp

import (
  "context"
  "crypto/tls"
  "errors"
  "fmt"
  "github.com/xpwu/go-pipe/pipe"
  "golang.org/x/sync/semaphore"
  "net"
  "sync"
)

func NetListen(listen *Listen) (net.Listener, error) {
  switch listen.Network() {
  case pipeNetwork:
    return pipe.Listen(listen.String())
  case unixNetwork, tcpNetwork:
    return net.Listen(listen.Network(), listen.String())
  }

  return nil, errors.New(fmt.Sprintf("not support protocol: %s and address: %s",
    listen.Network(), listen.String()))
}

func NetListenTLS(ln net.Listener, file *TlsFile) (l net.Listener, err error) {
  cert, err := tls.LoadX509KeyPair(file.RealCertFile(), file.RealKeyFile())
  if err != nil {
    return nil, err
  }

  conf := &tls.Config{Certificates: []tls.Certificate{cert}}
  l = tls.NewListener(ln, conf)

  return
}

func NetListenConcurrent(ctx context.Context, ln net.Listener,
  maxConnection ConnectionNum) (l net.Listener, err error) {

  return NetListenConcurrentAndName(ctx, ln, maxConnection, "xtcp")
}

func NetListenConcurrentAndName(ctx context.Context, ln net.Listener,
  maxConnection ConnectionNum, name string) (l net.Listener, err error) {

  return newConL(context.WithValue(ctx, coonNameKey{}, name), ln, maxConnection), nil
}

type conListener struct {
  net.Listener
  sem     *semaphore.Weighted
  ctx     context.Context
  cancelF context.CancelFunc
  once    sync.Once
}

func newConL(ctx context.Context, listener net.Listener, maxConnection ConnectionNum) net.Listener {
  ret := &conListener{
    Listener: listener,
    sem:      semaphore.NewWeighted(maxConnection.Value()),
  }

  ret.ctx,ret.cancelF = context.WithCancel(ctx)

  return ret
}

func (cl *conListener) Accept() (c net.Conn, err error) {
  err = cl.sem.Acquire(cl.ctx, 1)
  if err != nil {
    return
  }

  conn, err := cl.Listener.Accept()
  if err != nil {
    cl.sem.Release(1)
    return
  }

  return newConn(cl.ctx, conn, cl.sem), nil
}

func (cl *conListener) Close() error {
  var err error
  cl.once.Do(func() {
    cl.cancelF()
    err = cl.Listener.Close()
  })

  return err
}
