package connid

import (
  "fmt"
  "strconv"
  "sync/atomic"
  "time"
)

type Id uint64

func (id Id) String() string {
  return fmt.Sprintf("%x", uint64(id))
}

func ResumeIdFrom(str string) (id Id, err error) {
  i, err := strconv.ParseUint(str, 16, 64)
  if err != nil {
    return
  }

  id = Id(i)
  return
}

var sequence uint32 = 0

// 连接的ID不能直接使用地址的值，go的优化策略可能会运行中改变此变量的地址
//  由此可能就会引起不同时间的连接在同一个地址

// 用过即废，尽最大可能不重复id，防止窜连接
func New() Id {
  t := uint64(time.Now().Unix())
  return Id((t << 32) + uint64(atomic.AddUint32(&sequence, 1)))
}
