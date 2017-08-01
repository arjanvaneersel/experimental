package counting

import (
  "sync/atomic"
)

var count uint64

func Count() uint64 {
  atomic.AddUint64(&count, 1)
  return count
}

func Fib(n int) int {
  if n < 2 {
    return n
  }
  return Fib(n-1) + Fib(n-2)
}
