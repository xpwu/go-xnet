package xtcp

import (
  "github.com/xpwu/go-cmd/exe"
  "math"
  "path/filepath"
)

type Net struct {
  Listen         Listen `conf:",1、xxx.xxx.xxx.xxx:[0-9] 2、:[0-9] 3、pipe:[0-9] 4、unix:|xxx|xxx|xxx|xxx.socket:0"`
  MaxConnections ConnectionNum `conf:",-1:not limit"`
  TLS            bool
  TlsFile        TlsFile
}

const MaxConnNotLimit = -1

func DefaultNetConfig() *Net {
  return &Net{
    Listen:         "",
    MaxConnections: MaxConnNotLimit,
    TLS:            false,
    TlsFile:        TlsFile{
      PrivateKeyPEMFile: "",
      CertPEMFile:       ""},
  }
}

type ConnectionNum int64

func (n ConnectionNum) Value() int64 {
  if n <= 0 {
    return math.MaxInt64
  }

  return int64(n)
}

type TlsFile struct {
  PrivateKeyPEMFile string  `conf:",support relative path, must PEM encode data"`
  CertPEMFile       string  `conf:",support relative path, must PEM encode data"`
}

func (t *TlsFile) RealKeyFile() string {
  if filepath.IsAbs(t.PrivateKeyPEMFile) {
    return t.PrivateKeyPEMFile
  }

  return filepath.Join(exe.Exe.AbsDir, t.PrivateKeyPEMFile)
}

func (t *TlsFile) RealCertFile() string {
  if filepath.IsAbs(t.CertPEMFile) {
    return t.CertPEMFile
  }

  return filepath.Join(exe.Exe.AbsDir, t.CertPEMFile)
}
