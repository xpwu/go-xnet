package xhttp

import (
  "context"
  "github.com/xpwu/go-xnet/xtcp"
  "net"
  "net/http"
  "net/url"
)

// 如果返回了，一定是发生错误了
func SeverAndBlock(ctx context.Context, srv *http.Server, netC *xtcp.Net) error {
  srv.Addr = netC.Listen.String()

  ln, err := xtcp.NetListen(&netC.Listen)
  if err != nil {
    return err
  }
  defer func(ln net.Listener) {
    _ = ln.Close()
  }(ln)

  ln,err = xtcp.NetListenConcurrent(ctx, ln, netC.MaxConnections)
  if err != nil {
    return err
  }
  defer func(ln net.Listener) {
    _ = ln.Close()
  }(ln)

  if !netC.TLS || !netC.Listen.CanTLS() {
    return srv.Serve(ln)
  }

  return srv.ServeTLS(ln, netC.TlsFile.RealCertFile(), netC.TlsFile.RealKeyFile())
}

// Zero means no limit.
func NewClient(maxConn int) *http.Client {
  tr := clone(DefaultTransport)
  tr.MaxConnsPerHost = maxConn

  return &http.Client{Transport:tr}
}

var DefaultTransport http.RoundTripper

// no limit maxConn
var DefaultClient *http.Client

func clone(r http.RoundTripper) *http.Transport {
  return r.(*http.Transport).Clone()
}

func init() {
  tr := clone(http.DefaultTransport)

  tr.DialContext = xtcp.Dial

  oldProxy := tr.Proxy

  tr.Proxy = func(request *http.Request) (url *url.URL, err error) {
    addr := request.URL.Host
    if xtcp.CanProxy(addr) {
      return oldProxy(request)
    }

    return nil, nil
  }

  DefaultTransport = tr

  DefaultClient = &http.Client{Transport: DefaultTransport}
}


