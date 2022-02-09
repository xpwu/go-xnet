package xtcp

import (
  "context"
  "fmt"
  "github.com/xpwu/go-log/log"
  "golang.org/x/sync/semaphore"
  "net"
  "sync"
  "time"
  "unsafe"
)

type Conn struct {
  net.Conn
  mu         chan struct{}
  ctx        context.Context
  cancelFunc context.CancelFunc
  sem        *semaphore.Weighted
  once       sync.Once
  time       time.Time
}

// 主要使用在client dial后的net.Conn 生成新的xtcp.Conn
func NewConn(ctx context.Context, c net.Conn) *Conn {
  ctx,fun := context.WithCancel(ctx)
  ctx,logger := log.WithCtx(ctx)

  ret := &Conn{
    Conn:       c,
    mu:         make(chan struct{}, 1),
    ctx:        ctx,
    sem:        nil,
    cancelFunc: fun,
    time: time.Now(),
  }

  logger.PushPrefix(fmt.Sprintf("tcp conn(%s) to %s", ret.Id(), c.RemoteAddr().String()))

  return ret
}

func newConn(ctx context.Context, c net.Conn, sem *semaphore.Weighted) *Conn {
  ctx,fun := context.WithCancel(ctx)
  ctx,logger := log.WithCtx(ctx)

  ret := &Conn{
    Conn:       c,
    mu:         make(chan struct{}, 1),
    ctx:        ctx,
    sem:        sem,
    cancelFunc: fun,
    time: time.Now(),
  }

  logger.PushPrefix(fmt.Sprintf("tcp conn(%s) from %s", ret.Id(), c.RemoteAddr().String()))

  return ret
}

func (c *Conn) Write(b []byte) (n int, err error) {
  c.mu <- struct{}{}
  defer func() {
    <-c.mu
  }()

  n, err = c.Conn.Write(b)
  return
}

// 对于不支持 writev 的系统，将会使用 循环的方式，为了保证 buffer 是一个整体被写入，所以加一个同步
func (c *Conn) WriteBuffers(buffers net.Buffers) (n int, err error) {
  c.mu <- struct{}{}
  defer func() {
    <-c.mu
  }()

  n64, err := buffers.WriteTo(c.Conn)
  n = int(n64)
  return
}

func (c *Conn) Close() error {
  var err error
  c.once.Do(func() {
    if c.sem != nil {
      c.sem.Release(1)
    }
    c.cancelFunc()
    _,logger := log.WithCtx(c.ctx)
    logger.Debug("close connection")
    err = c.Conn.Close()
  })
  return err
}

func (c *Conn) Context() context.Context {
  return c.ctx
}

func (c *Conn) GetVar(name string) (value string, ok bool) {
  return GetVarValue(c, name)
}

// 直接使用地址，可能引起不同时间的连接在同一个地址
func (c *Conn) Id() string {
  return fmt.Sprintf("%x+%x", unsafe.Pointer(c), c.time.UnixNano())
}
